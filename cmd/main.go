package main

import (
	"log"

	"github.com/Gabriel-Schiestl/authgate/internal/src/module"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading env: %v", err)
	}

	app := fx.New(module.Module())

	app.Run()
}