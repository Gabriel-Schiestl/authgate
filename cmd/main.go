package main

import (
	"log"
	"os"

	"github.com/Gabriel-Schiestl/authgate/internal/src/infra/database/connection"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading env: %v", err)
	}

	sqlDb := connection.SetupConfig(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	defer sqlDb.Close()
}