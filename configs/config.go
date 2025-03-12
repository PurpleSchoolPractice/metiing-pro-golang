package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	return &Config{
		Server: ServerConfig{
			Port: os.Getenv("PORT"),
		},
	}
}
