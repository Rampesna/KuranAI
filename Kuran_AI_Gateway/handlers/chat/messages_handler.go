package chat

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rampesna/kuran-ai-gateway/core"
	"github.com/rampesna/kuran-ai-gateway/middleware"
	"github.com/rampesna/kuran-ai-gateway/services"
)

// MessagesHandler chat mesajları handler.
type MessagesHandler struct {
	sessionSvc *services.ChatSessionService
	messageSvc *services.ChatMessageService
}

// NewMessagesHandler MessagesHandler oluşturur.
func NewMessagesHandler(sessionSvc *services.ChatSessionService, messageSvc *services.ChatMessageService) *MessagesHandler {
	return &MessagesHandler{sessionSvc: sessionSvc, messageSvc: messageSvc}
}

// ListMessages GET /api/v1/chat/sessions/:id/messages - mesaj listesi (pagination).
func (h *MessagesHandler) ListMessages(c *gin.Context) {
	userIDVal, ok := c.Get(string(middleware.ContextKeyUserID))
	if !ok {
		c.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	userID := userIDVal.(uint)

	idStr := c.Param("id")
	sessionID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "Invalid session id", 400, nil))
		return
	}

	_, err = h.sessionSvc.GetByID(c.Request.Context(), uint(sessionID), services.AuthTypeUser, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, core.NewServiceResponse(false, "Session not found", 404, nil))
		return
	}

	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	messages, total, err := h.messageSvc.ListBySessionID(c.Request.Context(), uint(sessionID), pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.NewServiceResponse(false, err.Error(), 500, nil))
		return
	}

	c.JSON(http.StatusOK, core.NewServiceResponse(true, "ok", 200, map[string]interface{}{
		"items": messages,
		"total": total,
	}))
}
