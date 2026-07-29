package models

import "time"

// GooglePlayPurchase is the minimum audit record retained for an independently
// verified Play Billing subscription. The raw purchase token is never stored.
type GooglePlayPurchase struct {
	TokenHash           string    `gorm:"type:char(64);primaryKey" json:"-"`
	UserID              uint      `gorm:"index;not null" json:"userId"`
	PackageName         string    `gorm:"size:160;not null" json:"packageName"`
	ProductID           string    `gorm:"size:160;not null" json:"productId"`
	OrderID             string    `gorm:"size:160;index" json:"orderId"`
	ObfuscatedAccountID string    `gorm:"size:64" json:"-"`
	Environment         string    `gorm:"size:24;not null" json:"environment"`
	PurchasedAt         time.Time `gorm:"not null" json:"purchasedAt"`
	ExpiresAt           time.Time `gorm:"index;not null" json:"expiresAt"`
	Acknowledged        bool      `gorm:"not null;default:false" json:"acknowledged"`
	RawSHA256           string    `gorm:"type:char(64);not null" json:"-"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

func (GooglePlayPurchase) TableName() string {
	return "google_play_purchases"
}
