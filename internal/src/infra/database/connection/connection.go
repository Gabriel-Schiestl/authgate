package connection

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/Gabriel-Schiestl/authgate/internal/src/config"
	"github.com/Gabriel-Schiestl/authgate/internal/src/infra/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DbConfig *config.DbConfig
var Db *gorm.DB

func SetupConfig(host, user, password, port, name string) *sql.DB {
	dbPort, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Error converting DB_PORT to int")
	}

	DbConfig = config.NewDbConfig(host, user, password, name, dbPort)

	Db, err = gorm.Open(postgres.Open(DbConfig.ToString()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	sqlDb, err := Db.DB()
	if err != nil {
		log.Fatalf("Error getting DB connection: %v", err)
	}

	Db.AutoMigrate(entities.Auth{}, entities.UserInfo{})

	return sqlDb
}