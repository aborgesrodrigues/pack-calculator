package db

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"testing"

	"pack-calculator/internal/common"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/stretchr/testify/assert"
)

func testDB(t *testing.T, sqlDB *sql.DB) *DB {
	t.Helper()
	return &DB{
		logger: slog.New(slog.DiscardHandler),
		db:     sqlDB,
	}
}

func TestGetPackSizes(t *testing.T) {
	const getPackSizesQuery = `SELECT COALESCE\(array_agg`

	tests := map[string]struct {
		rows     *sqlmock.Rows
		expected []int
		err      error
	}{
		"success": {
			rows:     sqlmock.NewRows([]string{"array_agg"}).AddRow("{5000,2000,1000}"),
			expected: []int{5000, 2000, 1000},
		},
		"empty": {
			rows:     sqlmock.NewRows([]string{"array_agg"}).AddRow("{}"),
			expected: []int{},
		},
		"error": {
			err: errors.New("some error"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			assert.NoError(t, err)

			t.Cleanup(func() { _ = sqlDB.Close() })

			query := mock.ExpectQuery(getPackSizesQuery)
			if tt.rows != nil {
				query.WillReturnRows(tt.rows)
			} else {
				query.WillReturnError(tt.err)
			}

			sizes, err := testDB(t, sqlDB).GetPackSizes(context.Background())
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.expected, sizes)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestClearPackSize(t *testing.T) {
	tests := map[string]struct {
		expected error
		err      error
	}{
		"success": {},
		"error": {
			expected: errors.New("some error"),
			err:      errors.New("some error"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			assert.NoError(t, err)

			t.Cleanup(func() { _ = sqlDB.Close() })

			if tt.err == nil {
				mock.ExpectExec(`DELETE FROM pack_size`).WillReturnResult(sqlmock.NewResult(0, 3))
			} else {
				mock.ExpectExec(`DELETE FROM pack_size`).WillReturnError(tt.err)
			}

			err = testDB(t, sqlDB).ClearPackSize(context.Background())
			assert.Equal(t, tt.expected, err)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestSavePackSize(t *testing.T) {
	tests := map[string]struct {
		batch    *common.PackSizeBatch
		expected error
		err      error
	}{
		"success": {
			batch: &common.PackSizeBatch{Sizes: []int{5000, 2000, 1000}},
		},
		"error": {
			batch:    &common.PackSizeBatch{Sizes: []int{5000, 2000, 1000}},
			expected: errors.New("some error"),
			err:      errors.New("some error"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			assert.NoError(t, err)

			t.Cleanup(func() { _ = sqlDB.Close() })

			if tt.err == nil {
				mock.ExpectExec(`INSERT INTO public.pack_size`).
					WithArgs(sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 3))
			} else {
				mock.ExpectExec(`INSERT INTO public.pack_size`).
					WithArgs(sqlmock.AnyArg()).
					WillReturnError(tt.err)
			}

			err = testDB(t, sqlDB).SavePackSize(context.Background(), tt.batch)
			assert.Equal(t, tt.expected, err)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestNewDB(t *testing.T) {
	tests := map[string]struct {
		pingErr error
		wantErr bool
	}{
		"success": {},
		"database down": {
			pingErr: errors.New("connection refused"),
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			assert.NoError(t, err)

			t.Cleanup(func() { _ = sqlDB.Close() })

			ping := mock.ExpectPing()
			if tt.pingErr != nil {
				ping.WillReturnError(tt.pingErr)
			}

			db, err := newDBWithConn(slog.New(slog.DiscardHandler), sqlDB)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
