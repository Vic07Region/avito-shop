package service

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/Vic07Region/avito-shop/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginUser(t *testing.T) {
	mockStorage := new(MockStorage)
	svc := Service{
		Storage: mockStorage,
		log:     newTestLogger(),
	}

	ctx := context.Background()
	os.Setenv("SECRET_KEY", "secret")
	testUsername := "test_user"
	testPassword := "securepassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	userID := uuid.New()

	t.Run("User not found - new user created", func(t *testing.T) {
		mockStorage.On("GetUserAuthData", ctx, testUsername).
			Return((*storage.AuthData)(nil), storage.ErrUserNotFound)
		mockStorage.On("NewUser", ctx, testUsername, mock.Anything).Return(userID, nil)

		token, err := svc.LoginUser(ctx, UserData{Username: testUsername, Password: testPassword})
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		mockStorage.AssertExpectations(t)
	})

	t.Run("Error generating password hash", func(t *testing.T) {
		mockStorage.ExpectedCalls = nil
		mockStorage.On("GetUserAuthData", ctx, testUsername).
			Return((*storage.AuthData)(nil), storage.ErrUserNotFound)
		// Ломаем bcrypt, передавая слишком длинный пароль
		longPassword := string(make([]byte, 100_000))
		_, err := svc.LoginUser(ctx, UserData{Username: testUsername, Password: longPassword})
		assert.Error(t, err)
		assert.EqualError(t, err, ErrGeneratePasswordHash.Error())
	})

	t.Run("Error creating new user", func(t *testing.T) {
		mockStorage.ExpectedCalls = nil
		mockStorage.On("GetUserAuthData", ctx, testUsername).
			Return((*storage.AuthData)(nil), storage.ErrUserNotFound)
		mockStorage.On("NewUser", ctx, testUsername, mock.Anything).
			Return(uuid.UUID{}, errors.New("db error"))
		_, err := svc.LoginUser(ctx, UserData{Username: testUsername, Password: testPassword})
		log.Println(err)
		assert.Error(t, err)
	})

	t.Run("User found - successful login", func(t *testing.T) {
		mockStorage.ExpectedCalls = nil
		mockStorage.On("GetUserAuthData", ctx, testUsername).Return(&storage.AuthData{UserID: userID, PasswordHash: string(hashedPassword)}, nil)
		token, err := svc.LoginUser(ctx, UserData{Username: testUsername, Password: testPassword})
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		mockStorage.AssertExpectations(t)
	})

	t.Run("Incorrect password", func(t *testing.T) {
		mockStorage.ExpectedCalls = nil
		mockStorage.On("GetUserAuthData", ctx, testUsername).Return(&storage.AuthData{UserID: userID, PasswordHash: string(hashedPassword)}, nil)

		_, err := svc.LoginUser(ctx, UserData{Username: testUsername, Password: "wrongpassword"})
		assert.Error(t, err)
		assert.Equal(t, ErrInvelidPassword, err)
	})

}
