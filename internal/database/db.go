// internal/database/db.go
package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to the database!")
	createTables(db)

	return db, nil
}

func createTables(db *sql.DB) {
	// 1. Table for the daily available flowers
	dailyInventoryTable := `
	CREATE TABLE IF NOT EXISTS daily_inventory (
		id SERIAL PRIMARY KEY,
		date DATE UNIQUE NOT NULL,
		flower_assets JSONB NOT NULL
	);`

	// 2. Table for the saved bouquets
	bouquetsTable := `
	CREATE TABLE IF NOT EXISTS bouquets (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		title VARCHAR(255) NOT NULL,
		is_special BOOLEAN DEFAULT FALSE
	);`

	// 3. Table for the exact position of each flower in a bouquet
	bouquetItemsTable := `
	CREATE TABLE IF NOT EXISTS bouquet_items (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		bouquet_id UUID REFERENCES bouquets(id) ON DELETE CASCADE,
		asset_url VARCHAR(255) NOT NULL,
		x_position FLOAT NOT NULL,
		y_position FLOAT NOT NULL,
		rotation FLOAT NOT NULL,
		z_index INT NOT NULL
	);`

	tables := []string{dailyInventoryTable, bouquetsTable, bouquetItemsTable}

	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			log.Fatalf("Error creating table: %v", err)
		}
	}
	fmt.Println("Database schemas initialized.")
}
