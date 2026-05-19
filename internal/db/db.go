package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"pack-calculator/internal/common"

	"github.com/lib/pq"
)

var errDuplicatedEvent = errors.New("error inserting duplicated event")

type DB struct {
	logger *slog.Logger
	db     *sql.DB
}

type DBInterface interface {
	ClearPackSize(ctx context.Context) error
	SavePackSize(ctx context.Context, packSizeBatch *common.PackSizeBatch) error
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

func (db *DB) ClearPackSize(ctx context.Context) error {
	_, err := db.db.ExecContext(
		ctx,
		`	
			DELETE FROM pack_size;
		`,
	)

	return err
}

func (db *DB) SavePackSize(ctx context.Context, packSizeBatch *common.PackSizeBatch) error {
	_, err := db.db.ExecContext(
		ctx,
		`	
			INSERT INTO public.pack_size(size)
        	SELECT unnest($1::int[])
		`,
		pq.Array(packSizeBatch.Sizes),
	)

	return err
}
