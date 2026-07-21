package user

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rampesna/kuran-ai-gateway/core"
	"github.com/rampesna/kuran-ai-gateway/middleware"
	"github.com/rampesna/kuran-ai-gateway/services"
)

const passkeyTransactionHeader = "X-Passkey-Transaction"

type PasskeyHandler struct {
	service *services.PasskeyService
}

func NewPasskeyHandler(service *services.PasskeyService) *PasskeyHandler {
	return &PasskeyHandler{service: service}
}

func (h *PasskeyHandler) BeginLogin(c *gin.Context) {
	writePasskeyResponse(c, h.service.BeginLogin(c.Request.Context()))
}

func (h *PasskeyHandler) FinishLogin(c *gin.Context) {
	transactionID := strings.TrimSpace(c.GetHeader(passkeyTransactionHeader))
	if transactionID == "" {
		writePasskeyResponse(c, core.NewServiceResponse(false, "Passkey transaction header is required", http.StatusBadRequest, nil))
		return
	}
	writePasskeyResponse(c, h.service.FinishLogin(c.Request.Context(), transactionID, c.Request, core.DeviceInfoFromContext(c)))
}

func (h *PasskeyHandler) BeginRegistration(c *gin.Context) {
	withPasskeyUser(c, func(userID uint) *core.ServiceResponse {
		return h.service.BeginRegistration(c.Request.Context(), userID)
	})
}

func (h *PasskeyHandler) FinishRegistration(c *gin.Context) {
	transactionID := strings.TrimSpace(c.GetHeader(passkeyTransactionHeader))
	if transactionID == "" {
		writePasskeyResponse(c, core.NewServiceResponse(false, "Passkey transaction header is required", http.StatusBadRequest, nil))
		return
	}
	withPasskeyUser(c, func(userID uint) *core.ServiceResponse {
		return h.service.FinishRegistration(c.Request.Context(), userID, transactionID, c.Request)
	})
}

func (h *PasskeyHandler) Status(c *gin.Context) {
	withPasskeyUser(c, func(userID uint) *core.ServiceResponse {
		return h.service.Status(c.Request.Context(), userID)
	})
}

func (h *PasskeyHandler) DeleteAll(c *gin.Context) {
	withPasskeyUser(c, func(userID uint) *core.ServiceResponse {
		return h.service.DeleteAll(c.Request.Context(), userID)
	})
}

func withPasskeyUser(c *gin.Context, callback func(uint) *core.ServiceResponse) {
	value, exists := c.Get(string(middleware.ContextKeyUserID))
	userID, ok := value.(uint)
	if !exists || !ok {
		writePasskeyResponse(c, core.NewServiceResponse(false, "Unauthorized", http.StatusUnauthorized, nil))
		return
	}
	writePasskeyResponse(c, callback(userID))
}

func writePasskeyResponse(c *gin.Context, response *core.ServiceResponse) {
	c.JSON(response.StatusCode, response)
}
