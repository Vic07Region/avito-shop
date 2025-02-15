package storage

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"time"
)

const (
	DBPostgres = "postgres"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrNotEnoughCoins = errors.New("not enough coins on balance")
)

type Queries struct {
	db  *sql.DB
	log *zap.Logger
}

type ConnectionParams struct {
	DbDriver         string
	ConnectionString string
	MaxOpenConns     int
	MsxIdleConns     int
	MaxLifeTime      time.Duration
}

func NewDBConection(params ConnectionParams) (*sql.DB, error) {
	db, err := sql.Open(params.DbDriver, params.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	db.SetMaxOpenConns(params.MaxOpenConns)
	db.SetMaxIdleConns(params.MsxIdleConns)
	db.SetConnMaxLifetime(params.MaxLifeTime)
	if err = db.Ping(); err != nil {
		err := db.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
	return db, nil
}

func New(db *sql.DB, zapLogger *zap.Logger) *Queries {
	return &Queries{db, zapLogger}
}
