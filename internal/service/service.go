package service

import (
	"context"
	"fmt"
	"log/slog"
	"pack-calculator/internal/common"
	"pack-calculator/internal/db"

	"github.com/google/uuid"
)

type Service struct {
	logger *slog.Logger
	db     db.DBInterface
}

type SVCInterface interface {
	GetPackSizes(context.Context) (*common.PackSizeBatch, error)
	SavePackSize(context.Context, *common.PackSizeBatch) error
	Calculate(context.Context, *common.Order) (*common.Order, error)
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

func (s *Service) GetPackSizes(ctx context.Context) (*common.PackSizeBatch, error) {
	sizes, err := s.db.GetPackSizes(ctx)
	if err != nil {
		return nil, fmt.Errorf("error get pack sizes: %v", err)
	}

	return &common.PackSizeBatch{Sizes: sizes}, nil
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

func (s *Service) Calculate(ctx context.Context, order *common.Order) (*common.Order, error) {
	order.ID = uuid.New()

	sizes, err := s.db.GetPackSizes(ctx)
	if err != nil {
		return nil, fmt.Errorf("error get pack sizes: %v", err)
	}

	packs := calculate(order.AmountItems, sizes, make(map[int]int))
	s.logger.Info("Result", "packs", packs)

	order.Packs = packs

	if err := s.db.SaveOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("error saving order: %v", err)
	}

	return order, nil
}

func calculate(amountItems int, sizes []int, used map[int]int) map[int]int {
	rest := amountItems
	for _, size := range sizes {
		// find the first pack size that can support the size
		if amountItems >= size {
			used[size] += 1
			rest -= size

			// if rest is negative means that calculation is done
			if rest <= 0 {
				rest = 0
				break
			}

			// calculate over the rest
			if rest > 0 {
				return calculate(rest, sizes, used)
			}
		}
	}

	// if there is a rest use the smallest size
	if rest > 0 {
		used[sizes[len(sizes)-1]] += 1
	}

	return used
}
