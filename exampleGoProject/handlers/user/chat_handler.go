package user

import (
	"github.com/gin-gonic/gin"
	"github.com/talhacan/gateway/ai_gateway_services"
	"github.com/talhacan/gateway/ai_gateway_services/types"
	"github.com/talhacan/gateway/core"
	"github.com/talhacan/gateway/database"
	"github.com/talhacan/gateway/middleware"
	"github.com/talhacan/gateway/models"
	userRequest "github.com/talhacan/gateway/requests/user"
)

// ChatHandler user chat handler; session/history + AI gateway.
type ChatHandler struct{}

// NewChatHandler ChatHandler oluşturur.
func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}

// SendMessage session varsa geçmişi yükler, messages ile gateway'e gönderir; user + assistant mesajlarını DB'ye yazar.
func (chatHandler *ChatHandler) SendMessage(ginContext *gin.Context) {
	core.WithValidatedJSONAndUser(ginContext, string(middleware.ContextKeyUserID), func(req *userRequest.SendMessageRequest) *core.ServiceResponse {
		ctx := ginContext.Request.Context()
		userID := *req.UserID

		var session models.ChatSession
		if req.SessionID != nil && *req.SessionID > 0 {
			if err := database.DB.Where("id = ? AND user_id = ?", *req.SessionID, userID).First(&session).Error; err != nil {
				return core.NewServiceResponse(false, "Session not found", 404, nil)
			}
		} else {
			session = models.ChatSession{UserID: userID}
			if err := database.DB.Create(&session).Error; err != nil {
				return core.NewServiceResponse(false, "Failed to create session", 500, nil)
			}
		}

		var history []models.ChatMessage
		if err := database.DB.Where("session_id = ?", session.ID).Order("id").Find(&history).Error; err != nil {
			return core.NewServiceResponse(false, "Failed to load history", 500, nil)
		}
		messages := make([]types.Message, 0, len(history)+1)
		for _, m := range history {
			messages = append(messages, types.Message{Role: m.Role, Content: m.Content})
		}
		messages = append(messages, types.Message{Role: "user", Content: req.Prompt})

		resp, err := ai_gateway_services.SendMessage(ctx, req.Provider, req.Prompt, req.Model, req.ApiKey, messages)
		if err != nil {
			return core.NewServiceResponse(false, err.Error(), 500, nil)
		}

		if err := database.DB.Create(&models.ChatMessage{SessionID: session.ID, Role: "user", Content: req.Prompt}).Error; err != nil {
			return core.NewServiceResponse(false, "Failed to save message", 500, nil)
		}
		if err := database.DB.Create(&models.ChatMessage{SessionID: session.ID, Role: "assistant", Content: resp.Content}).Error; err != nil {
			return core.NewServiceResponse(false, "Failed to save response", 500, nil)
		}

		out := map[string]interface{}{"content": resp.Content, "sessionId": session.ID}
		return core.NewServiceResponse(true, "OK", 200, out)
	})
}
