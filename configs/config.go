package configs

import (
	"os"
	"fmt"
	"github.com/joho/godotenv"
)

type Config struct {
	Db DbConfig
	Auth AuthConfig
}

type DbConfig struct {
	DSN string
}

type AuthConfig struct {
	Secret string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	return &Config{
		Db: DbConfig{
			DSN: os.Getenv("DSN"),
		},
		Auth: AuthConfig{
			Secret: os.Getenv("SECRET"),
		},
	}
}