package services

import (
	"context"

	"github.com/rampesna/kuran-ai-gateway/database"
	"github.com/rampesna/kuran-ai-gateway/models"
)

const (
	AuthTypeUser   = "user"
	AuthTypeManager = "manager"
)

// ChatSessionService chat oturumu servisi.
type ChatSessionService struct{}

// NewChatSessionService ChatSessionService oluşturur.
func NewChatSessionService() *ChatSessionService {
	return &ChatSessionService{}
}

// Create yeni chat oturumu oluşturur. firstMessage boşsa güncellenir ilk mesajda.
func (s *ChatSessionService) Create(ctx context.Context, authType string, authID uint, firstMessage string) (*models.ChatSession, error) {
	session := &models.ChatSession{
		AuthType:     authType,
		AuthID:       authID,
		FirstMessage: firstMessage,
	}
	if err := database.DB.WithContext(ctx).Create(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

// GetByID ID ile oturum getirir.
func (s *ChatSessionService) GetByID(ctx context.Context, id uint, authType string, authID uint) (*models.ChatSession, error) {
	var session models.ChatSession
	err := database.DB.WithContext(ctx).
		Where("id = ? AND auth_type = ? AND auth_id = ?", id, authType, authID).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// List oturumları sayfalı listeler. pageIndex 0-based, pageSize varsayılan 20, max 100.
func (s *ChatSessionService) List(ctx context.Context, authType string, authID uint, pageIndex, pageSize int) ([]models.ChatSession, int64, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if pageIndex < 0 {
		pageIndex = 0
	}
	offset := pageIndex * pageSize

	var total int64
	if err := database.DB.WithContext(ctx).Model(&models.ChatSession{}).
		Where("auth_type = ? AND auth_id = ?", authType, authID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var sessions []models.ChatSession
	err := database.DB.WithContext(ctx).
		Where("auth_type = ? AND auth_id = ?", authType, authID).
		Order("updated_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&sessions).Error
	return sessions, total, err
}

// UpdateFirstMessage first_message günceller (boşsa).
func (s *ChatSessionService) UpdateFirstMessage(ctx context.Context, id uint, firstMessage string) error {
	return database.DB.WithContext(ctx).Model(&models.ChatSession{}).
		Where("id = ? AND (first_message = '' OR first_message IS NULL)", id).
		Update("first_message", firstMessage).Error
}
