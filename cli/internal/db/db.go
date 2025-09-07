package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Global database connection pool
var pool *pgxpool.Pool

// Connect establishes a connection to the PostgreSQL database
func Connect(dbPath string) {
	if dbPath == "" {
		log.Println("DB_PATH not set, skipping DB connection")
		return // Running in dry mode, no DB connection
	}

	var err error
	pool, err = pgxpool.New(context.Background(), dbPath)
	if err != nil {
		log.Println("Unable to connect to database:", err)
		pool = nil
		return
	}

	fmt.Println("Connected to PostgreSQL database!")
}

// InsertAppUsage inserts/updates an application usage record in the database
func InsertAppUsage(appName string, minutes float32, source string) {
	if pool == nil {
		return // No DB connection, skip insertion
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int
	// Check if there's an existing record for today
	err := pool.QueryRow(ctx,
		`SELECT id
		 FROM daily_usage
		 WHERE name=$1 AND source=$2 AND usage_date= CURRENT_DATE`,
		appName, source).Scan(&id)
	
	if err != nil {
		// No existing record, insert new
		_, err := pool.Exec(ctx,
			`INSERT INTO daily_usage (name, usage_date, minutes_used, source)
 			 VALUES ($1, CURRENT_DATE, $2, $3)`,
			appName, minutes, source)
		
		if err != nil {
			log.Println("Error inserting app usage:", err)
		}
		return
	}

	// Existing record found, update it
	_, err = pool.Exec(ctx,
		`UPDATE daily_usage
		 SET minutes_used = minutes_used + $1
		 WHERE id = $2`,
		minutes, id)
	
	if err != nil {
		log.Println("Error updating app usage:", err)
	}
}

// Close closes the database connection pool
func Close() {
	if pool != nil {
		pool.Close()
		fmt.Println("Database connection closed.")
	}
}