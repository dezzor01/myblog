// internal/config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SiteTitle   string
	SiteTagline string // например: "Заметки о Go и жизни"
	FooterText  string

	AdminPassword string

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	ServerPort string
}

func Load() *Config {
	_ = godotenv.Load() // открываем .env

	cfg := &Config{
		SiteTitle:     mustEnv("SITE_TITLE"),
		SiteTagline:   mustEnv("SITE_TAGLINE"),
		FooterText:    mustEnv("FOOTER_TEXT"),
		AdminPassword: mustEnv("ADMIN_PASSWORD"),

		DBHost:     mustEnv("DB_HOST"),
		DBPort:     mustEnv("DB_PORT"),
		DBName:     mustEnv("DB_NAME"),
		DBUser:     mustEnv("DB_USER"),
		DBPassword: mustEnv("DB_PASSWORD"),
		DBSSLMode:  mustEnv("DB_SSLMODE"),

		ServerPort: mustEnv("SERVER_PORT"),
	}
	log.Printf("Конфиг загружен: %s (порт %s)", cfg.SiteTitle, cfg.ServerPort)
	return cfg
}

// проверяем наличие переменной
func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("ОБЯЗАТЕЛЬНАЯ переменная окружения не задана: %s", key)
	}
	return value
}
