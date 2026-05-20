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
	sizes, err := s.db.GetPackSizes(ctx)
	if err != nil {
		return nil, fmt.Errorf("error get pack sizes: %v", err)
	}

	order.ID = uuid.New()
	order.Packs = choosePacks(order.AmountItems, sizes)

	if err := s.db.SaveOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("error saving order: %v", err)
	}

	return order, nil
}

// choosePacks picks whole packs for ordered items following the challenge rules:
//  1. only whole packs; 2. minimise total items shipped; 3. minimise number of packs.
func choosePacks(ordered int, sizes []int) map[int]int {
	// Nothing to fulfil without a positive order or at least one pack size.
	if ordered <= 0 || len(sizes) == 0 {
		return map[int]int{}
	}

	// The optimal shipment never needs more than one extra pack beyond the ordered amount.
	maxSize := sizes[0]
	for _, size := range sizes[1:] {
		if size > maxSize {
			maxSize = size
		}
	}

	// minPacks[total] = fewest packs needed to ship exactly `total` items.
	// choice[total] = pack size added last when reaching `total`.
	const inf = int(^uint(0) >> 1)
	limit := ordered + maxSize
	minPacks := make([]int, limit+1)
	choice := make([]int, limit+1)
	for i := range minPacks {
		minPacks[i] = inf
	}
	minPacks[0] = 0

	// Unbounded knapsack: each pack size may be used as many times as needed.
	for total := 1; total <= limit; total++ {
		for _, size := range sizes {
			prev := total - size
			if prev < 0 || minPacks[prev] == inf {
				continue
			}
			if n := minPacks[prev] + 1; n < minPacks[total] {
				minPacks[total] = n
				choice[total] = size
			}
		}
	}

	// Rule 2: smallest total that still covers the order.
	bestTotal := -1
	for total := ordered; total <= limit; total++ {
		if minPacks[total] != inf {
			bestTotal = total
			break
		}
	}
	if bestTotal < 0 {
		return map[int]int{}
	}

	// Rule 3 is encoded in minPacks[bestTotal]; rebuild pack counts from `choice`.
	packs := make(map[int]int)
	for total := bestTotal; total > 0; {
		size := choice[total]
		packs[size]++
		total -= size
	}
	return packs
}
