package db

import (
	"os"

	"github.com/chentihe/zk-rollup-lite/operator/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDB(config *config.Postgres) (*gorm.DB, error) {
	host := os.Getenv("POSTGRES_HOST")
	if host != "" {
		config.Host = host
	}

	db, err := gorm.Open(postgres.Open(config.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
