package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBUrl      string
	ServerPort string
	DBMaxConn  int
}

func Load() *Config {
	return &Config{
		DBUrl:      getEnv("DATABASE_URL", ""),
		ServerPort: getEnv("SERVER_PORT", ":8080"),
		DBMaxConn:  getEnvAsInt("DB_MAX_CONN", 25),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}