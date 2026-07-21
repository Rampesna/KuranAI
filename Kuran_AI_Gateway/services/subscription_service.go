package services

import (
	"context"
	"fmt"
	"time"

	"github.com/rampesna/kuran-ai-gateway/database"
	"github.com/rampesna/kuran-ai-gateway/models"
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
	}
	if err := database.DB.WithContext(ctx).Create(&sub).Error; err != nil {
		return nil, err
	}
	return &sub, nil
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
