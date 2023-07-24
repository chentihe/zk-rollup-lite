package db

import (
	"github.com/chentihe/zk-rollup-lite/operator/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDB(config *config.Postgres) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
