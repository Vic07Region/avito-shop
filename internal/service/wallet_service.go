package service

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/Vic07Region/avito-shop/internal/storage" //nolint:gci
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) GetWalletInfo(ctx context.Context, userID uuid.UUID) (*FullInfo, error) {
	var wg sync.WaitGroup
	var fullInfo FullInfo

	errChan := make(chan error, 4)
	wg.Add(4)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		defer wg.Done()
		coins, err := s.Storage.GetBalance(ctx, userID)
		if err != nil {
			s.log.Error("GetWalletInfo GetBalance error:", zap.Error(err))
			errChan <- err
			cancel()
		}
		fullInfo.Coins = coins
	}()

	go func() {
		defer wg.Done()
		inventories, err := s.Storage.GetInventories(ctx, userID)
		if err != nil {
			s.log.Error("GetWalletInfo GetInventories error:", zap.Error(err))
			errChan <- err
			cancel()
		}
		var invList []Inventory
		for _, i := range inventories {
			invList = append(invList, Inventory{
				Type:     i.Name,
				Quantity: i.Quantity,
			})
		}
		fullInfo.Inventory = invList
	}()

	go func() {
		defer wg.Done()
		sendedCoins, err := s.Storage.GetSendedCoins(ctx, userID)
		if err != nil {
			s.log.Error("GetWalletInfo GetSendedCoins error:", zap.Error(err))
			errChan <- err
			cancel()
		}
		var sendedList []Sent
		for _, i := range sendedCoins {
			sendedList = append(sendedList, Sent{
				ToUser: i.Username,
				Amount: i.Amount,
			})
		}
		fullInfo.CoinHistory.Sent = sendedList
	}()

	go func() {
		defer wg.Done()
		receivedCoins, err := s.Storage.GetReceivedCoins(ctx, userID)
		if err != nil {
			s.log.Error("GetWalletInfo GetReceivedCoins error:", zap.Error(err))
			errChan <- err
			cancel()
		}
		var receivedList []Received
		for _, i := range receivedCoins {
			receivedList = append(receivedList, Received{
				FromUser: i.Username,
				Amount:   i.Amount,
			})
		}
		fullInfo.CoinHistory.Received = receivedList
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return &fullInfo, nil
}

func (s *Service) SendCoins(ctx context.Context, userID uuid.UUID, toUsername string, amount int) error {
	user, err := s.Storage.FindUser(ctx, toUsername)
	if err != nil {

		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}
		s.log.Error("SendCoins FindUser error:", zap.Error(err))
		return err
	}

	err = s.Storage.SendCoinsTransaction(ctx, userID, user.EmployeeId, amount)
	if err != nil {
		if errors.Is(err, storage.ErrNotEnoughCoins) {
			return ErrNotEnoughCoins
		}
		s.log.Error("SendCoins SendCoinsTransaction error:", zap.Error(err))
		return err
	}

	return nil
}

func (s *Service) PurchaseMerch(ctx context.Context, userID uuid.UUID, merchName string, quantity int) error {
	merch, err := s.Storage.GetMerchItems(ctx, merchName)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrMerchNotFound
		}
		s.log.Error("PurchaseMerch GetMerchItems error:", zap.Error(err))
		return err
	}

	err = s.Storage.PurchaseMerchTransaction(ctx, userID, storage.MerchInfo{
		MerchID: merch.MerchID,
		Price:   merch.Price,
		Amount:  quantity,
	})
	if err != nil {
		if errors.Is(err, storage.ErrNotEnoughCoins) {
			return ErrNotEnoughCoins
		}
		s.log.Error("PurchaseMerch PurchaseMerchTransaction error:", zap.Error(err))
		return err
	}

	return nil
}
