package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMerchItems(t *testing.T) {
	db, mock, queries := setupMockDB(t)
	defer db.Close()

	ctx := context.Background()
	merchName := "t-shirt"
	fakeMerchName := "fake-merch"
	merchID := 1
	merchPrice := 80

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT item_id, price FROM merch_items WHERE name = \$1`).
			WithArgs(merchName).
			WillReturnRows(sqlmock.NewRows([]string{"item_id", "price"}).AddRow(merchID, merchPrice))

		merchItem, err := queries.GetMerchItems(ctx, merchName)
		assert.NoError(t, err)
		assert.NotNil(t, merchItem)
		assert.Equal(t, merchID, merchItem.MerchID)
		assert.Equal(t, merchPrice, merchItem.Price)
	})

	t.Run("NoRows", func(t *testing.T) {
		mock.ExpectQuery(`SELECT item_id, price FROM merch_items WHERE name = \$1`).
			WithArgs(fakeMerchName).
			WillReturnError(sql.ErrNoRows)

		merchItem, err := queries.GetMerchItems(ctx, fakeMerchName)
		assert.Error(t, err)
		assert.Nil(t, merchItem)
		assert.True(t, errors.Is(err, sql.ErrNoRows))
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}
