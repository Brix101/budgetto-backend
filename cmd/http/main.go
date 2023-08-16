package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5" // Import the PostgreSQL driver
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	database_url := os.Getenv("DATABASE_URL")

	log.Println("Hello world")
	log.Println(database_url)

	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	query := "SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE';"
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		panic(err)
	}
	defer rows.Close() // Close the result set when you're done

	// Iterate over the result set
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Table Name: %s\n", tableName)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
}
