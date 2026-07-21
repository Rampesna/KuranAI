package chat

import (
	"github.com/gin-gonic/gin"
	"github.com/rampesna/kuran-ai-gateway/core"
	"github.com/rampesna/kuran-ai-gateway/middleware"
	chatRequest "github.com/rampesna/kuran-ai-gateway/requests/chat"
	"github.com/rampesna/kuran-ai-gateway/services"
)

// SessionsHandler chat oturumları handler.
type SessionsHandler struct {
	sessionSvc *services.ChatSessionService
}

// NewSessionsHandler SessionsHandler oluşturur.
func NewSessionsHandler(sessionSvc *services.ChatSessionService) *SessionsHandler {
	return &SessionsHandler{sessionSvc: sessionSvc}
}

// ListSessions GET /api/v1/chat/sessions - oturum listesi (pagination).
func (h *SessionsHandler) ListSessions(c *gin.Context) {
	core.WithValidatedQueryAndUser[chatRequest.ListSessionsRequest](c, string(middleware.ContextKeyUserID), func(req *chatRequest.ListSessionsRequest) *core.ServiceResponse {
		sessions, total, err := h.sessionSvc.List(c.Request.Context(), services.AuthTypeUser, req.UserID, req.PageIndex, req.PageSize)
		if err != nil {
			return core.NewServiceResponse(false, err.Error(), 500, nil)
		}
		return core.NewServiceResponse(true, "ok", 200, map[string]interface{}{
			"items": sessions,
			"total": total,
		})
	})
}
