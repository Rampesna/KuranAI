package manager

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rampesna/kuran-ai-gateway/core"
	"github.com/rampesna/kuran-ai-gateway/middleware"
	"github.com/rampesna/kuran-ai-gateway/services"
)

// AuthHandler manager auth handler.
type AuthHandler struct {
	managerService *services.ManagerService
}

// NewAuthHandler AuthHandler oluşturur.
func NewAuthHandler(managerService *services.ManagerService) *AuthHandler {
	return &AuthHandler{managerService: managerService}
}

// Login manager giriş yapar.
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.ManagerLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if sr := core.ValidationErrorResponse(err); sr != nil {
			c.JSON(sr.StatusCode, sr)
			return
		}
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "Invalid request", 400, nil))
		return
	}
	sr := h.managerService.Login(c.Request.Context(), &req)
	c.JSON(sr.StatusCode, sr)
}

// Me mevcut manager bilgisini döner.
func (h *AuthHandler) Me(c *gin.Context) {
	managerIDVal, ok := c.Get(string(middleware.ContextKeyManagerID))
	if !ok {
		c.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	managerID, ok := managerIDVal.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	sr := h.managerService.Me(c.Request.Context(), managerID)
	c.JSON(sr.StatusCode, sr)
}
