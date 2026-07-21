package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	OAuth    OAuthConfig
	WebAuthn WebAuthnConfig
	OpenAI   OpenAIConfig
	AppleIAP AppleIAPConfig
}

type OpenAIConfig struct {
	APIKey string
}

type AppleIAPConfig struct {
	SharedSecret string // App Store Connect → App → App Information → App-Specific Shared Secret
	BundleID     string
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

type OAuthConfig struct {
	GoogleClientIDs    []string
	GoogleClientSecret string
	AppleClientIDs     []string
	AppleTeamID        string
	AppleKeyID         string
	ApplePrivateKey    string
}

type WebAuthnConfig struct {
	RPID      string
	RPName    string
	RPOrigins []string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "kuranai_user"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "kuranai_pgsql"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			AccessSecret:     getEnv("JWT_ACCESS_SECRET", ""),
			RefreshSecret:    getEnv("JWT_REFRESH_SECRET", ""),
			AccessExpireMin:  getEnvInt("JWT_ACCESS_EXPIRE_MIN", 1440),
			RefreshExpireDay: getEnvInt("JWT_REFRESH_EXPIRE_DAY", 7),
			Issuer:           getEnv("JWT_ISSUER", "kuran-ai-gateway"),
		},
		OAuth: OAuthConfig{
			GoogleClientIDs:    getEnvList("GOOGLE_CLIENT_IDS", getEnv("GOOGLE_CLIENT_ID", "")),
			GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			AppleClientIDs:     getEnvList("APPLE_CLIENT_IDS", getEnv("APPLE_CLIENT_ID", "")),
			AppleTeamID:        getEnv("APPLE_TEAM_ID", ""),
			AppleKeyID:         getEnv("APPLE_KEY_ID", ""),
			ApplePrivateKey:    getEnv("APPLE_PRIVATE_KEY", ""),
		},
		WebAuthn: WebAuthnConfig{
			RPID:      getEnv("WEBAUTHN_RP_ID", "kuranai.talhacan.com"),
			RPName:    getEnv("WEBAUTHN_RP_NAME", "Kuran AI"),
			RPOrigins: getEnvList("WEBAUTHN_RP_ORIGINS", "https://kuranai.talhacan.com"),
		},
		OpenAI: OpenAIConfig{
			APIKey: getEnv("OPENAI_API_KEY", ""),
		},
		AppleIAP: AppleIAPConfig{
			SharedSecret: getEnv("APPLE_IAP_SHARED_SECRET", ""),
			BundleID:     getEnv("APPLE_BUNDLE_ID", "talhacan.KuranAIChat"),
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

func getEnvList(key, fallback string) []string {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		raw = fallback
	}
	values := make([]string, 0)
	for _, value := range strings.Split(raw, ",") {
		if value = strings.TrimSpace(value); value != "" {
			values = append(values, value)
		}
	}
	return values
}
