package models

import (
	"time"

	"github.com/pgvector/pgvector-go"
	"gorm.io/datatypes"
)

// SurahExegesisEmbedding surah exegesis embeddings.
type SurahExegesisEmbedding struct {
	ID          uint            `gorm:"primarykey" json:"id"`
	SurahNumber int             `gorm:"column:surah_number;not null;uniqueIndex:idx_surah_exegesis_surah" json:"surahNumber"`
	Data        datatypes.JSON  `gorm:"type:jsonb;not null" json:"data"`
	Embedding  pgvector.Vector `gorm:"type:vector(3072)" json:"-"`
	TokenCount *int            `json:"tokenCount"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  *time.Time      `json:"updatedAt"`
}

// TableName GORM tablo adı.
func (SurahExegesisEmbedding) TableName() string {
	return "surah_exegesis_embeddings"
}
