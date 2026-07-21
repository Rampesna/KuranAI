package manager

import (
	"github.com/gin-gonic/gin"
	"github.com/talhacan/gateway/core"
	managerRequest "github.com/talhacan/gateway/requests/manager"
	"github.com/talhacan/gateway/services"
)

// ManagerHandler manager CRUD handler; Request → ToDTO → ManagerService.
type ManagerHandler struct {
	managerService *services.ManagerService
}

// NewManagerHandler ManagerHandler oluşturur.
func NewManagerHandler(managerService *services.ManagerService) *ManagerHandler {
	return &ManagerHandler{managerService: managerService}
}

// GetAll ManagerService.Index (tümü; query opsiyonel).
func (managerHandler *ManagerHandler) GetAll(ginContext *gin.Context) {
	core.WithValidatedQuery(ginContext, func(managerIndexRequest *managerRequest.ManagerIndexRequest) *core.ServiceResponse {
		return managerHandler.managerService.Index(ginContext.Request.Context(), managerIndexRequest)
	})
}

// GetByID ManagerService.GetByID (id URI'den).
func (managerHandler *ManagerHandler) GetByID(ginContext *gin.Context) {
	core.WithValidatedUri(ginContext, func(managerGetByIDRequest *managerRequest.ManagerGetByIDRequest) *core.ServiceResponse {
		return managerHandler.managerService.GetByID(ginContext.Request.Context(), managerGetByIDRequest)
	})
}

// Index ManagerService.Index.
func (managerHandler *ManagerHandler) Index(ginContext *gin.Context) {
	core.WithValidatedQuery(ginContext, func(managerIndexRequest *managerRequest.ManagerIndexRequest) *core.ServiceResponse {
		return managerHandler.managerService.Index(ginContext.Request.Context(), managerIndexRequest)
	})
}

// Create ManagerService.Create.
func (managerHandler *ManagerHandler) Create(ginContext *gin.Context) {
	core.WithValidatedJSON(ginContext, func(managerCreateRequest *managerRequest.ManagerCreateRequest) *core.ServiceResponse {
		return managerHandler.managerService.Create(ginContext.Request.Context(), managerCreateRequest)
	})
}

// Update ManagerService.Update (id URI + body JSON).
func (managerHandler *ManagerHandler) Update(ginContext *gin.Context) {
	core.WithValidatedUriAndJSON(ginContext, func(managerUpdateRequest *managerRequest.ManagerUpdateRequest) *core.ServiceResponse {
		return managerHandler.managerService.Update(ginContext.Request.Context(), managerUpdateRequest)
	})
}

// Delete ManagerService.Delete (id URI'den).
func (managerHandler *ManagerHandler) Delete(ginContext *gin.Context) {
	core.WithValidatedUri(ginContext, func(managerDeleteRequest *managerRequest.ManagerDeleteRequest) *core.ServiceResponse {
		return managerHandler.managerService.Delete(ginContext.Request.Context(), managerDeleteRequest)
	})
}
