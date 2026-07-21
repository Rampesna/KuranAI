package models

import (
	"time"

	"gorm.io/gorm"
)

// ChatSession kullanıcının AI chat oturumu.
// Sütun sırası: id, user_id, created_at, updated_at, deleted_at.
type ChatSession struct {
	ID     uint `gorm:"primarykey" json:"id"`
	UserID uint `gorm:"not null;index" json:"userId"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName GORM tablo adı.
func (ChatSession) TableName() string {
	return "chat_sessions"
}
