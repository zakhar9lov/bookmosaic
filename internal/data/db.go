package data

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	DB *pgx.Conn
}

func NewDBConnection(ctx context.Context, connString string) (*DB, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	log.Println("Successfuly connected to database!")
	return &DB{DB: conn}, nil
}

func (db *DB) Close(ctx context.Context) {
	db.DB.Close(ctx)
}
