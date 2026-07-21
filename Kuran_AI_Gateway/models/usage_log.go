package models

import (
	"time"

	"gorm.io/gorm"
)

// UsageAction kullanım türü.
type UsageAction string

const (
	UsageActionMessage UsageAction = "message"
	UsageActionTTS     UsageAction = "tts"
)

// UsageLog mesaj ve sesli okuma kullanım kaydı (günlük limit hesabı için).
type UsageLog struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"userId"`
	Action    UsageAction    `gorm:"size:16;not null" json:"action"`
	CreatedAt time.Time      `gorm:"not null" json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName GORM tablo adı.
func (UsageLog) TableName() string {
	return "usage_logs"
}
