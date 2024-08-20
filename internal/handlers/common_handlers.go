package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"bookmosaic/internal/data"
	"bookmosaic/internal/models"
)

type Handlers struct {
	DB *data.DB
}

func NewHandlers(db *data.DB) *Handlers {
	return &Handlers{DB: db}
}

// Route /book/
func (h *Handlers) BookHandler(res http.ResponseWriter, req *http.Request) {
	idStr := req.URL.Path[len("/book/"):]

	_, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(res, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch req.Method {
	case http.MethodGet:
		h.GetBookByIDHandler(res, req)
	case http.MethodDelete:
		h.DeleteBookByIDHandler(res, req)
	case http.MethodPost:
		h.UpdateBookByIDHandler(res, req)
	default:
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// GET /books&limit=n&offset=m
func (h *Handlers) GetBooksHandler(res http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	limitStr := query.Get("limit")
	offsetStr := query.Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 0
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	books, err := h.DB.GetBooks(context.Background(), limit, offset)
	if err != nil {
		books = []models.SummaryBook{}
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(res).Encode(books); err != nil {
		http.Error(res, "Error encoding json", http.StatusInternalServerError)
		return
	}
}

// GET /book/{id}
func (h *Handlers) GetBookByIDHandler(res http.ResponseWriter, req *http.Request) {
	idStr := req.URL.Path[len("/book/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(res, "Invalid ID", http.StatusBadRequest)
		return
	}

	book, err := h.DB.GetBookByID(context.Background(), id)
	if err != nil {
		http.Error(res, "Error retrieving book", http.StatusInternalServerError)
		log.Printf("Error retrieving book: %v", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(book)
}

// DELETE /book/{id}
func (h *Handlers) DeleteBookByIDHandler(res http.ResponseWriter, req *http.Request) {
	idStr := req.URL.Path[len("/book/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(res, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.DB.DeleteBookByID(context.Background(), id)
	if err != nil {
		http.Error(res, "Error delete book", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}

// PUT /book/{id}
func (h *Handlers) UpdateBookByIDHandler(res http.ResponseWriter, req *http.Request) {
	idStr := req.URL.Path[len("/book/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(res, "Invalid ID", http.StatusBadRequest)
		return
	}

	var book models.Book

	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&book)
	if err != nil {
		http.Error(res, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = h.DB.UpdateBookByID(context.Background(), id, book)
	if err != nil {
		http.Error(res, "Failed to update book", http.StatusNotFound)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}

// POST /book
func (h *Handlers) AddNewBookHandler(res http.ResponseWriter, req *http.Request) {
	var book models.Book
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&book)
	if err != nil {
		http.Error(res, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = h.DB.AddNewBook(context.Background(), book)
	if err != nil {
		http.Error(res, "Failed to add new book", http.StatusInternalServerError)
		log.Printf("Failed to insert new book: %v", err)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}
