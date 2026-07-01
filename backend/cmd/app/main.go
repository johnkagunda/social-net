package main

import (
	"database/sql"
	"fmt"
	"log"
	"social/pkg/db/sqlite"
	"social/pkg/sqllite"
	"social/server"

	_ "modernc.org/sqlite"
)

func main() {
	// Initialize database
	db, err := sql.Open("sqlite", "./social_network.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()
	sqlite.DB = db

	// Run migrations
	if err := sqllite.RunMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Create and start server
	srv := server.NewServer(db)

	fmt.Println("Server starting on port 8080...")
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
