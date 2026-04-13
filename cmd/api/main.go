package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Saident/The-Daily-Bouquet-Crafter-Backend/internal/database"
	"github.com/Saident/The-Daily-Bouquet-Crafter-Backend/internal/game"
	"github.com/joho/godotenv"
)

// Defines a single flower dropped into the vase
type BouquetItem struct {
	AssetURL  string  `json:"asset_url"`
	XPosition float64 `json:"x_position"`
	YPosition float64 `json:"y_position"`
	Rotation  float64 `json:"rotation"`
	ZIndex    int     `json:"z_index"`
}

type BouquetResponse struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Date      string        `json:"date"`
	IsSpecial bool          `json:"is_special"`
	Items     []BouquetItem `json:"flowers"`
}

// The full payload sent from React when she clicks "Save"
type SaveBouquetRequest struct {
	Title string        `json:"title"`
	Items []BouquetItem `json:"items"`
}

// enableCORS allows the React frontend to talk to this Go server
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func main() {
	// 1. Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// 2. Initialize Database
	dbURL := os.Getenv("DB_URL")
	db, err := database.InitDB(dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close() // Ensure the connection closes when the app stops

	// 3. Setup Routes

	// GET: /api/inventory
	http.HandleFunc("/api/inventory", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// (Future feature: check DB to see if she unlocked the secret)
		hasUnlockedSecret := false
		dailyFlowers := game.GetDailyInventory(hasUnlockedSecret)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"date":    time.Now().Format("2006-01-02"),
			"flowers": dailyFlowers,
		})
	}))

	// GET: /api/bouquets (Fetch the Greenhouse Gallery)
	http.HandleFunc("/api/bouquets", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 1. Grab all the overarching bouquets, newest first
		rows, err := db.Query("SELECT id, title, created_at, is_special FROM bouquets ORDER BY created_at DESC")
		if err != nil {
			http.Error(w, "Failed to fetch bouquets", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var bouquets []BouquetResponse
		for rows.Next() {
			var b BouquetResponse
			var createdAt time.Time
			if err := rows.Scan(&b.ID, &b.Title, &createdAt, &b.IsSpecial); err != nil {
				continue
			}
			// Format the date to match the aesthetic UI (e.g., "April 13, 2026")
			b.Date = createdAt.Format("January 2, 2006")
			bouquets = append(bouquets, b)
		}

		// 2. Loop through and grab the specific flowers for each bouquet
		for i := range bouquets {
			itemRows, err := db.Query(
				"SELECT asset_url, x_position, y_position, rotation, z_index FROM bouquet_items WHERE bouquet_id = $1",
				bouquets[i].ID,
			)
			if err != nil {
				continue
			}

			var items []BouquetItem
			for itemRows.Next() {
				var it BouquetItem
				itemRows.Scan(&it.AssetURL, &it.XPosition, &it.YPosition, &it.Rotation, &it.ZIndex)
				items = append(items, it)
			}
			itemRows.Close()
			bouquets[i].Items = items
		}

		w.Header().Set("Content-Type", "application/json")
		// Prevent returning 'null' if the database is completely empty
		if bouquets == nil {
			bouquets = []BouquetResponse{}
		}
		json.NewEncoder(w).Encode(bouquets)
	}))

	// POST: /api/bouquets/save
	http.HandleFunc("/api/bouquets/save", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req SaveBouquetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// The Easter Egg Logic
		isSpecial := false
		message := "Bouquet saved successfully!"
		if strings.ToLower(req.Title) == "bellaa" {
			isSpecial = true
			message = "Secret unlocked! A rare glowing lotus has been added to your inventory for tomorrow."
		}

		// Insert the overarching bouquet record
		var bouquetID string
		err = db.QueryRow(
			"INSERT INTO bouquets (title, is_special) VALUES ($1, $2) RETURNING id",
			req.Title, isSpecial,
		).Scan(&bouquetID)

		if err != nil {
			log.Printf("Error saving bouquet to DB: %v", err)
			http.Error(w, "Failed to save bouquet", http.StatusInternalServerError)
			return
		}

		// Loop through the flowers and save their exact positions
		for _, item := range req.Items {
			_, err = db.Exec(
				"INSERT INTO bouquet_items (bouquet_id, asset_url, x_position, y_position, rotation, z_index) VALUES ($1, $2, $3, $4, $5, $6)",
				bouquetID, item.AssetURL, item.XPosition, item.YPosition, item.Rotation, item.ZIndex,
			)
			if err != nil {
				log.Printf("Error saving bouquet item: %v", err)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     "success",
			"message":    message,
			"special":    isSpecial,
			"bouquet_id": bouquetID,
		})
	}))

	// DELETE: /api/bouquets/delete
	http.HandleFunc("/api/bouquets/delete", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Grab the ID from the URL (e.g., /api/bouquets/delete?id=123)
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Missing bouquet ID", http.StatusBadRequest)
			return
		}

		// The ON DELETE CASCADE in our schema automatically deletes the flowers from bouquet_items too!
		_, err := db.Exec("DELETE FROM bouquets WHERE id = $1", id)
		if err != nil {
			log.Printf("Error deleting bouquet: %v", err)
			http.Error(w, "Failed to delete bouquet", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Bouquet completely removed",
		})
	}))

	// 4. Start the Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Backend server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
