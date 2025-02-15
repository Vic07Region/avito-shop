package storage

import (
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

func (q *Queries) GetMerchItems(ctx context.Context, merchName string) (*MerchItem, error) {
	sqlBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sqlQuery := sqlBuilder.Select("item_id", "price").
		From("merch_items").
		Where(sq.Eq{"name": merchName})

	var merchItem MerchItem
	err := sqlQuery.RunWith(q.db).QueryRowContext(ctx).Scan(&merchItem.MerchID, &merchItem.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		q.log.Error("GetMerchItems QueryRowContext scan error:", zap.Error(err))
		return nil, err
	}
	merchItem.Name = merchName

	return &merchItem, nil
}
