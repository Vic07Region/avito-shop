package storage

import (
	"context"
	"database/sql" //nolint:gci
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (q *Queries) GetUser4UserID(ctx context.Context, userID uuid.UUID) (*Employee, error) {
	sqlquery := sq.Select("employee_id", "username", "email", "created_at").
		PlaceholderFormat(sq.Dollar).From("employees").Where(sq.Eq{"employee_id": userID})
	var user Employee
	err := sqlquery.RunWith(q.db).QueryRowContext(ctx).Scan(&user.EmployeeId, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		q.log.Error("GetUser4UserID QueryRowContext scan error:", zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (q *Queries) GetUserAuthData(ctx context.Context, username string) (*AuthData, error) {
	sqlBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sqlquery := sqlBuilder.Select("employee_id", "password_hash").
		From("employees").
		Where(sq.Expr("username ILIKE ?", username))

	var data AuthData

	err := sqlquery.RunWith(q.db).QueryRowContext(ctx).Scan(&data.UserID, &data.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		q.log.Error("GetUserAuthData QueryRowContext scan error:", zap.Error(err))
		return nil, err
	}
	return &data, nil
}

func (q *Queries) NewUser(ctx context.Context, username string, passwordHash string) (uuid.UUID, error) {
	txOptions := sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	}
	sqlBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	tx, err := q.db.BeginTx(ctx, &txOptions)
	if err != nil {
		q.log.Error("NewUser BeginTX error:", zap.Error(err))
		return uuid.Nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback() //nolint:errcheck
		}
	}()

	UserQuery := sqlBuilder.Insert("employees").
		Columns("username", "password_hash").
		Values(username, passwordHash).
		Suffix("RETURNING employee_id")

	var userID uuid.UUID

	err = UserQuery.RunWith(tx).QueryRowContext(ctx).Scan(&userID)
	if err != nil {
		q.log.Error("NewUser UserQuery error:", zap.Error(err))
		return uuid.Nil, err
	}

	WalletQuery := sqlBuilder.Insert("wallets").
		Columns("employee_id", "balance").
		Values(userID, 1000)

	result, err := WalletQuery.RunWith(tx).ExecContext(ctx)
	if err != nil {
		q.log.Error("NewUser WalletQuery error:", zap.Error(err))
		return uuid.Nil, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		q.log.Error("NewUser RowsAffected error:", zap.Error(err))
		return uuid.Nil, err
	}

	if affected < 1 {
		q.log.Error("NewUser affected rows:", zap.Int64("affected", affected))
		return uuid.Nil, err
	}

	err = tx.Commit()
	if err != nil {
		q.log.Error("NewUser Commit error:", zap.Error(err))
		return uuid.Nil, err
	}

	return userID, nil
}

func (q *Queries) FindUser(ctx context.Context, username string) (*Employee, error) {
	sqlBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sqlquery := sqlBuilder.Select("employee_id", "username", "created_at").
		From("employees").
		Where(sq.Expr("username ILIKE ?", username))

	var user Employee
	err := sqlquery.RunWith(q.db).
		QueryRowContext(ctx).
		Scan(&user.EmployeeId, &user.Name, &user.CreatedAt)
	if err != nil {
		q.log.Error("FindUser QueryRowContext scan error:", zap.Error(err))
		return nil, err
	}
	return &user, nil
}
