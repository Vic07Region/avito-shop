package storage

import (
	"github.com/google/uuid"
	"time"
)

type Employee struct {
	EmployeeId uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      *string   `json:"email"`
	Password   string    `json:"password"`
	CreatedAt  time.Time `json:"created_at"`
}

type Wallet struct {
	EmployeeID uuid.UUID `json:"id"`
	Balance    int       `json:"balance"`
}

type InventoryItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type SenderInfo struct {
	Username string `json:"username"`
	Amount   int    `json:"amount"`
}

type MerchInfo struct {
	MerchID int `json:"merchID"`
	Price   int `json:"price"`
	Amount  int `json:"amount"`
}

type MerchItem struct {
	MerchID int    `json:"merchID"`
	Name    string `json:"name"`
	Price   int    `json:"price"`
}

type AuthData struct {
	UserID       uuid.UUID `json:"userID"`
	PasswordHash string    `json:"passwordHash"`
}
