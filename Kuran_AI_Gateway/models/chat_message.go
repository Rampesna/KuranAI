package models

import (
	"time"

	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

// ChatMessage sohbet mesajı. role: "user" | "assistant", sourceType: "meal" | "tefsir" | "hadis" vb. (genişletilebilir).
type ChatMessage struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	ChatSessionID uint           `gorm:"not null;index:idx_chat_message_session" json:"chatSessionId"`
	Role          string         `gorm:"type:varchar(20);not null" json:"role"`       // user | assistant
	Content       string         `gorm:"type:text;not null" json:"content"`           // orijinal metin
	SourceType    string         `gorm:"type:varchar(50);not null" json:"sourceType"` // meal | tefsir | hadis | ...
	Embedding     *pgvector.Vector `gorm:"type:vector(3072)" json:"-"`                // user ve assistant mesajlarında (her ikisi de embed edilir)
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	ChatSession *ChatSession `gorm:"foreignKey:ChatSessionID" json:"-"`
}

// TableName GORM tablo adı.
func (ChatMessage) TableName() string {
	return "chat_messages"
}
