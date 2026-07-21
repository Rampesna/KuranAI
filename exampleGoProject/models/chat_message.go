package models

import (
	"time"

	"gorm.io/gorm"
)

// ChatMessage oturuma ait tek mesaj (user veya assistant).
// Sütun sırası: id, session_id, role, content, created_at, updated_at, deleted_at.
type ChatMessage struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	SessionID uint   `gorm:"not null;index" json:"sessionId"`
	Role      string `gorm:"not null;size:20" json:"role"` // user | assistant
	Content   string `gorm:"type:text;not null" json:"content"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName GORM tablo adı.
func (ChatMessage) TableName() string {
	return "chat_messages"
}
