package storage

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetBalance(t *testing.T) {
	ctx := context.Background()
	db, mock, queries := setupMockDB(t)
	defer db.Close()

	userID := uuid.New()
	expectedBalance := 100

	mock.ExpectQuery(`SELECT balance FROM wallets WHERE employee_id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(expectedBalance))

	balance, err := queries.GetBalance(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetInventories(t *testing.T) {
	ctx := context.Background()
	db, mock, queries := setupMockDB(t)
	defer db.Close()

	userID := uuid.New()

	mock.ExpectQuery(`SELECT name, SUM\(quantity\) as quantity FROM purchases INNER JOIN merch_items using\(item_id\) WHERE employee_id = \$1 GROUP BY name ORDER BY name`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"name", "quantity"}).
			AddRow("Item1", 2).
			AddRow("Item2", 5))

	inventories, err := queries.GetInventories(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, inventories, 2)
	assert.Equal(t, "Item1", inventories[0].Name)
	assert.Equal(t, 2, inventories[0].Quantity)
	assert.Equal(t, "Item2", inventories[1].Name)
	assert.Equal(t, 5, inventories[1].Quantity)

	assert.NoError(t, mock.ExpectationsWereMet()) // Проверка, что все ожидания выполнены
}

func TestGetReceivedCoins(t *testing.T) {
	ctx := context.Background()
	db, mock, queries := setupMockDB(t)
	defer db.Close()

	userID := uuid.New()
	mock.ExpectQuery(`SELECT username, sum\(amount\) FROM transactions INNER JOIN employees on employee_id = sender_id WHERE receiver_id = \$1 GROUP BY username ORDER BY username`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"username", "sum"}).
			AddRow("Alice", 100).
			AddRow("Bob", 200))

	receivedCoins, err := queries.GetReceivedCoins(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, receivedCoins, 2)
	assert.Equal(t, "Alice", receivedCoins[0].Username)
	assert.Equal(t, 100, receivedCoins[0].Amount)
	assert.Equal(t, "Bob", receivedCoins[1].Username)
	assert.Equal(t, 200, receivedCoins[1].Amount)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetReceivedCoins_NoRows(t *testing.T) {
	ctx := context.Background()
	db, mock, queries := setupMockDB(t)
	defer db.Close()

	userID := uuid.New()
	mock.ExpectQuery(`SELECT username, sum\(amount\) FROM transactions INNER JOIN employees on employee_id = sender_id WHERE receiver_id = \$1 GROUP BY username ORDER BY username`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"username", "sum"})) // Пустой результат

	receivedCoins, err := queries.GetReceivedCoins(ctx, userID)
	assert.NoError(t, err)
	assert.Empty(t, receivedCoins)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSendedCoins(t *testing.T) {
	ctx := context.Background()
	db, mock, queries := setupMockDB(t)
	defer db.Close()

	userID := uuid.New()
	mock.ExpectQuery(`SELECT username, sum\(amount\) FROM transactions INNER JOIN employees on employee_id = receiver_id WHERE sender_id = \$1 GROUP BY username ORDER BY username`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"username", "sum"}).
			AddRow("Charlie", 300).
			AddRow("David", 400))

	sendedCoins, err := queries.GetSendedCoins(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, sendedCoins, 2)
	assert.Equal(t, "Charlie", sendedCoins[0].Username)
	assert.Equal(t, 300, sendedCoins[0].Amount)
	assert.Equal(t, "David", sendedCoins[1].Username)
	assert.Equal(t, 400, sendedCoins[1].Amount)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSendedCoins_NoRows(t *testing.T) {
	ctx := context.Background()
	db, mock, queries := setupMockDB(t)
	defer db.Close()

	userID := uuid.New()
	mock.ExpectQuery(`SELECT username, sum\(amount\) FROM transactions INNER JOIN employees on employee_id = receiver_id WHERE sender_id = \$1 GROUP BY username ORDER BY username`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"username", "sum"})) // Пустой результат

	sendedCoins, err := queries.GetSendedCoins(ctx, userID)
	assert.NoError(t, err)
	assert.Empty(t, sendedCoins)

	assert.NoError(t, mock.ExpectationsWereMet())
}
