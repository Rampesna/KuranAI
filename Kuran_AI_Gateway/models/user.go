package models

import (
	"time"

	"gorm.io/gorm"
)

// User kullanıcı modeli. Apple/Google OAuth için apple_id, google_id alanları.
type User struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	Email    string `gorm:"uniqueIndex" json:"email"`
	Username string `gorm:"uniqueIndex" json:"username"`
	Password string `json:"-"` // OAuth ile girişte boş olabilir
	Name     string `gorm:"not null" json:"name"`
	Status   int    `gorm:"default:1" json:"status"` // 0: inactive, 1: active

	// OAuth providers (NULL when not used; boş string unique ihlaline yol açar)
	GoogleID *string `gorm:"uniqueIndex;size:255" json:"-"`
	AppleID  *string `gorm:"uniqueIndex;size:255" json:"-"`

	AvatarURL string `gorm:"size:512" json:"avatarUrl,omitempty"`

	// HasPassword DB'de tutulmaz, JSON response'ta şifre varlığını bildirir.
	HasPassword bool `gorm:"-" json:"hasPassword"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName GORM tablo adı.
func (User) TableName() string {
	return "users"
}
