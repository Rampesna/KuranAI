package models

import (
	"time"

	"gorm.io/gorm"
)

// AiProvider yapay zeka sağlayıcısı.
// Sütun sırası: id, name, slug, created_at, updated_at, deleted_at.
type AiProvider struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Name string `gorm:"not null;uniqueIndex" json:"name"`
	Slug string `gorm:"size:50;uniqueIndex" json:"slug"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName GORM tablo adı.
func (AiProvider) TableName() string {
	return "ai_providers"
}
