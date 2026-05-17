package config

import (
	"log"
	"os"
	"time"
)

type AppConfig struct {
	JWTSecret            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func LoadAppConfig() AppConfig {

	// parse duration, kalau error pakai default
	// kalau JWTSecret kosong, fatal
	accessDuration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		accessDuration = time.Minute * 15
	}
	refreshDuration, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_DURATION")) // 168h = 7 hari
	if err != nil {
		refreshDuration = time.Hour * 24 * 7
	}
	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	return AppConfig{
		JWTSecret:            jwtSecret,
		AccessTokenDuration:  accessDuration,
		RefreshTokenDuration: refreshDuration,
	}

}
