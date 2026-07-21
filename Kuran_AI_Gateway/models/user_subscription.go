package models

import (
	"time"

	"gorm.io/gorm"
)

// UserSubscription kullanıcının aktif aboneliği.
type UserSubscription struct {
	ID                   uint           `gorm:"primarykey" json:"id"`
	UserID               uint           `gorm:"index;not null" json:"userId"`
	PackageID            uint           `gorm:"not null" json:"packageId"`
	IsMonthly            bool           `gorm:"not null" json:"isMonthly"`
	StartedAt            time.Time      `gorm:"not null" json:"startedAt"`
	ExpiresAt            time.Time      `gorm:"not null" json:"expiresAt"`
	AppleTransactionID   string         `gorm:"size:255;index" json:"-"` // Apple orijinal transaction ID (yenileme takibi)
	AppleProductID       string         `gorm:"size:128" json:"-"`       // StoreKit product ID
	CreatedAt            time.Time      `json:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName GORM tablo adı.
func (UserSubscription) TableName() string {
	return "user_subscriptions"
}
