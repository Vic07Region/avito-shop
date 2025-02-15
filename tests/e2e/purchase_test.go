package e2e

import (
	"net/http"
	"testing"

	"github.com/Vic07Region/avito-shop/tests"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
)

//merch name and prices
//name:		price
//t-shirt		80
//cup			20
//book			50
//pen			10
//powerbank		200
//hoody			300
//umbrella		200
//socks			10
//wallet		50
//pink-hoody	500

func TestPurchaseMerch(t *testing.T) {
	var token string
	// Регистрация пользователя
	t.Run("Авторизация/регистрация", func(t *testing.T) {
		username := faker.Username()
		userToken, err := tests.AuthOrRegister(username, "password_1")
		//проверка что пользователь создан или авторизован
		assert.NoError(t, err, "Не удалось авторизоваться/зарегистрироваться")
		assert.NotEmpty(t, userToken, "Токен неполучен")
		token = userToken
	})

	merchName := "powerbank"

	t.Run("Покупка мерча", func(t *testing.T) {
		status, err := tests.BuyMerch(token, merchName)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
	})

	// Проверяем баланс
	t.Run("проверяем баланс", func(t *testing.T) {
		balance, err := tests.GetUserBalance(token)
		assert.NoError(t, err)
		assert.Equal(t, 800, balance)
	})

}

func TestPurchaseMerchNoFoundMerch(t *testing.T) {
	var token string
	// Регистрация пользователя
	t.Run("Авторизация/регистрация", func(t *testing.T) {
		username := faker.Username()
		userToken, err := tests.AuthOrRegister(username, "password_1")
		//проверка что пользователь создан или авторизован
		assert.NoError(t, err, "Не удалось авторизоваться/зарегистрироваться")
		assert.NotEmpty(t, userToken, "Токен неполучен")
		token = userToken
	})

	merchName := "powerbank_avito"

	t.Run("Покупка мерча которого нет", func(t *testing.T) {
		status, err := tests.BuyMerch(token, merchName)
		assert.EqualError(t, err, "merch not found")
		assert.Equal(t, http.StatusBadRequest, status)
	})

	// Проверяем баланс
	t.Run("проверяем баланс", func(t *testing.T) {
		balance, err := tests.GetUserBalance(token)
		assert.NoError(t, err)
		assert.Equal(t, 1000, balance)
	})

}

func TestPurchaseMerchNoBalance(t *testing.T) {
	var token string
	// Регистрация пользователя
	t.Run("Авторизация/регистрация", func(t *testing.T) {
		username := faker.Username()
		userToken, err := tests.AuthOrRegister(username, "password_1")
		//проверка что пользователь создан или авторизован
		assert.NoError(t, err, "Не удалось авторизоваться/зарегистрироваться")
		assert.NotEmpty(t, userToken, "Токен неполучен")
		token = userToken
	})

	merchName := "hoody"
	// Покупка мерча 3 раз что бы опустошить баланс
	for i := 0; i < 3; i++ {
		tests.BuyMerch(token, merchName)
	}

	t.Run("Покупка мерча когда нехватает монет", func(t *testing.T) {
		status, err := tests.BuyMerch(token, merchName)
		assert.EqualError(t, err, "not enough coins on balance")
		assert.Equal(t, http.StatusBadRequest, status)
	})

	// Проверяем баланс
	t.Run("проверяем баланс", func(t *testing.T) {
		balance, err := tests.GetUserBalance(token)
		assert.NoError(t, err)
		assert.Equal(t, 100, balance)
	})

}
