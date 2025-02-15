package service

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/Vic07Region/avito-shop/internal/storage" //nolint:gci
)

func TestGetWalletInfo_Success(t *testing.T) {
	mockStorage := new(MockStorage)
	svc := Service{
		Storage: mockStorage,
		log:     newTestLogger(),
	}

	ctx := context.Background()
	userID := uuid.New()

	mockStorage.On("GetBalance", mock.Anything, userID).Return(1000, nil)
	mockStorage.On("GetInventories", mock.Anything, userID).Return([]storage.InventoryItem{}, nil)
	mockStorage.On("GetSendedCoins", mock.Anything, userID).Return([]storage.SenderInfo{}, nil)
	mockStorage.On("GetReceivedCoins", mock.Anything, userID).Return([]storage.SenderInfo{}, nil)

	info, err := svc.GetWalletInfo(ctx, userID)
	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, 1000, info.Coins)
}

func TestSendCoins_Success(t *testing.T) {
	mockStorage := new(MockStorage)
	svc := Service{Storage: mockStorage}
	ctx := context.Background()
	userID := uuid.New()
	toUserID := uuid.New()
	toUsername := "test_user"
	amount := 500

	mockStorage.On("FindUser", ctx, toUsername).Return(&storage.Employee{EmployeeId: toUserID}, nil)
	mockStorage.On("SendCoinsTransaction", ctx, userID, toUserID, amount).Return(nil)

	err := svc.SendCoins(ctx, userID, toUsername, amount)
	assert.NoError(t, err)
}

func TestSendCoins_UserNotFound(t *testing.T) {
	mockStorage := new(MockStorage)
	svc := Service{Storage: mockStorage}
	ctx := context.Background()
	userID := uuid.New()
	toUsername := "unknown_user"

	mockStorage.On("FindUser", ctx, toUsername).Return(&storage.Employee{}, sql.ErrNoRows)

	err := svc.SendCoins(ctx, userID, toUsername, 500)
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestPurchaseMerch_Success(t *testing.T) {
	mockStorage := new(MockStorage)
	svc := Service{Storage: mockStorage,
		log: newTestLogger(),
	}
	ctx := context.Background()
	userID := uuid.New()
	merchName := "Cool T-Shirt"
	quantity := 2

	mockStorage.On("GetMerchItems", ctx, merchName).Return(&storage.MerchItem{MerchID: 1, Price: 200}, nil)
	mockStorage.On("PurchaseMerchTransaction", ctx, userID, mock.Anything).Return(nil)

	err := svc.PurchaseMerch(ctx, userID, merchName, quantity)
	assert.NoError(t, err)
}

func TestPurchaseMerch_NotEnoughCoins(t *testing.T) {
	mockStorage := new(MockStorage)
	svc := Service{Storage: mockStorage,
		log: newTestLogger()}
	ctx := context.Background()
	userID := uuid.New()
	merchName := "Expensive Item"
	quantity := 5

	mockStorage.On("GetMerchItems", ctx, merchName).Return(&storage.MerchItem{MerchID: 2, Price: 1000}, nil)
	mockStorage.On("PurchaseMerchTransaction", ctx, userID, mock.Anything).Return(storage.ErrNotEnoughCoins)

	err := svc.PurchaseMerch(ctx, userID, merchName, quantity)
	assert.ErrorIs(t, err, ErrNotEnoughCoins)
}
