package main

import (
	//nolint:gci
	"github.com/Vic07Region/avito-shop/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		err = godotenv.Load()
		if err != nil {
			return //nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	a, err := app.New()
	if err != nil {
		panic(err)
	}

	err = a.Run()
	if err != nil {
		panic(err)
	}
}
