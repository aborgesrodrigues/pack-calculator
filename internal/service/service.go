package service

import (
	"fmt"
	"log/slog"
	"pack-calculator/internal/db"
)

type Service struct {
	logger *slog.Logger
	db     db.DBInterface
}

type SVCInterface interface {
}

func NewService(logger *slog.Logger) (*Service, error) {
	db, err := db.NewDB(logger)

	if err != nil {
		return nil, fmt.Errorf("error instantiating DB: %w", err)
	}

	return &Service{
		logger: logger,
		db:     db,
	}, nil
}
