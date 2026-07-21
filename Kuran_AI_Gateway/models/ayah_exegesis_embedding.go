package models

import (
	"time"

	"github.com/pgvector/pgvector-go"
	"gorm.io/datatypes"
)

// AyahExegesisEmbedding ayah exegesis embeddings. Ayah text olabilir ("1", "14-16").
type AyahExegesisEmbedding struct {
	ID         uint            `gorm:"primarykey" json:"id"`
	Surah      int64           `gorm:"column:surah;not null;index:idx_ayah_exegesis_surah_ayah" json:"surah"`
	Ayah       string          `gorm:"column:ayah;not null;type:text;index:idx_ayah_exegesis_surah_ayah" json:"ayah"`
	Data       datatypes.JSON  `gorm:"type:jsonb;not null" json:"data"`
	Embedding  pgvector.Vector `gorm:"type:vector(3072)" json:"-"`
	TokenCount *int            `json:"tokenCount"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  *time.Time      `json:"updatedAt"`
}

// TableName GORM tablo adı.
func (AyahExegesisEmbedding) TableName() string {
	return "ayah_exegesis_embeddings"
}
