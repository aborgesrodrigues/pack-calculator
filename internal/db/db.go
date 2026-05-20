package db

import (
	"context"
	"database/sql"
	"encoding/json"
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
	ClearPackSize(context.Context) error
	SavePackSize(context.Context, *common.PackSizeBatch) error
	GetPackSizes(context.Context) ([]int, error)
	SaveOrder(context.Context, *common.Order) error
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

	// check if database is up
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Database is down: %v", err)
	}

	return &DB{
		logger: logger,
		db:     db,
	}, nil
}

func (db *DB) GetPackSizes(ctx context.Context) ([]int, error) {
	var sizes64 []int64

	if err := db.db.QueryRowContext(
		ctx,
		`	
			SELECT COALESCE(array_agg(size ORDER BY size DESC), '{}') FROM pack_size;
		`,
	).Scan(pq.Array(&sizes64)); err != nil {
		return nil, err
	}

	// cast to []int
	sizes := make([]int, len(sizes64))
	for i := range sizes {
		sizes[i] = int(sizes64[i])
	}

	return sizes, nil
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

func (db *DB) SaveOrder(ctx context.Context, order *common.Order) error {
	packs := order.Packs
	if packs == nil {
		packs = map[int]int{}
	}

	packsJSON, err := json.Marshal(packs)
	if err != nil {
		return fmt.Errorf("error marshaling packs: %w", err)
	}

	_, err = db.db.ExecContext(
		ctx,
		`INSERT INTO public."order" (id, amount_items, packs) VALUES ($1, $2, $3::jsonb)`,
		order.ID,
		order.AmountItems,
		packsJSON,
	)

	return err
}
