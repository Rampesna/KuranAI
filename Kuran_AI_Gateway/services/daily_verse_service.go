package services

import (
	"context"
	"fmt"
	"time"

	"github.com/rampesna/kuran-ai-gateway/core"
	"github.com/rampesna/kuran-ai-gateway/database"
	"github.com/rampesna/kuran-ai-gateway/models"
	verseRequest "github.com/rampesna/kuran-ai-gateway/requests/verse"
	"gorm.io/gorm"
)

// DailyVerseService günün ayeti iş mantığı.
type DailyVerseService struct{}

// NewDailyVerseService yeni DailyVerseService oluşturur.
func NewDailyVerseService() *DailyVerseService {
	return &DailyVerseService{}
}

// GetToday bugünün ayetini döndürür. lang parametresi ile sadece ilgili dil çevirisi filtrelenir.
func (s *DailyVerseService) GetToday(ctx context.Context, lang string) *core.ServiceResponse {
	db := database.DB.WithContext(ctx)
	today := time.Now().Format("2006-01-02")

	var verse models.DailyVerse

	// 1. Bugüne sabitlenmiş ayet var mı?
	err := db.Where("display_date = ? AND is_active = ?", today, true).
		Preload("Translations", func(tx *gorm.DB) *gorm.DB {
			if lang != "" {
				return tx.Where("language_code = ?", lang)
			}
			return tx
		}).
		First(&verse).Error

	if err == nil {
		return core.NewServiceResponse(true, "ok", 200, verse)
	}

	// 2. Havuzdaki aktif ayetlerden gün bazlı seç
	var pool []models.DailyVerse
	db.Where("display_date IS NULL AND is_active = ?", true).
		Order("sort_order ASC, id ASC").
		Find(&pool)

	if len(pool) == 0 {
		return core.NewServiceResponse(false, "No verses available", 404, nil)
	}

	dayOfYear := time.Now().YearDay()
	picked := pool[dayOfYear%len(pool)]

	// Çevirileri yükle
	db.Where("daily_verse_id = ?", picked.ID).Find(&picked.Translations)
	if lang != "" {
		filtered := make([]models.DailyVerseTranslation, 0)
		for _, t := range picked.Translations {
			if t.LanguageCode == lang {
				filtered = append(filtered, t)
			}
		}
		picked.Translations = filtered
	}

	return core.NewServiceResponse(true, "ok", 200, picked)
}

// List tüm ayetleri listeler (manager için).
func (s *DailyVerseService) List(ctx context.Context) *core.ServiceResponse {
	db := database.DB.WithContext(ctx)
	var verses []models.DailyVerse
	db.Preload("Translations").Order("sort_order ASC, id ASC").Find(&verses)
	return core.NewServiceResponse(true, "ok", 200, verses)
}

// GetByID tek ayet detayı.
func (s *DailyVerseService) GetByID(ctx context.Context, id uint) *core.ServiceResponse {
	db := database.DB.WithContext(ctx)
	var verse models.DailyVerse
	if err := db.Preload("Translations").First(&verse, id).Error; err != nil {
		return core.NewServiceResponse(false, "Verse not found", 404, nil)
	}
	return core.NewServiceResponse(true, "ok", 200, verse)
}

// Create yeni günün ayeti oluşturur.
func (s *DailyVerseService) Create(ctx context.Context, req *verseRequest.CreateDailyVerseRequest) *core.ServiceResponse {
	db := database.DB.WithContext(ctx)

	verse := models.DailyVerse{
		SurahNumber: req.SurahNumber,
		AyahNumber:  req.AyahNumber,
		ArabicText:  req.ArabicText,
		DisplayDate: req.DisplayDate,
		IsActive:    true,
	}
	if req.IsActive != nil {
		verse.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		verse.SortOrder = *req.SortOrder
	}

	tx := db.Begin()

	if err := tx.Create(&verse).Error; err != nil {
		tx.Rollback()
		return core.NewServiceResponse(false, fmt.Sprintf("Create failed: %v", err), 500, nil)
	}

	for _, t := range req.Translations {
		translation := models.DailyVerseTranslation{
			DailyVerseID:    verse.ID,
			LanguageCode:    t.LanguageCode,
			TranslationText: t.TranslationText,
			SourceName:      t.SourceName,
		}
		if err := tx.Create(&translation).Error; err != nil {
			tx.Rollback()
			return core.NewServiceResponse(false, fmt.Sprintf("Translation create failed: %v", err), 500, nil)
		}
	}

	tx.Commit()

	// Çevirilerle birlikte dön
	db.Preload("Translations").First(&verse, verse.ID)
	return core.NewServiceResponse(true, "Verse created", 201, verse)
}

// Update mevcut ayeti günceller.
func (s *DailyVerseService) Update(ctx context.Context, id uint, req *verseRequest.UpdateDailyVerseRequest) *core.ServiceResponse {
	db := database.DB.WithContext(ctx)

	var verse models.DailyVerse
	if err := db.First(&verse, id).Error; err != nil {
		return core.NewServiceResponse(false, "Verse not found", 404, nil)
	}

	tx := db.Begin()

	updates := map[string]interface{}{}
	if req.SurahNumber != nil {
		updates["surah_number"] = *req.SurahNumber
	}
	if req.AyahNumber != nil {
		updates["ayah_number"] = *req.AyahNumber
	}
	if req.ArabicText != nil {
		updates["arabic_text"] = *req.ArabicText
	}
	if req.DisplayDate != nil {
		if *req.DisplayDate == "" {
			updates["display_date"] = nil // Havuza geri al
		} else {
			updates["display_date"] = *req.DisplayDate
		}
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}

	if len(updates) > 0 {
		if err := tx.Model(&verse).Updates(updates).Error; err != nil {
			tx.Rollback()
			return core.NewServiceResponse(false, fmt.Sprintf("Update failed: %v", err), 500, nil)
		}
	}

	// Çeviriler verilmişse mevcut olanları sil ve yenilerini ekle
	if req.Translations != nil && len(req.Translations) > 0 {
		tx.Where("daily_verse_id = ?", id).Delete(&models.DailyVerseTranslation{})
		for _, t := range req.Translations {
			translation := models.DailyVerseTranslation{
				DailyVerseID:    id,
				LanguageCode:    t.LanguageCode,
				TranslationText: t.TranslationText,
				SourceName:      t.SourceName,
			}
			if err := tx.Create(&translation).Error; err != nil {
				tx.Rollback()
				return core.NewServiceResponse(false, fmt.Sprintf("Translation update failed: %v", err), 500, nil)
			}
		}
	}

	tx.Commit()

	db.Preload("Translations").First(&verse, id)
	return core.NewServiceResponse(true, "Verse updated", 200, verse)
}

// Delete ayeti siler.
func (s *DailyVerseService) Delete(ctx context.Context, id uint) *core.ServiceResponse {
	db := database.DB.WithContext(ctx)

	var verse models.DailyVerse
	if err := db.First(&verse, id).Error; err != nil {
		return core.NewServiceResponse(false, "Verse not found", 404, nil)
	}

	tx := db.Begin()
	tx.Where("daily_verse_id = ?", id).Delete(&models.DailyVerseTranslation{})
	if err := tx.Delete(&verse).Error; err != nil {
		tx.Rollback()
		return core.NewServiceResponse(false, "Delete failed", 500, nil)
	}
	tx.Commit()

	return core.NewServiceResponse(true, "Verse deleted", 200, nil)
}
