package main

import (
	"log"

	"github.com/faris-muhammed/e-commerce/order-service/config"
)

func main() {
	_, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("Failed to connect database")
	}
}
