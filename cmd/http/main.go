package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

func main() {
	database_url := os.Getenv("DATABASE_URL")

	log.Println("Hello world")
	log.Println(database_url)

	db, err := sql.Open("postgres", database_url)
	if err != nil {
		panic(err)
	}
	defer db.Close() // Close the database connection when you're done
	// Perform a query
	query := "SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE';"
	rows, err := db.Query(query)
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
