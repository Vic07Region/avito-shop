package service

import (
	"context"
	"errors"
	"github.com/Vic07Region/avito-shop/internal/storage"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type StorageInterface interface {
	GetUser4UserID(ctx context.Context, userID uuid.UUID) (*storage.Employee, error)
	GetUserAuthData(ctx context.Context, username string) (*storage.AuthData, error)
	NewUser(ctx context.Context, username string, passwordHash string) (uuid.UUID, error)
	GetBalance(ctx context.Context, userID uuid.UUID) (int, error)
	GetInventories(ctx context.Context, userID uuid.UUID) ([]storage.InventoryItem, error)
	GetReceivedCoins(ctx context.Context, userID uuid.UUID) ([]storage.SenderInfo, error)
	GetSendedCoins(ctx context.Context, userID uuid.UUID) ([]storage.SenderInfo, error)
	SendCoinsTransaction(ctx context.Context, senderID uuid.UUID, receiverID uuid.UUID, amount int) error
	PurchaseMerchTransaction(ctx context.Context, userID uuid.UUID, merch storage.MerchInfo) error
	FindUser(ctx context.Context, username string) (*storage.Employee, error)
	GetMerchItems(ctx context.Context, merchName string) (*storage.MerchItem, error)
}

var (
	ErrGenerateJWT          = errors.New("token generate error")
	ErrGeneratePasswordHash = errors.New("password hash generate error")
	ErrInvelidPassword      = errors.New("invalid password")
	ErrUserNotFound         = errors.New("user not found")
	ErrMerchNotFound        = errors.New("merch not found")
	ErrNotEnoughCoins       = errors.New("not enough coins on balance")
)

type Service struct {
	Storage StorageInterface
	log     *zap.Logger
}

func New(storage StorageInterface, zapLogger *zap.Logger) *Service {
	return &Service{
		Storage: storage,
		log:     zapLogger,
	}
}
