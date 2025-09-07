package main

import (
	"log"

	"github.com/NathanBurns3/Focus-Tracker/cmd"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if present
    if err := godotenv.Load("../.env"); err != nil {
        log.Println("No .env file found or error loading .env file")
    }
	cmd.Execute()
}
