package main

import (
	//nolint:gci
	"github.com/Vic07Region/avito-shop/internal/app"
)

func main() {

	a, err := app.New()
	if err != nil {
		panic(err)
	}

	err = a.Run()
	if err != nil {
		panic(err)
	}
}
