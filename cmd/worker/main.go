package main

import (
	"fmt"

	"github.com/vnnyx/golang-dot-api/injector/wire"
)

func main() {
	controler, err := wire.InitializeController(".env")
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to kafka broker")
	controler.UserController.HandleMessage()
}
