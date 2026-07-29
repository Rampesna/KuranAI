package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/rampesna/kuran-ai-gateway/database"
	"github.com/rampesna/kuran-ai-gateway/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SubscriptionService abonelik işlemleri.
type SubscriptionService struct{}

// NewSubscriptionService SubscriptionService oluşturur.
func NewSubscriptionService() *SubscriptionService {
	return &SubscriptionService{}
}

// GetActiveSubscription kullanıcının geçerli aboneliğini döner (expires_at > now). Yoksa nil.
func (s *SubscriptionService) GetActiveSubscription(ctx context.Context, userID uint) (*models.UserSubscription, *models.SubscriptionPackage, error) {
	var sub models.UserSubscription
	err := database.DB.WithContext(ctx).
		Where("user_id = ? AND expires_at > ?", userID, time.Now().UTC()).
		Order("expires_at DESC").
		First(&sub).Error
	if err != nil {
		return nil, nil, nil // hata = abonelik yok
	}
	var pkg models.SubscriptionPackage
	if err := database.DB.WithContext(ctx).First(&pkg, sub.PackageID).Error; err != nil {
		return nil, nil, err
	}
	return &sub, &pkg, nil
}

// GetPackages tüm paketleri döner.
func (s *SubscriptionService) GetPackages(ctx context.Context) ([]models.SubscriptionPackage, error) {
	var pkgs []models.SubscriptionPackage
	err := database.DB.WithContext(ctx).Order("sort_order ASC, id ASC").Find(&pkgs).Error
	return pkgs, err
}

// CreateSubscription Apple satın alma sonrası abonelik kaydı oluşturur.
// Restore için: aynı transaction ID varsa mevcut kaydı döner (idempotent).
func (s *SubscriptionService) CreateSubscription(ctx context.Context, userID uint, packageID uint, isMonthly bool, appleTransactionID, appleProductID string, expiresAtOverride *time.Time) (*models.UserSubscription, error) {
	var existing models.UserSubscription
	if appleTransactionID != "" && database.DB.WithContext(ctx).Where("user_id = ? AND apple_transaction_id = ?", userID, appleTransactionID).First(&existing).Error == nil {
		updates := map[string]interface{}{
			"apple_product_id": appleProductID,
			"package_id":       packageID,
			"is_monthly":       isMonthly,
			"provider":         "apple",
		}
		if expiresAtOverride != nil && expiresAtOverride.After(existing.ExpiresAt) {
			updates["expires_at"] = *expiresAtOverride
		}
		if err := database.DB.WithContext(ctx).Model(&existing).Updates(updates).Error; err != nil {
			return nil, err
		}
		return &existing, nil
	}

	var pkg models.SubscriptionPackage
	if err := database.DB.WithContext(ctx).First(&pkg, packageID).Error; err != nil {
		return nil, fmt.Errorf("package not found: %w", err)
	}

	now := time.Now().UTC()
	expiresAt := now
	if expiresAtOverride != nil && expiresAtOverride.After(now) {
		expiresAt = *expiresAtOverride
	} else if isMonthly {
		expiresAt = now.AddDate(0, 1, 0)
	} else {
		expiresAt = now.AddDate(1, 0, 0)
	}

	sub := models.UserSubscription{
		UserID:             userID,
		PackageID:          packageID,
		IsMonthly:          isMonthly,
		StartedAt:          now,
		ExpiresAt:          expiresAt,
		AppleTransactionID: appleTransactionID,
		AppleProductID:     appleProductID,
		Provider:           "apple",
	}
	if err := database.DB.WithContext(ctx).Create(&sub).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

// CreateGoogleSubscription stores a Google-verified entitlement and binds a
// purchase token to one Kuran AI account. Only a SHA-256 token digest is stored.
func (s *SubscriptionService) CreateGoogleSubscription(
	ctx context.Context,
	userID uint,
	packageID uint,
	isMonthly bool,
	info GooglePlayPurchaseInfo,
) (*models.UserSubscription, error) {
	if stringsTrimmed := fmt.Sprint(userID); info.ObfuscatedAccountID != "" && info.ObfuscatedAccountID != stringsTrimmed {
		return nil, ErrGoogleAccountConflict
	}
	if info.PurchaseToken == "" || info.ProductID == "" || !info.ExpiresAt.After(time.Now().UTC()) {
		return nil, ErrInvalidGooglePurchase
	}

	tokenDigest := sha256.Sum256([]byte(info.PurchaseToken))
	tokenHash := hex.EncodeToString(tokenDigest[:])
	rawDigest := sha256.Sum256(info.RawJSON)
	rawHash := hex.EncodeToString(rawDigest[:])
	now := time.Now().UTC()
	var saved models.UserSubscription

	err := database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var purchase models.GooglePlayPurchase
		findErr := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("token_hash = ?", tokenHash).
			First(&purchase).Error
		switch {
		case findErr == nil:
			if purchase.UserID != userID {
				return ErrGoogleAccountConflict
			}
			if err := tx.Model(&purchase).Updates(map[string]interface{}{
				"package_name":          info.PackageName,
				"product_id":            info.ProductID,
				"order_id":              info.OrderID,
				"obfuscated_account_id": info.ObfuscatedAccountID,
				"environment":           info.Environment,
				"purchased_at":          info.PurchasedAt,
				"expires_at":            info.ExpiresAt,
				"acknowledged":          info.Acknowledged,
				"raw_sha256":            rawHash,
			}).Error; err != nil {
				return err
			}
		case errors.Is(findErr, gorm.ErrRecordNotFound):
			purchase = models.GooglePlayPurchase{
				TokenHash:           tokenHash,
				UserID:              userID,
				PackageName:         info.PackageName,
				ProductID:           info.ProductID,
				OrderID:             info.OrderID,
				ObfuscatedAccountID: info.ObfuscatedAccountID,
				Environment:         info.Environment,
				PurchasedAt:         info.PurchasedAt,
				ExpiresAt:           info.ExpiresAt,
				Acknowledged:        info.Acknowledged,
				RawSHA256:           rawHash,
			}
			if err := tx.Create(&purchase).Error; err != nil {
				return err
			}
		default:
			return findErr
		}

		var existing models.UserSubscription
		findSubErr := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND google_purchase_hash = ?", userID, tokenHash).
			First(&existing).Error
		if findSubErr == nil {
			if err := tx.Model(&existing).Updates(map[string]interface{}{
				"package_id":           packageID,
				"is_monthly":           isMonthly,
				"expires_at":           info.ExpiresAt,
				"provider":             "google",
				"google_product_id":    info.ProductID,
				"google_purchase_hash": tokenHash,
			}).Error; err != nil {
				return err
			}
			saved = existing
			saved.PackageID = packageID
			saved.IsMonthly = isMonthly
			saved.ExpiresAt = info.ExpiresAt
			saved.Provider = "google"
			saved.GoogleProductID = info.ProductID
			saved.GooglePurchaseHash = tokenHash
			return nil
		}
		if !errors.Is(findSubErr, gorm.ErrRecordNotFound) {
			return findSubErr
		}

		saved = models.UserSubscription{
			UserID:             userID,
			PackageID:          packageID,
			IsMonthly:          isMonthly,
			StartedAt:          info.PurchasedAt,
			ExpiresAt:          info.ExpiresAt,
			Provider:           "google",
			GooglePurchaseHash: tokenHash,
			GoogleProductID:    info.ProductID,
		}
		if saved.StartedAt.IsZero() {
			saved.StartedAt = now
		}
		return tx.Create(&saved).Error
	})
	if err != nil {
		return nil, err
	}
	return &saved, nil
}

