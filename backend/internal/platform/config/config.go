package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL          string
	CORSOrigins          []string
	JWTSecret            string
	WebPushPublicKey     string
	WebPushPrivateKey    string
	WebPushSubject       string
	NotificationsEnabled bool
}

func Load() (Config, error) {
	// .env — удобство для локальной разработки; отсутствие файла не ошибка
	_ = godotenv.Load()

	var databaseURL string
	if url := os.Getenv("DATABASE_URL"); url != "" {
		databaseURL = url
	} else {
		host := envOrDefault("DB_HOST", "localhost")
		port := envOrDefault("DB_PORT", "5433") // docker-compose пробрасывает 5433→5432
		user := envOrDefault("DB_USER", "wishlist")
		password := envOrDefault("DB_PASSWORD", "wishlist")
		name := envOrDefault("DB_NAME", "wishlist")
		sslmode := envOrDefault("DB_SSLMODE", "disable")

		databaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			user, password, host, port, name, sslmode,
		)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		return Config{}, fmt.Errorf("JWT_SECRET must contain at least 32 characters")
	}

	webPushPublicKey := os.Getenv("WEB_PUSH_PUBLIC_KEY")
	webPushPrivateKey := os.Getenv("WEB_PUSH_PRIVATE_KEY")
	notificationsEnabled := webPushPublicKey != "" && webPushPrivateKey != ""
	if (webPushPublicKey == "") != (webPushPrivateKey == "") {
		return Config{}, fmt.Errorf("WEB_PUSH_PUBLIC_KEY and WEB_PUSH_PRIVATE_KEY must be configured together")
	}

	return Config{
		DatabaseURL: databaseURL,
		CORSOrigins: parseListEnv(
			"CORS_ORIGINS",
			[]string{"http://localhost:5173"},
		),
		JWTSecret:            jwtSecret,
		WebPushPublicKey:     webPushPublicKey,
		WebPushPrivateKey:    webPushPrivateKey,
		WebPushSubject:       envOrDefault("WEB_PUSH_SUBJECT", "mailto:admin@our-nest.ru"),
		NotificationsEnabled: notificationsEnabled,
	}, nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func parseListEnv(key string, fallback []string) []string {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}

	values := strings.Split(raw, ",")
	result := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return fallback
	}
	return result
}
