package db

import (
	"gorm.io/gorm"
	"6-project/configs"
	"gorm.io/driver/postgres"
)

type DB struct {
	*gorm.DB
}

func NewDb(config *configs.Config) *DB {
	db, err := gorm.Open(postgres.Open(config.Db.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &DB{db}
}