package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
)

var errDuplicatedEvent = errors.New("error inserting duplicated event")

type DB struct {
	logger *slog.Logger
	db     *sql.DB
}

type DBInterface interface {
}

func NewDB(logger *slog.Logger) (*DB, error) {
	connString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"), // TODO: password must be from a secret manager
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
		"disable", // for production it should be confifured properly
	)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return &DB{
		logger: logger,
		db:     db,
	}, nil
}
