package user

import (
	"github.com/gin-gonic/gin"
	"github.com/talhacan/gateway/core"
	"github.com/talhacan/gateway/requests/user"
	"github.com/talhacan/gateway/services"
)

// UserHandler user CRUD handler; Request → ToDTO → UserService.
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler UserHandler oluşturur.
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetAll UserService.Index (tümü; query opsiyonel).
func (userHandler *UserHandler) GetAll(ginContext *gin.Context) {
	core.WithValidatedQuery(ginContext, func(userIndexRequest *user.UserIndexRequest) *core.ServiceResponse {
		return userHandler.userService.Index(ginContext.Request.Context(), userIndexRequest)
	})
}

// GetByID UserService.GetByID (id URI'den).
func (userHandler *UserHandler) GetByID(ginContext *gin.Context) {
	core.WithValidatedUri(ginContext, func(userGetByIDRequest *user.UserGetByIDRequest) *core.ServiceResponse {
		return userHandler.userService.GetByID(ginContext.Request.Context(), userGetByIDRequest)
	})
}

// Index UserService.Index (sayfalı).
func (userHandler *UserHandler) Index(ginContext *gin.Context) {
	core.WithValidatedQuery(ginContext, func(userIndexRequest *user.UserIndexRequest) *core.ServiceResponse {
		return userHandler.userService.Index(ginContext.Request.Context(), userIndexRequest)
	})
}

// Create UserService.Create.
func (userHandler *UserHandler) Create(ginContext *gin.Context) {
	core.WithValidatedJSON(ginContext, func(userCreateRequest *user.UserCreateRequest) *core.ServiceResponse {
		return userHandler.userService.Create(ginContext.Request.Context(), userCreateRequest)
	})
}

// Update UserService.Update (id URI + body JSON).
func (userHandler *UserHandler) Update(ginContext *gin.Context) {
	core.WithValidatedUriAndJSON(ginContext, func(userUpdateRequest *user.UserUpdateRequest) *core.ServiceResponse {
		return userHandler.userService.Update(ginContext.Request.Context(), userUpdateRequest)
	})
}

// Delete UserService.Delete (id URI'den).
func (userHandler *UserHandler) Delete(ginContext *gin.Context) {
	core.WithValidatedUri(ginContext, func(userDeleteRequest *user.UserDeleteRequest) *core.ServiceResponse {
		return userHandler.userService.Delete(ginContext.Request.Context(), userDeleteRequest)
	})
}