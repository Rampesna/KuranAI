package verse

// TranslationInput tek bir çeviri girdisi.
type TranslationInput struct {
	LanguageCode    string `json:"languageCode" binding:"required"`
	TranslationText string `json:"translationText" binding:"required"`
	SourceName      string `json:"sourceName" binding:"required"`
}

// CreateDailyVerseRequest yeni günün ayeti oluşturma isteği.
type CreateDailyVerseRequest struct {
	SurahNumber  int                `json:"surahNumber" binding:"required,min=1"`
	AyahNumber   int                `json:"ayahNumber" binding:"required,min=1"`
	ArabicText   string             `json:"arabicText" binding:"required"`
	DisplayDate  *string            `json:"displayDate"`  // "2026-02-09" formatında, nil = havuzda
	IsActive     *bool              `json:"isActive"`     // nil ise default true
	SortOrder    *int               `json:"sortOrder"`
	Translations []TranslationInput `json:"translations" binding:"required,min=1,dive"`
}

// UpdateDailyVerseRequest güncelleme isteği.
type UpdateDailyVerseRequest struct {
	SurahNumber  *int               `json:"surahNumber"`
	AyahNumber   *int               `json:"ayahNumber"`
	ArabicText   *string            `json:"arabicText"`
	DisplayDate  *string            `json:"displayDate"`
	IsActive     *bool              `json:"isActive"`
	SortOrder    *int               `json:"sortOrder"`
	Translations []TranslationInput `json:"translations"` // Verilirse mevcut çeviriler silinip yenisi eklenir
}
