package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	Logging     LoggingConfig
	Auth        AuthConfig
	EmailSender EmailSender
}

type ServerConfig struct {
	Port string
}
type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	DSN      string
}
type LoggingConfig struct {
	Level  string
	Format string
}

type AuthConfig struct {
	Secret string
}
type EmailSender struct {
	Email        string
	Password     string
	SmtpWithPort string
	Smtp         string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &Config{
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},
		Database: DatabaseConfig{
			Host:     os.Getenv("DATABASE_HOST"),
			Port:     os.Getenv("DATABASE_PORT"),
			Username: os.Getenv("DATABASE_USERNAME"),
			Password: os.Getenv("DATABASE_PASSWORD"),
			Database: os.Getenv("DATABASE_NAME"),
			DSN:      os.Getenv("DATABASE_DSN"),
		},
		Logging: LoggingConfig{
			Level:  os.Getenv("LOGGING_LEVEL"),
			Format: os.Getenv("LOGGING_FORMAT"),
		},
		Auth: AuthConfig{
			Secret: os.Getenv("KEY"),
		}, EmailSender: EmailSender{
			Email:        os.Getenv("EMAIL"),
			Password:     os.Getenv("PASSWORD"),
			SmtpWithPort: os.Getenv("SMTPWITHPORT"),
			Smtp:         os.Getenv("SMTP"),
		},
	}
}
