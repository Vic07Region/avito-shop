package storage

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *Queries) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	logger, _ := zap.NewDevelopment()
	queries := &Queries{db: db, log: logger}

	return db, mock, queries
}
