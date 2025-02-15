package handlers

type Inventory struct {
	Type     string `json:"type" `
	Quantity int    `json:"quantity"`
}

type Received struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type Sent struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type CoinHistory struct {
	Received []Received `json:"received"`
	Sent     []Sent     `json:"sent"`
}

type FullInfo struct {
	Coins       int         `json:"coins" `
	Inventory   []Inventory `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}
