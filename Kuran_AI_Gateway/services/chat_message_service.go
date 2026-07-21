package services

import (
	"context"

	"github.com/pgvector/pgvector-go"
	"github.com/rampesna/kuran-ai-gateway/database"
	"github.com/rampesna/kuran-ai-gateway/models"
)

// ChatMessageService chat mesaj servisi.
type ChatMessageService struct{}

// NewChatMessageService ChatMessageService oluşturur.
func NewChatMessageService() *ChatMessageService {
	return &ChatMessageService{}
}

// Create yeni mesaj oluşturur.
func (s *ChatMessageService) Create(ctx context.Context, sessionID uint, role, content, sourceType string, embedding []float32) (*models.ChatMessage, error) {
	msg := &models.ChatMessage{
		ChatSessionID: sessionID,
		Role:          role,
		Content:       content,
		SourceType:    sourceType,
	}
	if len(embedding) > 0 {
		v := pgvector.NewVector(embedding)
		msg.Embedding = &v
	}
	if err := database.DB.WithContext(ctx).Create(msg).Error; err != nil {
		return nil, err
	}
	return msg, nil
}

// UpdateEmbedding mesajın embedding alanını günceller (arka plan embedding için).
func (s *ChatMessageService) UpdateEmbedding(ctx context.Context, messageID uint, embedding []float32) error {
	if len(embedding) == 0 {
		return nil
	}
	v := pgvector.NewVector(embedding)
	return database.DB.WithContext(ctx).Model(&models.ChatMessage{}).
		Where("id = ?", messageID).
		Update("embedding", v).Error
}

// ListBySessionID oturuma ait mesajları sayfalı listeler.
func (s *ChatMessageService) ListBySessionID(ctx context.Context, sessionID uint, pageIndex, pageSize int) ([]models.ChatMessage, int64, error) {
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
	if err := database.DB.WithContext(ctx).Model(&models.ChatMessage{}).
		Where("chat_session_id = ?", sessionID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var messages []models.ChatMessage
	err := database.DB.WithContext(ctx).
		Where("chat_session_id = ?", sessionID).
		Order("created_at ASC").
		Offset(offset).Limit(pageSize).
		Find(&messages).Error
	return messages, total, err
}
