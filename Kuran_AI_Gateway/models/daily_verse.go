package models

import (
	"time"

	"gorm.io/gorm"
)

// DailyVerse günün ayeti ana modeli. Dile bağımsız alanlar burada.
type DailyVerse struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	SurahNumber int            `gorm:"not null" json:"surahNumber"`
	AyahNumber  int            `gorm:"not null" json:"ayahNumber"`
	ArabicText  string         `gorm:"type:text;not null" json:"arabicText"`
	DisplayDate *string        `gorm:"type:date;uniqueIndex" json:"displayDate,omitempty"` // YYYY-MM-DD, nil = havuzda
	IsActive    bool           `gorm:"default:true" json:"isActive"`
	SortOrder   int            `gorm:"default:0" json:"sortOrder"`

	Translations []DailyVerseTranslation `gorm:"foreignKey:DailyVerseID;constraint:OnDelete:CASCADE" json:"translations,omitempty"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName GORM tablo adı.
func (DailyVerse) TableName() string {
	return "daily_verses"
}

// DailyVerseTranslation çeviri modeli. İleride dil sayısı artabilir (tr, en, ar, de, fr, ...).
type DailyVerseTranslation struct {
	ID              uint   `gorm:"primarykey" json:"id"`
	DailyVerseID    uint   `gorm:"not null;index" json:"dailyVerseId"`
	LanguageCode    string `gorm:"size:8;not null" json:"languageCode"`    // "tr", "en", "ar" vb.
	TranslationText string `gorm:"type:text;not null" json:"translationText"`
	SourceName      string `gorm:"size:128" json:"sourceName"` // Surenin o dildeki adı: "Fatiha", "Al-Fatiha"

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName GORM tablo adı.
func (DailyVerseTranslation) TableName() string {
	return "daily_verse_translations"
}
