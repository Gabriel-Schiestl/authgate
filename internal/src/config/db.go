package config

import "fmt"

type DbConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

func NewDbConfig(host, user, password, dbname string, port int) *DbConfig {
	return &DbConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DbName:   dbname,
	}
}

func (db *DbConfig) ToString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", db.Host, db.Port, db.User, db.Password, db.DbName)
}