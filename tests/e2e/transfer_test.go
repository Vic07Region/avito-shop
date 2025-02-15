package e2e

import (
	"github.com/Vic07Region/avito-shop/tests"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
)

func TestTransferCoins(t *testing.T) {
	var user_token1 string
	var username1 string
	// Регистрация пользователя
	t.Run("Авторизация/регистрация", func(t *testing.T) {
		username := faker.Username()
		userToken, err := tests.AuthOrRegister(username, "password_1")
		assert.NoError(t, err, "Не удалось авторизоваться/зарегистрироваться")
		assert.NotEmpty(t, userToken, "Токен неполучен")
		user_token1 = userToken
		username1 = username
	})
	//создаем пользователя для проверки
	username2 := faker.Username()
	userToken2, err := tests.AuthOrRegister(username2, "password_1")
	assert.NoError(t, err, "Не удалось авторизоваться/зарегистрироваться")
	assert.NotEmpty(t, userToken2, "Токен неполучен")

	t.Run("отправка монет другому пользователю:"+username2, func(t *testing.T) {
		status, err := tests.SendCoin(user_token1, username2, 100)
		log.Println(err)
		assert.NoError(t, err, "Монеты не переданы")
		assert.Equal(t, http.StatusOK, status)
	})

	// Проверяем баланс 1 пользователя
	t.Run("проверяем баланс пользователя:"+username1, func(t *testing.T) {
		balance, err := tests.GetUserBalance(user_token1)
		assert.NoError(t, err)
		assert.Equal(t, 900, balance)
	})
	// Проверяем баланс 2 пользователя
	t.Run("проверяем баланс пользователя:"+username2, func(t *testing.T) {
		balance, err := tests.GetUserBalance(userToken2)
		assert.NoError(t, err)
		assert.Equal(t, 1100, balance)
	})
}

func TestTransferCoinsNoUser(t *testing.T) {
	var user_token1 string
	var username1 string
	// Регистрация пользователя
	t.Run("Авторизация/регистрация", func(t *testing.T) {
		username := faker.Username()
		userToken, err := tests.AuthOrRegister(username, "password_1")
		assert.NoError(t, err, "Не удалось авторизоваться/зарегистрироваться")
		assert.NotEmpty(t, userToken, "Токен неполучен")
		user_token1 = userToken
		username1 = username
	})
	//создаем пользователя для проверки
	username2 := faker.Username()

	t.Run("отправка монет другому пользователю(не существует):"+username2, func(t *testing.T) {
		status, err := tests.SendCoin(user_token1, username2, 100)
		log.Println(err)
		assert.EqualError(t, err, "user not found")
		assert.Equal(t, http.StatusBadRequest, status)
	})

	// Проверяем баланс 1 пользователя
	t.Run("проверяем баланс пользователя:"+username1, func(t *testing.T) {
		balance, err := tests.GetUserBalance(user_token1)
		assert.NoError(t, err)
		assert.Equal(t, 1000, balance)
	})

}

func TestTransferCoinsNotEnoughCoins(t *testing.T) {
	var user_token1 string
	var username1 string
	t.Run("Авторизация/регистрация", func(t *testing.T) {
		username := faker.Username()
		userToken, err := tests.AuthOrRegister(username, "password_1")
		assert.NoError(t, err, "Не удалось авторизоваться/зарегистрироваться")
		assert.NotEmpty(t, userToken, "Токен неполучен")
		user_token1 = userToken
		username1 = username
	})

	username2 := faker.Username()
	userToken2, err := tests.AuthOrRegister(username2, "password_1")
	assert.NoError(t, err, "Не удалось авторизоваться/зарегистрироваться")
	assert.NotEmpty(t, userToken2, "Токен неполучен")

	t.Run("отправка монет другому пользователю:"+username2, func(t *testing.T) {
		status, err := tests.SendCoin(user_token1, username2, 2000)
		log.Println(err)
		assert.EqualError(t, err, "not enough coins on balance")
		assert.Equal(t, http.StatusBadRequest, status)
	})

	// Проверяем баланс 1 пользователя
	t.Run("проверяем баланс пользователя:"+username1, func(t *testing.T) {
		balance, err := tests.GetUserBalance(user_token1)
		assert.NoError(t, err)
		assert.Equal(t, 1000, balance)
	})
	// Проверяем баланс 2 пользователя
	t.Run("проверяем баланс пользователя:"+username2, func(t *testing.T) {
		balance, err := tests.GetUserBalance(userToken2)
		assert.NoError(t, err)
		assert.Equal(t, 1000, balance)
	})
}
