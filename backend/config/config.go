package config

import (
	"os"
)

type Config struct {
	Port               string
	DBPath             string
	JWTSecret          string
	PayPalClientID     string
	PayPalClientSecret string
	PayPalMode         string // "sandbox" or "live"
	AppBaseURL         string
}

var AppConfig *Config

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "wallpaper.db"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "wallart-super-secret-jwt-token-2026"
	}

	paypalClientID := os.Getenv("PAYPAL_CLIENT_ID")
	if paypalClientID == "" {
		paypalClientID = "sb" // PayPal standard sandbox test client ID
	}

	paypalSecret := os.Getenv("PAYPAL_CLIENT_SECRET")

	paypalMode := os.Getenv("PAYPAL_MODE")
	if paypalMode == "" {
		paypalMode = "sandbox"
	}

	appBaseURL := os.Getenv("APP_BASE_URL")
	if appBaseURL == "" {
		appBaseURL = "http://localhost:5173"
	}

	AppConfig = &Config{
		Port:               port,
		DBPath:             dbPath,
		JWTSecret:          jwtSecret,
		PayPalClientID:     paypalClientID,
		PayPalClientSecret: paypalSecret,
		PayPalMode:         paypalMode,
		AppBaseURL:         appBaseURL,
	}

	return AppConfig
}
