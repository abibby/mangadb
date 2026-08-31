package config

import (
	"errors"
	"os"

	"abibby.com/salusa/database"
	"abibby.com/salusa/database/dialects/postgres"
	"abibby.com/salusa/email"
	"abibby.com/salusa/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Port     int
	BasePath string

	Database database.Config
	Mail     email.Config
}

func Load() *Config {
	err := godotenv.Load("./.env")
	if errors.Is(err, os.ErrNotExist) {
		// fall through
	} else if err != nil {
		panic(err)
	}

	return &Config{
		Port:     env.Int("PORT", 2303),
		BasePath: env.String("BASE_PATH", ""),
		Database: &postgres.Config{
			Username:   env.String("DATABASE_USERNAME", "icbmdb"),
			Password:   env.String("DATABASE_PASSWORD", "icbmdb"),
			Host:       env.String("DATABASE_HOST", "localhost"),
			Database:   env.String("DATABASE_NAME", "icbmdb"),
			DisableSSL: true,
		},
		Mail: &email.SMTPConfig{
			From:     env.String("MAIL_FROM", "salusa@example.com"),
			Host:     env.String("MAIL_HOST", "sandbox.smtp.mailtrap.io"),
			Port:     env.Int("MAIL_PORT", 2525),
			Username: env.String("MAIL_USERNAME", "user"),
			Password: env.String("MAIL_PASSWORD", "pass"),
		},
	}
}

func (c *Config) GetHTTPPort() int {
	return c.Port
}
func (c *Config) GetBaseURL() string {
	return c.BasePath
}
