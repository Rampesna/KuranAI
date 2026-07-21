package models

import (
	"time"

	"gorm.io/gorm"
)

// SubscriptionPackage abonelik paketi tanımı.
type SubscriptionPackage struct {
	ID                 uint           `gorm:"primarykey" json:"id"`
	Name               string         `gorm:"size:64;not null" json:"name"`
	NameEn             string         `gorm:"size:64;not null" json:"nameEn"`
	Slug               string         `gorm:"uniqueIndex;size:32;not null" json:"slug"`
	DailyQuestionLimit *int           `json:"dailyQuestionLimit,omitempty"` // nil = sınırsız
	DailyTTSLimit      *int           `json:"dailyTTSLimit,omitempty"`     // nil = sınırsız
	MonthlyPriceTL     float64        `gorm:"type:decimal(10,2);not null" json:"monthlyPriceTL"`
	YearlyPriceTL      float64        `gorm:"type:decimal(10,2);not null" json:"yearlyPriceTL"`
	SortOrder          int            `gorm:"default:0" json:"sortOrder"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName GORM tablo adı.
func (SubscriptionPackage) TableName() string {
	return "subscription_packages"
}
