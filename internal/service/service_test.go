package service

import (
	"context"
	"github.com/Vic07Region/avito-shop/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetUser4UserID(ctx context.Context, userID uuid.UUID) (*storage.Employee, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*storage.Employee), args.Error(1)
}

func (m *MockStorage) FindUser(ctx context.Context, username string) (*storage.Employee, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*storage.Employee), args.Error(1)
}

func (m *MockStorage) GetUserAuthData(ctx context.Context, username string) (*storage.AuthData, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*storage.AuthData), args.Error(1)
}

func (m *MockStorage) NewUser(ctx context.Context, username string, passwordHash string) (uuid.UUID, error) {
	args := m.Called(ctx, username, passwordHash)
	if args.Get(0) == nil {
		return uuid.Nil, args.Error(1)
	}
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockStorage) GetBalance(ctx context.Context, userID uuid.UUID) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockStorage) GetInventories(ctx context.Context, userID uuid.UUID) ([]storage.InventoryItem, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]storage.InventoryItem), args.Error(1)
}

func (m *MockStorage) GetSendedCoins(ctx context.Context, userID uuid.UUID) ([]storage.SenderInfo, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]storage.SenderInfo), args.Error(1)
}

func (m *MockStorage) GetReceivedCoins(ctx context.Context, userID uuid.UUID) ([]storage.SenderInfo, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]storage.SenderInfo), args.Error(1)
}

func (m *MockStorage) SendCoinsTransaction(ctx context.Context, senderID uuid.UUID, receiverID uuid.UUID, amount int) error {
	args := m.Called(ctx, senderID, receiverID, amount)
	return args.Error(0)
}

func (m *MockStorage) PurchaseMerchTransaction(ctx context.Context, userID uuid.UUID, merch storage.MerchInfo) error {
	args := m.Called(ctx, userID, merch)
	return args.Error(0)
}

func (m *MockStorage) GetMerchItems(ctx context.Context, merchName string) (*storage.MerchItem, error) {
	args := m.Called(ctx, merchName)
	return args.Get(0).(*storage.MerchItem), args.Error(1)
}

// Инициализируем моковый логгер
func newTestLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}
