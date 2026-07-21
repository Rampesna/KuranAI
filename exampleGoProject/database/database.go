package database

import (
	"fmt"
	"log"
	"os"

	"github.com/talhacan/gateway/config"
	"github.com/talhacan/gateway/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB global GORM instance.
var DB *gorm.DB

// Connect PostgreSQL bağlantısı ve migration.
func Connect(appConfig *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		appConfig.Database.Host,
		appConfig.Database.Port,
		appConfig.Database.User,
		appConfig.Database.Password,
		appConfig.Database.DBName,
		appConfig.Database.SSLMode,
	)

	var connectErr error
	DB, connectErr = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if connectErr != nil {
		return fmt.Errorf("database connect: %w", connectErr)
	}

	// pgvector extension (vektör alanları için)
	if vectorExtErr := DB.Exec("CREATE EXTENSION IF NOT EXISTS vector").Error; vectorExtErr != nil {
		log.Printf("warning: pgvector extension: %v (skip if not needed)", vectorExtErr)
	}

	if os.Getenv("DB_RESET") == "1" {
		log.Println("DB_RESET=1: dropping all tables and re-migrating")
		if resetErr := Reset(); resetErr != nil {
			return fmt.Errorf("database reset: %w", resetErr)
		}
	} else if migrateErr := runMigrations(); migrateErr != nil {
		return fmt.Errorf("migrations: %w", migrateErr)
	}

	if seedErr := runSeedAI(); seedErr != nil {
		return fmt.Errorf("seed AI: %w", seedErr)
	}

	log.Println("database connected and migrated")
	return nil
}

// Reset tüm tabloları siler ve AutoMigrate ile baştan oluşturur (DB_RESET=1 ile tetiklenir).
func Reset() error {
	migrator := DB.Migrator()
	tables := []interface{}{
		&models.AiModel{},
		&models.AiProvider{},
		&models.ChatMessage{},
		&models.ChatSession{},
		&models.PersonalAccessToken{},
		&models.Todo{},
		&models.User{},
		&models.Manager{},
	}
	for _, table := range tables {
		if migrator.HasTable(table) {
			if dropErr := migrator.DropTable(table); dropErr != nil {
				return fmt.Errorf("drop table: %w", dropErr)
			}
		}
	}
	if migrateErr := runMigrations(); migrateErr != nil {
		return migrateErr
	}
	return runSeedAI()
}

func runMigrations() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Manager{},
		&models.Todo{},
		&models.PersonalAccessToken{},
		&models.AiProvider{},
		&models.AiModel{},
		&models.ChatSession{},
		&models.ChatMessage{},
	)
}
