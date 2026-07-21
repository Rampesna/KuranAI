package models

import (
	"time"

	"gorm.io/gorm"
)

// ChatSession sohbet oturumu. auth_type: "user" | "manager", auth_id: users.id veya managers.id.
type ChatSession struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	AuthType     string         `gorm:"type:varchar(50);not null;index:idx_chat_session_auth" json:"authType"`
	AuthID       uint           `gorm:"not null;index:idx_chat_session_auth" json:"authId"`
	FirstMessage string         `gorm:"type:text" json:"firstMessage"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Messages []ChatMessage `gorm:"foreignKey:ChatSessionID" json:"-"`
}

// TableName GORM tablo adı.
func (ChatSession) TableName() string {
	return "chat_sessions"
}
