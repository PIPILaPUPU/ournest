package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
}

func Load() (Config, error) {
	// .env — удобство для локальной разработки; отсутствие файла не ошибка
	_ = godotenv.Load()

	if url := os.Getenv("DATABASE_URL"); url != "" {
		return Config{DatabaseURL: url}, nil
	}

	host := envOrDefault("DB_HOST", "localhost")
	port := envOrDefault("DB_PORT", "5433") // docker-compose пробрасывает 5433→5432
	user := envOrDefault("DB_USER", "wishlist")
	password := envOrDefault("DB_PASSWORD", "wishlist")
	name := envOrDefault("DB_NAME", "wishlist")
	sslmode := envOrDefault("DB_SSLMODE", "disable")

	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, name, sslmode,
	)

	return Config{DatabaseURL: url}, nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
