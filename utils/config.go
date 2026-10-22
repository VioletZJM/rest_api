package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port   string
	DBPath string
}

func LoadConfig() (Config, error) {
	path := getenvDefault("DB_PATH", "api.db")
	port := strings.TrimSpace(getenvDefault("PORT", "8090"))

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return Config{}, fmt.Errorf("invalid port : %w", err)
	} else if portInt < 1 || portInt > 65535 {
		return Config{}, fmt.Errorf("invalid port : %d want 1 < PORT < 65535", portInt)
	}

	cfg := Config{
		Port:   strconv.Itoa(portInt),
		DBPath: path,
	}
	return cfg, nil
}

func getenvDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
