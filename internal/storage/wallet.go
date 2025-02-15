package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq" //nolint:gci
	"go.uber.org/zap"
	"sync"

	sq "github.com/Masterminds/squirrel" //nolint:gci
	"github.com/google/uuid"
)

func (q *Queries) GetBalance(ctx context.Context, userID uuid.UUID) (int, error) {
	sqlBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sqlQuery := sqlBuilder.Select("balance").
		From("wallets").
		Where(sq.Eq{"employee_id": userID})
	var balance int
	if err := sqlQuery.RunWith(q.db).QueryRowContext(ctx).Scan(&balance); err != nil {
		q.log.Error("GetBalance QueryRowContext error:", zap.Error(err))
		return 0, err
	}
	return balance, nil
}

func (q *Queries) GetInventories(ctx context.Context, userID uuid.UUID) ([]InventoryItem, error) {
	sqlBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sqlQuery := sqlBuilder.Select("name", "SUM(quantity) as quantity").
		From("purchases").
		InnerJoin("merch_items using(item_id)").
		Where(sq.Eq{"employee_id": userID}).
		GroupBy("name").OrderBy("name")
	rows, err := sqlQuery.RunWith(q.db).QueryContext(ctx)
	if err != nil {
		q.log.Error("GetInventories QueryContext error:", zap.Error(err))
		return nil, err
	}

	defer rows.Close()

	var inventoryList []InventoryItem

	for rows.Next() {
		var i InventoryItem
		if err := rows.Scan(&i.Name, &i.Quantity); err != nil {
			q.log.Error("GetInventories rows.Scan error:", zap.Error(err))
			return nil, err
		}
		inventoryList = append(inventoryList, i)
	}

	if err = rows.Err(); err != nil {
		q.log.Error("GetInventories rows error:", zap.Error(err))
		return nil, err
	}

	return inventoryList, nil
}

func (q *Queries) GetReceivedCoins(ctx context.Context, userID uuid.UUID) ([]SenderInfo, error) {
	sqlBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sqlQuery := sqlBuilder.Select("username", "sum(amount)").
		From("transactions").
		InnerJoin("employees on employee_id = sender_id").
		Where(sq.Eq{"receiver_id": userID}).
		GroupBy("username").
		OrderBy("username")

	rows, err := sqlQuery.RunWith(q.db).QueryContext(ctx)
	if err != nil {
		q.log.Error("GetSendedCoins QueryContext error:", zap.Error(err))
		return nil, err
	}

	defer rows.Close()

	var senderInfoList []SenderInfo
	for rows.Next() {
		var i SenderInfo
		if err := rows.Scan(&i.Username, &i.Amount); err != nil {
			q.log.Error("GetSendedCoins rows.Scan error:", zap.Error(err))
			return nil, err
		}
		senderInfoList = append(senderInfoList, i)
	}
	if err = rows.Err(); err != nil {
		q.log.Error("GetReceivedCoins rows error:", zap.Error(err))
		return nil, err
	}

	return senderInfoList, nil
}

func (q *Queries) GetSendedCoins(ctx context.Context, userID uuid.UUID) ([]SenderInfo, error) {
	sqlBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sqlQuery := sqlBuilder.Select("username", "sum(amount)").
		From("transactions").
		InnerJoin("employees on employee_id = receiver_id").
		Where(sq.Eq{"sender_id": userID}).
		GroupBy("username").
		OrderBy("username")

	rows, err := sqlQuery.RunWith(q.db).QueryContext(ctx)
	if err != nil {
		q.log.Error("GetSendedCoins QueryContext error:", zap.Error(err))
		return nil, err
	}

	defer rows.Close()

	var senderInfoList []SenderInfo
	for rows.Next() {
		var i SenderInfo
		if err := rows.Scan(&i.Username, &i.Amount); err != nil {
			q.log.Error("GetSendedCoins rows.Scan error:", zap.Error(err))
			return nil, err
		}
		senderInfoList = append(senderInfoList, i)
	}
	if err = rows.Err(); err != nil {
		q.log.Error("GetSendedCoins rows error:", zap.Error(err))
		return nil, err
	}

	return senderInfoList, nil
}

