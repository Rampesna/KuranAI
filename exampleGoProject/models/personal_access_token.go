package models

import (
	"time"

	"gorm.io/gorm"
)

// TokenableType token sahibi tipi: "user" veya "manager".
type TokenableType string

const (
	TokenableTypeUser    TokenableType = "user"
	TokenableTypeManager TokenableType = "manager"
)

// PersonalAccessToken JWT token kaydı.
// Sütun sırası: id, tokenable_id, tokenable_type, access_token, refresh_token, expires_at, user_agent, ip_address, device_info, device_type, created_at, updated_at, deleted_at.
type PersonalAccessToken struct {
	ID             uint          `gorm:"primarykey" json:"id"`
	TokenableID    uint          `gorm:"not null;index:idx_tokenable" json:"tokenableId"`
	TokenableType  TokenableType `gorm:"not null;index:idx_tokenable;size:50" json:"tokenableType"`
	AccessToken    string        `gorm:"not null;uniqueIndex" json:"-"`
	RefreshToken   string        `gorm:"not null;uniqueIndex" json:"-"`
	ExpiresAt      time.Time     `gorm:"not null" json:"expiresAt"`
	UserAgent      string        `gorm:"size:500" json:"userAgent"`
	IPAddress      string        `gorm:"size:45" json:"ipAddress"`
	DeviceInfo     string        `gorm:"size:500" json:"deviceInfo"`
	DeviceType     string        `gorm:"size:50" json:"deviceType"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName GORM tablo adı.
func (PersonalAccessToken) TableName() string {
	return "personal_access_tokens"
}
