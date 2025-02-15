package storage

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUser4UserID(t *testing.T) {
	ctx := context.Background()

	db, mock, queries := setupMockDB(t)
	defer db.Close()

	userID := uuid.New()
	createdAt := time.Now() // Генерируем текущее время

	rows := sqlmock.NewRows([]string{"employee_id", "username", "email", "created_at"}).
		AddRow(userID, "testuser", "test@example.com", createdAt) // Используем time.Time

	mock.ExpectQuery(`SELECT employee_id, username, email, created_at FROM employees WHERE employee_id = \$1`).
		WithArgs(userID).
		WillReturnRows(rows)

	user, err := queries.GetUser4UserID(ctx, userID)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.EmployeeId)
	assert.Equal(t, "testuser", user.Name)
	assert.Equal(t, "test@example.com", *user.Email)
	assert.WithinDuration(t, createdAt, user.CreatedAt, time.Second) // Проверяем с учетом возможных миллисекундных отклонений
}

func TestGetUserAuthData(t *testing.T) {
	ctx := context.Background()

	db, mock, queries := setupMockDB(t)
	defer db.Close()

	username := "testuser"
	userID := uuid.New()
	passwordHash := "hashedpassword"

	rows := sqlmock.NewRows([]string{"employee_id", "password_hash"}).
		AddRow(userID, passwordHash)

	mock.ExpectQuery(`SELECT employee_id, password_hash FROM employees WHERE username ILIKE \$1`).
		WithArgs(username).
		WillReturnRows(rows)

	data, err := queries.GetUserAuthData(ctx, username)
	require.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, userID, data.UserID)
	assert.Equal(t, passwordHash, data.PasswordHash)
}

func TestNewUser(t *testing.T) {
	ctx := context.Background()

	db, mock, queries := setupMockDB(t)
	defer db.Close()

	username := "newuser"
	passwordHash := "newhash"
	newUserID := uuid.New()

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO employees (username,password_hash) VALUES ($1,$2) RETURNING employee_id`)).
		WithArgs(username, passwordHash).
		WillReturnRows(sqlmock.NewRows([]string{"employee_id"}).AddRow(newUserID))

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO wallets (employee_id,balance) VALUES ($1,$2)`)).
		WithArgs(newUserID, 1000).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	userID, err := queries.NewUser(ctx, username, passwordHash)
	require.NoError(t, err, "Ошибка при создании пользователя")
	require.NotEqual(t, uuid.Nil, userID, "userID не должен быть nil")
	assert.Equal(t, newUserID, userID)

	require.NoError(t, mock.ExpectationsWereMet(), "Не все ожидания были выполнены")
}

func TestFindUser(t *testing.T) {
	ctx := context.Background()

	db, mock, queries := setupMockDB(t)
	defer db.Close()

	username := "existinguser"
	userID := uuid.New()
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{"employee_id", "username", "created_at"}).
		AddRow(userID, username, createdAt)

	mock.ExpectQuery(`SELECT employee_id, username, created_at FROM employees WHERE username ILIKE \$1`).
		WithArgs(username).
		WillReturnRows(rows)

	user, err := queries.FindUser(ctx, username)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.EmployeeId)
	assert.Equal(t, username, user.Name)
	assert.WithinDuration(t, createdAt, user.CreatedAt, time.Second) // Проверяем дату с небольшой погрешностью
}
