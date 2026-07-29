package chat

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rampesna/kuran-ai-gateway/core"
	"github.com/rampesna/kuran-ai-gateway/middleware"
	chatRequest "github.com/rampesna/kuran-ai-gateway/requests/chat"
	"github.com/rampesna/kuran-ai-gateway/services"
)

// MessageHandler mesaj gönderme (SSE stream) handler.
type MessageHandler struct {
	orchestration *services.ChatOrchestrationService
	usageSvc      *services.UsageService
}

// NewMessageHandler MessageHandler oluşturur.
func NewMessageHandler(orchestration *services.ChatOrchestrationService, usageSvc *services.UsageService) *MessageHandler {
	return &MessageHandler{orchestration: orchestration, usageSvc: usageSvc}
}

// SendMessage POST /api/v1/chat/message - mesaj gönderir, cevabı SSE ile stream eder.
func (h *MessageHandler) SendMessage(c *gin.Context) {
	userIDVal, ok := c.Get(string(middleware.ContextKeyUserID))
	if !ok {
		c.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	userID := userIDVal.(uint)

	if err := h.usageSvc.CheckAndLogMessage(c.Request.Context(), userID); err != nil {
		if le, ok := services.IsLimitExceededError(err); ok {
			c.JSON(http.StatusPaymentRequired, core.NewServiceResponse(false, le.Message, 402, gin.H{
				"code":      le.Code,
				"limitType": le.LimitType,
				"used":      le.Used,
				"limit":     le.Limit,
				"packageId": le.PackageID,
			}))
			return
		}
		c.JSON(http.StatusInternalServerError, core.NewServiceResponse(false, "Kullanım kontrolü başarısız", 500, nil))
		return
	}

	var req chatRequest.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "Invalid request", 400, nil))
		return
	}
	if req.Content == "" {
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "content required", 400, nil))
		return
	}
	sourceType := req.SourceType
	if sourceType == "kuran" {
		sourceType = "meal"
	}
	validSources := map[string]bool{"meal": true, "tefsir": true, "hadis": true}
	if !validSources[sourceType] {
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "sourceType must be meal/kuran, tefsir, or hadis", 400, nil))
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	writeSSE := func(obj map[string]interface{}) {
		b, _ := json.Marshal(obj)
		c.Writer.Write([]byte("data: "))
		c.Writer.Write(b)
		c.Writer.Write([]byte("\n\n"))
		c.Writer.Flush()
	}

	sessionID, _, err := h.orchestration.ProcessMessage(c.Request.Context(), services.AuthTypeUser, userID, req.SessionID, sourceType, req.Content,
		func(sid uint) error {
			writeSSE(map[string]interface{}{"type": "started", "sessionId": sid})
			return nil
		},
		func(chunk string) error {
			writeSSE(map[string]interface{}{"type": "chunk", "content": chunk})
			return nil
		},
		func(content string) error {
			writeSSE(map[string]interface{}{"type": "completed", "content": content})
			return nil
		})
	if err != nil {
		log.Printf("[Chat] ProcessMessage error: %v", err)
		writeSSE(map[string]interface{}{"type": "error", "message": err.Error()})
		return
	}
	writeSSE(map[string]interface{}{"type": "done", "sessionId": sessionID})
}
