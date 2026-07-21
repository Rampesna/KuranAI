package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	AccessSecret     string
	RefreshSecret    string
	AccessExpireMin  int
	RefreshExpireDay int
	Issuer           string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "talhacan"),
			Password: getEnv("DB_PASSWORD", ".4br3x4cp6X."),
			DBName:   getEnv("DB_NAME", "talhacan_postgres"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			AccessSecret:     getEnv("JWT_ACCESS_SECRET", "access-secret-change-in-production"),
			RefreshSecret:    getEnv("JWT_REFRESH_SECRET", "refresh-secret-change-in-production"),
			AccessExpireMin:  getEnvInt("JWT_ACCESS_EXPIRE_MIN", 1440), // 24 saat
			RefreshExpireDay: getEnvInt("JWT_REFRESH_EXPIRE_DAY", 7),   // 1 hafta
			Issuer:           getEnv("JWT_ISSUER", "gateway"),
		},
	}
}

func getEnv(key, defaultVal string) string {
	if envValue := os.Getenv(key); envValue != "" {
		return envValue
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if envValue := os.Getenv(key); envValue != "" {
		if parsedInt, parseErr := strconv.Atoi(envValue); parseErr == nil {
			return parsedInt
		}
	}
	return defaultVal
}
