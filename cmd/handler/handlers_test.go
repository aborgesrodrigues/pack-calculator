package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pack-calculator/internal/common"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	savePackSize func(context.Context, *common.PackSizeBatch) error
	calculate    func(context.Context, *common.Order) (*common.Order, error)
}

func (m *mockService) SavePackSize(ctx context.Context, batch *common.PackSizeBatch) error {
	return m.savePackSize(ctx, batch)
}

func (m *mockService) Calculate(ctx context.Context, order *common.Order) (*common.Order, error) {
	return m.calculate(ctx, order)
}

func testHandler(svc *mockService) *Handler {
	return &Handler{
		logger:  slog.New(slog.DiscardHandler),
		service: svc,
	}
}

func TestSavePackSize(t *testing.T) {
	tests := map[string]struct {
		body           string
		mock           *mockService
		expectedStatus int
		expectedBody   string
	}{
		"success": {
			body: `{"sizes":[5000,2000]}`,
			mock: &mockService{
				savePackSize: func(_ context.Context, batch *common.PackSizeBatch) error {
					if len(batch.Sizes) != 2 {
						return errors.New("unexpected sizes length")
					}
					return nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Pack sizes saved",
		},
		"invalid JSON": {
			body:           `{invalid`,
			mock:           &mockService{},
			expectedStatus: http.StatusBadRequest,
		},
		"empty sizes": {
			body:           `{"sizes":[]}`,
			mock:           &mockService{},
			expectedStatus: http.StatusBadRequest,
		},
		"service error": {
			body: `{"sizes":[1000]}`,
			mock: &mockService{
				savePackSize: func(context.Context, *common.PackSizeBatch) error {
					return errors.New("db error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
		"decodes request body": {
			body: `{"sizes":[42]}`,
			mock: &mockService{
				savePackSize: func(_ context.Context, batch *common.PackSizeBatch) error {
					assert.Equal(t, []int{42}, batch.Sizes)
					return nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Pack sizes saved",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			h := testHandler(tt.mock)
			req := httptest.NewRequest(http.MethodPost, "/pack_size/batch", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			h.SavePackSize(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedBody != "" {
				var msg string
				err := json.Unmarshal(rec.Body.Bytes(), &msg)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, msg)
			}
		})
	}
}

func TestCalculate(t *testing.T) {
	orderID := uuid.New()

	tests := map[string]struct {
		body           string
		mock           *mockService
		expectedStatus int
		expectedOrder  *common.Order
	}{
		"success": {
			body: `{"items":12}`,
			mock: &mockService{
				calculate: func(_ context.Context, order *common.Order) (*common.Order, error) {
					order.ID = orderID
					return order, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedOrder:  &common.Order{ID: orderID, AmountItems: 12},
		},
		"invalid JSON": {
			body:           `not-json`,
			mock:           &mockService{},
			expectedStatus: http.StatusBadRequest,
		},
		"service error": {
			body: `{"items":1}`,
			mock: &mockService{
				calculate: func(context.Context, *common.Order) (*common.Order, error) {
					return nil, errors.New("calculation failed")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			h := testHandler(tt.mock)
			req := httptest.NewRequest(http.MethodPost, "/calculate", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			h.Calculate(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedOrder != nil {
				var got common.Order
				err := json.Unmarshal(rec.Body.Bytes(), &got)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOrder.ID, got.ID)
				assert.Equal(t, tt.expectedOrder.AmountItems, got.AmountItems)
			}
		})
	}
}

func TestWriteResponse(t *testing.T) {
	tests := map[string]struct {
		status  int
		payload any
		wantErr bool
		check   func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		"success": {
			status:  http.StatusCreated,
			payload: map[string]string{"status": "ok"},
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

				var payload map[string]string
				err := json.NewDecoder(rec.Body).Decode(&payload)
				assert.NoError(t, err)
				assert.Equal(t, "ok", payload["status"])
			},
		},
		"encode error": {
			status:  http.StatusOK,
			payload: make(chan int),
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			err := writeResponse(rec, tt.status, tt.payload)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.check != nil {
				tt.check(t, rec)
			}
		})
	}
}
