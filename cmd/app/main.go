package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"bookmosaic/internal/data"
	"bookmosaic/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {

	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Setting up a logger to write to a file
	log.SetOutput(file)

	// Loading environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// conn is a structure that defines interaction with the database
	DATABASE_URL := os.Getenv("DATABASE_URL")
	conn, err := data.NewDBConnection(context.Background(), DATABASE_URL)
	if err != nil {
		log.Fatalf("Error to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	// h — this is a structure that provides handlers
	h := handlers.NewHandlers(conn)
	mux := http.NewServeMux()
	mux.HandleFunc("/books", h.GetBooksHandler)
	mux.HandleFunc("/book/", h.BookHandler)
	mux.HandleFunc("/book", h.AddNewBookHandler)

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Failed to listen and serve: %v", err)
	}

}
