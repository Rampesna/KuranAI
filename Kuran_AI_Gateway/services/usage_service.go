package services

import (
	"context"
	"errors"
	"time"

	"github.com/rampesna/kuran-ai-gateway/database"
	"github.com/rampesna/kuran-ai-gateway/models"
)

// LimitExceededError sınır aşımı hatası.
type LimitExceededError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	LimitType  string `json:"limitType"` // "message" | "tts"
	Used       int    `json:"used"`
	Limit      int    `json:"limit"`
	PackageID  *uint  `json:"packageId,omitempty"`
}

func (e *LimitExceededError) Error() string {
	return e.Message
}

// UsageService günlük kullanım takibi.
type UsageService struct {
	subSvc *SubscriptionService
}

// NewUsageService UsageService oluşturur.
func NewUsageService(subSvc *SubscriptionService) *UsageService {
	return &UsageService{subSvc: subSvc}
}

// CheckAndLogMessage mesaj gönderme limitini kontrol eder, aşılmamışsa loglar.
func (u *UsageService) CheckAndLogMessage(ctx context.Context, userID uint) error {
	used, limit, pkgID, err := u.getDailyUsage(ctx, userID, models.UsageActionMessage)
	if err != nil {
		return err
	}
	if used >= limit {
		return &LimitExceededError{
			Code:      "MESSAGE_LIMIT_EXCEEDED",
			Message:   "Günlük soru sorma limitinize ulaştınız.",
			LimitType: "message",
			Used:      used,
			Limit:     limit,
			PackageID: pkgID,
		}
	}
	return u.logUsage(ctx, userID, models.UsageActionMessage)
}

// CheckAndLogTTS sesli okuma limitini kontrol eder, aşılmamışsa loglar.
func (u *UsageService) CheckAndLogTTS(ctx context.Context, userID uint) error {
	used, limit, pkgID, err := u.getDailyUsage(ctx, userID, models.UsageActionTTS)
	if err != nil {
		return err
	}
	if used >= limit {
		return &LimitExceededError{
			Code:      "TTS_LIMIT_EXCEEDED",
			Message:   "Günlük sesli okuma limitinize ulaştınız.",
			LimitType: "tts",
			Used:      used,
			Limit:     limit,
			PackageID: pkgID,
		}
	}
	return u.logUsage(ctx, userID, models.UsageActionTTS)
}

func (u *UsageService) getDailyUsage(ctx context.Context, userID uint, action models.UsageAction) (used int, limit int, packageID *uint, err error) {
	sub, pkg, _ := u.subSvc.GetActiveSubscription(ctx, userID)

	switch action {
	case models.UsageActionMessage:
		if sub == nil || pkg == nil {
			limit = 3 // paket yok: günde 3 soru
		} else if pkg.DailyQuestionLimit == nil {
			limit = 999999
		} else {
			limit = *pkg.DailyQuestionLimit
			p := sub.PackageID
			packageID = &p
		}
	case models.UsageActionTTS:
		if sub == nil || pkg == nil {
			limit = 0 // paket yok: sesli okuma yok
		} else if pkg.DailyTTSLimit == nil {
			limit = 999999
		} else {
			limit = *pkg.DailyTTSLimit
			p := sub.PackageID
			packageID = &p
		}
	}

	todayStart := time.Now().UTC().Truncate(24 * time.Hour)
	var count int64
	database.DB.WithContext(ctx).Model(&models.UsageLog{}).
		Where("user_id = ? AND action = ? AND created_at >= ?", userID, action, todayStart).
		Count(&count)
	used = int(count)
	return used, limit, packageID, nil
}

func (u *UsageService) logUsage(ctx context.Context, userID uint, action models.UsageAction) error {
	return database.DB.WithContext(ctx).Create(&models.UsageLog{
		UserID: userID,
		Action: action,
	}).Error
}

// GetDailyUsage kullanıcının bugünkü kullanımını döner.
func (u *UsageService) GetDailyUsage(ctx context.Context, userID uint) (messageUsed, messageLimit int, ttsUsed, ttsLimit int, packageID *uint, hasPackage bool, err error) {
	sub, pkg, _ := u.subSvc.GetActiveSubscription(ctx, userID)
	hasPackage = sub != nil && pkg != nil

	if hasPackage {
		if pkg.DailyQuestionLimit != nil {
			messageLimit = *pkg.DailyQuestionLimit
		} else {
			messageLimit = 999999
		}
		if pkg.DailyTTSLimit != nil {
			ttsLimit = *pkg.DailyTTSLimit
		} else {
			ttsLimit = 999999
		}
		p := sub.PackageID
		packageID = &p
	} else {
		messageLimit = 3
		ttsLimit = 0
	}

	todayStart := time.Now().UTC().Truncate(24 * time.Hour)

	var msgCount, ttsCount int64
	database.DB.WithContext(ctx).Model(&models.UsageLog{}).
		Where("user_id = ? AND action = ? AND created_at >= ?", userID, models.UsageActionMessage, todayStart).
		Count(&msgCount)
	database.DB.WithContext(ctx).Model(&models.UsageLog{}).
		Where("user_id = ? AND action = ? AND created_at >= ?", userID, models.UsageActionTTS, todayStart).
		Count(&ttsCount)

	return int(msgCount), messageLimit, int(ttsCount), ttsLimit, packageID, hasPackage, nil
}

// IsLimitExceededError LimitExceededError mı kontrol eder.
func IsLimitExceededError(err error) (*LimitExceededError, bool) {
	var le *LimitExceededError
	if errors.As(err, &le) {
		return le, true
	}
	return nil, false
}