func (q *Queries) SendCoinsTransaction(ctx context.Context, senderID uuid.UUID, receiverID uuid.UUID, amount int) error {
	sqlBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	txOptions := sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 3)
	wg.Add(3)

	tx, err := q.db.BeginTx(ctx, &txOptions)
	if err != nil {
		q.log.Error("SendCoins BeginTX error:", zap.Error(err))
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	SenderBalanceQuery := sqlBuilder.Update("wallets").
		Set("balance", sq.Expr("balance - ?", amount)).
		Where(sq.Eq{"employee_id": senderID})

	ReceiverBalanceQuery := sqlBuilder.Update("wallets").
		Set("balance", sq.Expr("balance + ?", amount)).
		Where(sq.Eq{"employee_id": receiverID})
	TransactionQuery := sqlBuilder.Insert("transactions").
		Columns("sender_id", "receiver_id", "amount").
		Values(senderID, receiverID, amount)

	var pgErr *pq.Error

	go func() {
		defer wg.Done()
		if _, err := SenderBalanceQuery.RunWith(tx).ExecContext(ctx); err != nil {
			if errors.As(err, &pgErr) {
				if pgErr.Code == "23514" {
					q.log.Error("SendCoins SenderBalanceQuery error:", zap.Error(ErrNotEnoughCoins))
					errChan <- ErrNotEnoughCoins
					cancel()
					return
				}
				if pgErr.Code == "25P02" {
					q.log.Error("SendCoins SenderBalanceQuery error:", zap.Error(err))
					cancel()
					return
				}

			}
			q.log.Error("SendCoins SenderBalanceQuery error:", zap.Error(err))
			errChan <- err
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := ReceiverBalanceQuery.RunWith(tx).ExecContext(ctx); err != nil {
			if errors.As(err, &pgErr) {
				if pgErr.Code == "25P02" {
					q.log.Error("SendCoins SenderBalanceQuery error:", zap.Error(err))
					cancel()
					return
				}

			}
			q.log.Error("SendCoins ReceiverBalanceQuery error:", zap.Error(err))
			errChan <- err
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := TransactionQuery.RunWith(tx).ExecContext(ctx); err != nil {
			if errors.As(err, &pgErr) {
				if pgErr.Code == "25P02" {
					q.log.Error("SendCoins SenderBalanceQuery error:", zap.Error(err))
					cancel()
					return
				}

			}
			q.log.Error("SendCoins TransactionQuery error:", zap.Error(err))
			errChan <- err
			cancel()
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		q.log.Error("SendCoins Commit error:", zap.Error(err))
		return err
	}
	return nil
}

func (q *Queries) PurchaseMerchTransaction(ctx context.Context, userID uuid.UUID, merch MerchInfo) error {
	sqlBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	txOptions := sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}

	var pgErr *pq.Error

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tx, err := q.db.BeginTx(ctx, &txOptions)
	if err != nil {
		q.log.Error("SendCoins BeginTX error:", zap.Error(err))
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	total := merch.Price * merch.Amount

	buyerBalanceQuery := sqlBuilder.Update("wallets").
		Set("balance", sq.Expr("balance - ?", total)).
		Where(sq.Eq{"employee_id": userID})

	purchaseQuery := sqlBuilder.Insert("purchases").
		Columns("employee_id", "item_id", "quantity").
		Values(userID, merch.MerchID, merch.Amount)

	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	wg.Add(2)

	go func() {
		defer wg.Done()
		if _, err := buyerBalanceQuery.RunWith(tx).ExecContext(ctx); err != nil {
			if errors.As(err, &pgErr) {
				if pgErr.Code == "23514" {
					q.log.Error("SendCoins buyerBalanceQuery error:", zap.Error(ErrNotEnoughCoins))
					errChan <- ErrNotEnoughCoins
					cancel()
					return
				}
				if pgErr.Code == "25P02" {
					q.log.Error("PurchaseMerch buyerBalanceQuery error:", zap.Error(err))
					cancel()
					return
				}
			}
			q.log.Error("PurchaseMerch buyerBalanceQuery error:", zap.Error(err))
			errChan <- err
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := purchaseQuery.RunWith(tx).ExecContext(ctx); err != nil {
			if errors.As(err, &pgErr) {
				if pgErr.Code == "25P02" {
					q.log.Error("PurchaseMerch purchaseQuery error:", zap.Error(err))
					cancel()
					return
				}
			}
			q.log.Error("PurchaseMerch purchaseQuery error:", zap.Error(err))
			errChan <- err
			cancel()
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		q.log.Error("PurchaseMerch Commit error:", zap.Error(err))
		return err
	}
	return nil
}
