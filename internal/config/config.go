package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port   string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

// LoadConfig загружает переменные из .env файла
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env файл не найден, используем переменные окружения")
	}

	return &Config{
		Port:   getEnv("PORT", "8080"),
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5432"),
		DBUser: getEnv("DB_USER", "effectivemobile"),
		DBPass: getEnv("DB_PASSWORD", ""),
		DBName: getEnv("DB_NAME", "effective_db"),
	}
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key string, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}
