package config

import (
	"os"
)

type ServiceConfig struct {
	Port     string
	LogLevel string
}

type GatewayConfig struct {
	Port       string
	KingdomURL string
	EconomyURL string
	MarketURL  string
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func LoadService() ServiceConfig {
	return ServiceConfig{
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

func LoadGateway() GatewayConfig {
	return GatewayConfig{
		Port:       getEnv("PORT", "8080"),
		KingdomURL: getEnv("KINGDOM_URL", "http://kingdom:8081"),
		EconomyURL: getEnv("ECONOMY_URL", "http://economy:8082"),
		MarketURL:  getEnv("MARKET_URL", "http://market:8083"),
	}
}
