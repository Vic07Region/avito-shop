package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type BalanceResponse struct {
	Coins int `json:"coins"`
}

func AuthOrRegister(username string, password string) (string, error) {
	client := &http.Client{}

	reqBody := map[string]interface{}{
		"username": username,
		"password": password,
	}
	body, _ := json.Marshal(reqBody)
	resp, err := client.Post("http://localhost:8080/api/auth", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("Ошибка при авторизации/регистрации пользователя: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return "", fmt.Errorf("не удалось декодировать ошибку: %w", err)
		}
		return "", fmt.Errorf("ошибка авторизации: %s", errResp.Errors)
	} else {
		var okResp AuthResponse
		if err := json.NewDecoder(resp.Body).Decode(&okResp); err != nil {
			return "", fmt.Errorf("не удалось декодировать токен: %w", err)
		}
		return okResp.Token, nil
	}
}

func GetUserBalance(token string) (int, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:8080/api/info", nil)
	if err != nil {
		return 0, fmt.Errorf("не удалось сформировать запрос: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("не удалось выполнить запрос: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return 0, fmt.Errorf("не удалось декодировать ошибку: %w", err)
		}
		return 0, fmt.Errorf("ошибка получения баланса: %s", errResp.Errors)
	} else {
		var okResp BalanceResponse
		if err := json.NewDecoder(resp.Body).Decode(&okResp); err != nil {
			return 0, fmt.Errorf("не удалось декодировать баланс: %w", err)
		}
		return okResp.Coins, nil
	}
}

func BuyMerch(token string, merchName string) (int, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/buy/"+merchName, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return resp.StatusCode, fmt.Errorf("не удалось декодировать ошибку: %w", err)
		}
		log.Println(errResp.Errors)
		return resp.StatusCode, fmt.Errorf("%s", errResp.Errors)
	}
	return resp.StatusCode, err
}

func SendCoin(token string, username string, coins int) (int, error) {
	client := &http.Client{}
	reqBody := map[string]interface{}{
		"toUser": username,
		"amount": coins,
	}
	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "http://localhost:8080/api/sendCoin", bytes.NewBuffer(body))
	if err != nil {
		return 0, fmt.Errorf("не удалось сформировать запрос: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("не удалось выполнить запрос: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return resp.StatusCode, fmt.Errorf("не удалось декодировать ошибку: %w", err)
		}
		return resp.StatusCode, fmt.Errorf("%s", errResp.Errors)
	}
	return resp.StatusCode, nil
}
