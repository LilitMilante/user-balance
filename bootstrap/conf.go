package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbName     string
	DbScheme   string
	DbLogin    string
	DbPassword string
	DbHost     string
	DbPort     string
	HttpPort   string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	c := Config{
		DbName:     os.Getenv("DB_NAME"),
		DbScheme:   os.Getenv("DB_SCHEME"),
		DbLogin:    os.Getenv("DB_LOGIN"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		HttpPort:   os.Getenv("HTTP_PORT"),
	}

	return &c, nil
}
