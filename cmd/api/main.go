// cmd/api/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/Saident/The-Daily-Bouquet-Crafter-Backend/internal/database"
)

// Request payload for saving a bouquet
type SaveBouquetRequest struct {
	Title string `json:"title"`
	// Items will be added here later when we wire up the frontend
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// Initialize Database
	dbURL := os.Getenv("DB_URL")
	_, err = database.InitDB(dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Setup a simple route to handle saving the bouquet
	http.HandleFunc("/api/bouquets/save", func(w http.ResponseWriter, r *http.Request) {
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
			// Future logic: flag tomorrow's daily_inventory row to include the rare asset
		}

		// (Here is where we will write the SQL INSERT statement later)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": message,
			"special": isSpecial,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Backend server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