// SeedPackages paketleri seed eder (yoksa oluşturur).
func (s *SubscriptionService) SeedPackages(ctx context.Context) error {
	pkgs := []models.SubscriptionPackage{
		{Name: "Bronze", NameEn: "Bronze", Slug: "bronze", DailyQuestionLimit: intPtr(100), DailyTTSLimit: intPtr(5), MonthlyPriceTL: 29.99, YearlyPriceTL: 299.99, SortOrder: 1},
		{Name: "Gümüş", NameEn: "Silver", Slug: "silver", DailyQuestionLimit: intPtr(250), DailyTTSLimit: intPtr(50), MonthlyPriceTL: 49.99, YearlyPriceTL: 499.99, SortOrder: 2},
		{Name: "Altın", NameEn: "Gold", Slug: "gold", DailyQuestionLimit: nil, DailyTTSLimit: intPtr(100), MonthlyPriceTL: 79.99, YearlyPriceTL: 799.99, SortOrder: 3},
		{Name: "Elmas", NameEn: "Diamond", Slug: "diamond", DailyQuestionLimit: nil, DailyTTSLimit: nil, MonthlyPriceTL: 149.99, YearlyPriceTL: 1499.99, SortOrder: 4},
	}
	for i := range pkgs {
		if err := database.DB.WithContext(ctx).Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "slug"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name", "name_en", "daily_question_limit", "daily_tts_limit",
				"monthly_price_tl", "yearly_price_tl", "sort_order", "updated_at", "deleted_at",
			}),
		}).Create(&pkgs[i]).Error; err != nil {
			return err
		}
	}
	return nil
}

func intPtr(i int) *int {
	return &i
}
