package connection

import (
	"log"
	"os"
	"strconv"

	"github.com/Gabriel-Schiestl/authgate/internal/src/config"
	"github.com/Gabriel-Schiestl/authgate/internal/src/infra/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupConfig() *gorm.DB {
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Error converting DB_PORT to int")
	}

	dbConfig := config.NewDbConfig(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), dbPort)

	db, err := gorm.Open(postgres.Open(dbConfig.ToString()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	db.AutoMigrate(entities.Auth{}, entities.UserInfo{})

	return db
}