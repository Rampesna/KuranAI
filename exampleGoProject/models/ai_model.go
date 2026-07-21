package models

import (
	"time"

	"gorm.io/gorm"
)

// AiModel sağlayıcıya ait API modeli.
// Sütun sırası: id, provider_id, name, api_id, created_at, updated_at, deleted_at.
type AiModel struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	ProviderID uint   `gorm:"not null;index" json:"providerId"`
	Name       string `gorm:"not null" json:"name"`
	ApiID      string `gorm:"not null;size:120" json:"apiId"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName GORM tablo adı.
func (AiModel) TableName() string {
	return "ai_models"
}
