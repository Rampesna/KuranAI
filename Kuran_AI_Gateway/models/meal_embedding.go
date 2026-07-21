package models

import (
	"time"

	"github.com/pgvector/pgvector-go"
	"gorm.io/datatypes"
)

// MealEmbedding meal (çeviri) embeddings tablosu. kuran_embeddings_structure.sql ile uyumlu.
type MealEmbedding struct {
	ID         uint            `gorm:"primarykey" json:"id"`
	Surah      int             `gorm:"not null;index:idx_meal_surah" json:"surah"`
	Ayah       int             `gorm:"not null;index:idx_meal_ayah" json:"ayah"`
	Data       datatypes.JSON  `gorm:"type:jsonb;not null" json:"data"`
	Embedding  pgvector.Vector `gorm:"type:vector(3072)" json:"-"`
	TokenCount *int            `json:"tokenCount"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  *time.Time      `json:"updatedAt"`
}

// TableName GORM tablo adı.
func (MealEmbedding) TableName() string {
	return "meal_embeddings"
}
