package main

import (
	"UserSVC/internal/app"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	app.StartService()
}
