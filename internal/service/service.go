package service

import (
	"context"
	"fmt"
	"log/slog"
	"pack-calculator/internal/common"
	"pack-calculator/internal/db"
)

type Service struct {
	logger *slog.Logger
	db     db.DBInterface
}

type SVCInterface interface {
	SavePackSize(context.Context, *common.PackSizeBatch) error
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

func (s *Service) SavePackSize(ctx context.Context, packSizeBatch *common.PackSizeBatch) error {
	// clear all previous pack sizes
	if err := s.db.ClearPackSize(ctx); err != nil {
		return fmt.Errorf("error clearing pack sizes: %v", err)
	}

	// save the pack size list passed
	if err := s.db.SavePackSize(ctx, packSizeBatch); err != nil {
		return fmt.Errorf("error inserting pack sizes: %v", err)
	}

	return nil
}
