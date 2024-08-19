package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	DATABASE_URL := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), DATABASE_URL)
	if err != nil {
		log.Fatalf("Unable tp connect to database: %v", err)
	}

	var title string
	err = conn.QueryRow(context.Background(), "SELECT title FROM books LIMIT 1").Scan(&title)
	if err != nil {
		log.Fatalf("QueryRow failed: %v", err)
	}

	fmt.Println(title)

}
