package main

import (
	"os"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"6-project/internal/auth"
	"6-project/internal/order"
	"github.com/joho/godotenv"
	"6-project/internal/product"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), (&gorm.Config{}))
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&auth.AuthByPhone{}, &product.Product{}, &order.Order{})
}