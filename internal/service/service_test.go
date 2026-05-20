package service

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"pack-calculator/internal/common"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockDB struct {
	clearPackSize func(context.Context) error
	savePackSize  func(context.Context, *common.PackSizeBatch) error
	getPackSizes  func(context.Context) ([]int, error)
}

func (m *mockDB) ClearPackSize(ctx context.Context) error {
	return m.clearPackSize(ctx)
}

func (m *mockDB) SavePackSize(ctx context.Context, batch *common.PackSizeBatch) error {
	return m.savePackSize(ctx, batch)
}

func (m *mockDB) GetPackSizes(ctx context.Context) ([]int, error) {
	return m.getPackSizes(ctx)
}

func testService(db *mockDB) *Service {
	return &Service{
		logger: slog.New(slog.DiscardHandler),
		db:     db,
	}
}

func TestSavePackSize(t *testing.T) {
	var cleared, saved bool

	tests := map[string]struct {
		db      *mockDB
		batch   *common.PackSizeBatch
		wantErr bool
	}{
		"success": {
			db: &mockDB{
				clearPackSize: func(context.Context) error {
					cleared = true
					return nil
				},
				savePackSize: func(_ context.Context, batch *common.PackSizeBatch) error {
					saved = true
					if len(batch.Sizes) != 2 {
						return errors.New("unexpected batch size")
					}
					return nil
				},
			},
			batch: &common.PackSizeBatch{Sizes: []int{5000, 2000}},
		},
		"clear error": {
			db: &mockDB{
				clearPackSize: func(context.Context) error {
					return errors.New("clear failed")
				},
			},
			batch:   &common.PackSizeBatch{Sizes: []int{1000}},
			wantErr: true,
		},
		"save error": {
			db: &mockDB{
				clearPackSize: func(context.Context) error { return nil },
				savePackSize: func(context.Context, *common.PackSizeBatch) error {
					return errors.New("save failed")
				},
			},
			batch:   &common.PackSizeBatch{Sizes: []int{1000}},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cleared, saved = false, false

			err := testService(tt.db).SavePackSize(context.Background(), tt.batch)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.True(t, cleared)
			assert.True(t, saved)
		})
	}
}

func TestCalculate(t *testing.T) {
	tests := map[string]struct {
		order         *common.Order
		packSizes     []int
		packSizesErr  error
		wantErr       bool
		wantIDNotNil  bool
	}{
		"success": {
			order:        &common.Order{AmountItems: 12},
			packSizes:    []int{5000, 2000, 1000},
			wantIDNotNil: true,
		},
		"error": {
			order:        &common.Order{AmountItems: 1},
			packSizesErr: errors.New("db unavailable"),
			wantErr:      true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			svc := testService(&mockDB{
				getPackSizes: func(context.Context) ([]int, error) {
					return tt.packSizes, tt.packSizesErr
				},
			})

			result, err := svc.Calculate(context.Background(), tt.order)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.wantIDNotNil {
				assert.NotEqual(t, uuid.Nil, result.ID)
			}
		})
	}
}

func TestCalculatePacking(t *testing.T) {
	tests := map[string]struct {
		amount int
		sizes  []int
		want   map[int]int
	}{
		"exact fit with largest size": {
			amount: 10,
			sizes:  []int{5, 3, 1},
			want:   map[int]int{5: 2},
		},
		"uses smaller sizes for remainder": {
			amount: 12,
			sizes:  []int{5, 3, 1},
			want:   map[int]int{5: 2, 1: 2},
		},
		"remainder uses smallest pack": {
			amount: 6,
			sizes:  []int{5, 3, 1},
			want:   map[int]int{5: 1, 1: 1},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := calculate(tt.amount, tt.sizes, make(map[int]int))
			assert.Equal(t, tt.want, got)
		})
	}
}
