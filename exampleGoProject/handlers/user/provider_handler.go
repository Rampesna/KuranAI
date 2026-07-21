package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/talhacan/gateway/core"
	"github.com/talhacan/gateway/database"
	"github.com/talhacan/gateway/middleware"
	"github.com/talhacan/gateway/models"
)

// ProviderHandler AI provider ve model listesi; DB'den okur.
type ProviderHandler struct{}

// NewProviderHandler ProviderHandler oluşturur.
func NewProviderHandler() *ProviderHandler {
	return &ProviderHandler{}
}

// GetAllProviders tüm AI provider'ları döner (protected).
func (h *ProviderHandler) GetAllProviders(ginContext *gin.Context) {
	core.WithAuthID(ginContext, string(middleware.ContextKeyUserID), func(_ uint) *core.ServiceResponse {
		var list []models.AiProvider
		if err := database.DB.Order("id").Find(&list).Error; err != nil {
			return core.NewServiceResponse(false, "Failed to load providers", 500, nil)
		}
		return core.NewServiceResponse(true, "OK", 200, list)
	})
}

// GetModelsByProviderID verilen provider id'ye ait modelleri döner (protected).
func (h *ProviderHandler) GetModelsByProviderID(ginContext *gin.Context) {
	idStr := ginContext.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ginContext.JSON(400, core.NewServiceResponse(false, "Invalid provider id", 400, nil))
		return
	}
	core.WithAuthID(ginContext, string(middleware.ContextKeyUserID), func(_ uint) *core.ServiceResponse {
		var list []models.AiModel
		if err := database.DB.Where("provider_id = ?", uint(id)).Order("id").Find(&list).Error; err != nil {
			return core.NewServiceResponse(false, "Failed to load models", 500, nil)
		}
		return core.NewServiceResponse(true, "OK", 200, list)
	})
}
