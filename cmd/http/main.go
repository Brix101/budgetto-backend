package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Brix101/budgetto-backend/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5" // Import the PostgreSQL driver
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	database_url := config.GetConfig().DATABASE_URL

	log.Println(database_url)

	db, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		panic(err)
	}

	query := `
	SELECT table_name 
	FROM information_schema.tables 
	WHERE table_schema='public' 
	AND table_type='BASE TABLE' 
	AND table_name NOT LIKE '%db_version%';`

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		panic(err)
	}
	defer rows.Close() // Close the result set when you're done

	var tableNames []string

	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			panic(err)
		}

		// Append the tableName to the slice
		tableNames = append(tableNames, tableName)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		responseJSON, err := json.Marshal(tableNames)
		if err != nil {
			http.Error(w, "Error converting to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	})
	http.ListenAndServe(":3000", r)
}
