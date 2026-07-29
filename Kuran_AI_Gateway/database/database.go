package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rampesna/kuran-ai-gateway/config"
	"github.com/rampesna/kuran-ai-gateway/models"
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

	customLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			ParameterizedQueries:      true, // parametreleri (embedding vektörü vb.) loglama
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
		},
	)
	var connectErr error
	DB, connectErr = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: customLogger.LogMode(logger.Info),
	})
	if connectErr != nil {
		return fmt.Errorf("database connect: %w", connectErr)
	}

	// pgvector extension
	if vectorExtErr := DB.Exec("CREATE EXTENSION IF NOT EXISTS vector").Error; vectorExtErr != nil {
		log.Printf("warning: pgvector extension: %v (skip if not needed)", vectorExtErr)
	}

	// Migration: RUN_MIGRATE=1 ise veya varsayılan olarak çalıştır. RUN_MIGRATE=0 ile kapatılabilir.
	if os.Getenv("RUN_MIGRATE") != "0" {
		if migrateErr := runMigrations(); migrateErr != nil {
			return fmt.Errorf("migrations: %w", migrateErr)
		}
	}

	log.Println("database connected and migrated")
	return nil
}

func runMigrations() error {
	// Rename Turkish column names to English (sure->surah, ayet->ayah, sura_number->surah_number)
	DB.Exec(`
		DO $$
		BEGIN
		  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema='public' AND table_name='ayah_exegesis_embeddings' AND column_name='sure') THEN
		    ALTER TABLE ayah_exegesis_embeddings RENAME COLUMN sure TO surah;
		  END IF;
		  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema='public' AND table_name='ayah_exegesis_embeddings' AND column_name='ayet') THEN
		    ALTER TABLE ayah_exegesis_embeddings RENAME COLUMN ayet TO ayah;
		  END IF;
		  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema='public' AND table_name='surah_exegesis_embeddings' AND column_name='sura_number') THEN
		    ALTER TABLE surah_exegesis_embeddings RENAME COLUMN sura_number TO surah_number;
		  END IF;
		END $$;
	`)
	// ayah: int -> text (14-16 aralık destegi), surah: int -> bigint
	DB.Exec(`
		DO $$
		BEGIN
		  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema='public' AND table_name='ayah_exegesis_embeddings' AND column_name='ayah' AND data_type='integer') THEN
		    ALTER TABLE ayah_exegesis_embeddings ALTER COLUMN ayah TYPE text USING ayah::text;
		  END IF;
		  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema='public' AND table_name='ayah_exegesis_embeddings' AND column_name='surah' AND data_type='integer') THEN
		    ALTER TABLE ayah_exegesis_embeddings ALTER COLUMN surah TYPE bigint USING surah::bigint;
		  END IF;
		END $$;
	`)

	return DB.AutoMigrate(
		&models.User{},
		&models.Manager{},
		&models.JWTToken{},
		&models.MealEmbedding{},
		&models.SurahExegesisEmbedding{},
		&models.AyahExegesisEmbedding{},
		&models.ChatSession{},
		&models.ChatMessage{},
		&models.SubscriptionPackage{},
		&models.UserSubscription{},
		&models.GooglePlayPurchase{},
		&models.UsageLog{},
		&models.DailyVerse{},
		&models.DailyVerseTranslation{},
		&models.SupportRequest{},
		&models.PasskeyCredential{},
	)
}
