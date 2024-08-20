package data

import (
	"context"
	"fmt"

	"bookmosaic/internal/models"
)

// Get books with offset
func (db *DB) GetBooks(ctx context.Context, limit, offset int) ([]models.SummaryBook, error) {
	summBooks := make([]models.SummaryBook, 0, limit)
	rows, err := db.DB.Query(ctx, "SELECT id, title, cover FROM books ORDER BY id LIMIT $1 OFFSET $2", limit, offset)

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var title string
		var cover string

		if err := rows.Scan(&id, &title, &cover); err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}

		bookElement := models.SummaryBook{
			ID:         id,
			Title:      title,
			CoverImage: cover,
		}

		summBooks = append(summBooks, bookElement)

		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("row error: %w", err)
		}
	}

	return summBooks, nil
}

// Get book by id
func (db *DB) GetBookByID(ctx context.Context, id int) (models.Book, error) {
	book := models.Book{}

	err := db.DB.QueryRow(
		ctx, "SELECT * FROM books WHERE id = $1", id,
	).Scan(
		&book.ID,
		&book.ISBN,
		&book.Title,
		&book.Author,
		&book.Year,
		&book.Summary,
		&book.CoverImage,
	)
	if err != nil {
		return book, fmt.Errorf("query row error: %w", err)
	}

	return book, nil
}

// Delete book by ID
func (db *DB) DeleteBookByID(ctx context.Context, id int) error {
	tag, err := db.DB.Exec(ctx, "DELETE FROM books WHERE id = $1", id)

	if err != nil {
		return fmt.Errorf("failed to delete book with id %d: %w", id, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no book found with id %d: %w", id, err)
	}

	return nil
}

// Update book data
func (db *DB) UpdateBookByID(ctx context.Context, id int, upd models.Book) error {
	tag, err := db.DB.Exec(
		ctx,
		"UPDATE books SET isbn=$1, title=$2, author=$3, year=$4, summary=$5, cover=$6 WHERE id=$7",
		&upd.ISBN, &upd.Title, &upd.Author, &upd.Year, &upd.Summary, &upd.CoverImage, id,
	)

	if err != nil {
		return fmt.Errorf("failed to update book with id %d: %w", id, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no book found with id %d: %w", id, err)
	}

	return nil
}

// Add new book
func (db *DB) AddNewBook(ctx context.Context, book models.Book) error {
	_, err := db.DB.Exec(
		ctx, "INSERT INTO books (isbn, title, author, year, summary, cover) VALUES ($1, $2, $3, $4, $5, $6)",
		book.ISBN, book.Title, book.Author, book.Year, book.Summary, book.CoverImage,
	)

	if err != nil {
		return fmt.Errorf("failed to insert new book: %w", err)
	}

	return nil
}
