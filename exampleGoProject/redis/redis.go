package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/talhacan/gateway/config"
)

// Client global Redis client.
var Client *redis.Client

// Connect Redis bağlantısı.
func Connect(appConfig *config.Config) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", appConfig.Redis.Host, appConfig.Redis.Port),
		Password: appConfig.Redis.Password,
		DB:       appConfig.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if pingErr := Client.Ping(ctx).Err(); pingErr != nil {
		return fmt.Errorf("redis connect: %w", pingErr)
	}

	log.Println("redis connected")
	return nil
}

// KeyPrefix JWT key prefix (token type + token id veya token string).
const KeyPrefixAccess = "jwt:access:"
const KeyPrefixRefresh = "jwt:refresh:"
