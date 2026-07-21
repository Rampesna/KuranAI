package manager

import (
	"github.com/gin-gonic/gin"
	"github.com/talhacan/gateway/core"
	"github.com/talhacan/gateway/middleware"
	managerRequest "github.com/talhacan/gateway/requests/manager"
	"github.com/talhacan/gateway/services"
)

// AuthHandler manager auth handler; ManagerService + TokenService kullanır.
type AuthHandler struct {
	managerService *services.ManagerService
	tokenService   *services.TokenService
}

// NewAuthHandler AuthHandler oluşturur.
func NewAuthHandler(managerService *services.ManagerService, tokenService *services.TokenService) *AuthHandler {
	return &AuthHandler{managerService: managerService, tokenService: tokenService}
}

// Login ManagerService.Login; DeviceInfo core tarafından request'e set edilir; ToLoginDTO() sadece request kullanır.
func (authHandler *AuthHandler) Login(ginContext *gin.Context) {
	core.WithValidatedJSONAndDeviceInfo(ginContext, func(loginRequest *managerRequest.LoginRequest) *core.ServiceResponse {
		return authHandler.managerService.Login(ginContext.Request.Context(), loginRequest)
	})
}

// RefreshToken TokenService.RefreshTokens; request → ToRefreshTokenDTO() (request üzerinden).
func (authHandler *AuthHandler) RefreshToken(ginContext *gin.Context) {
	core.WithValidatedJSONAndDeviceInfo(ginContext, func(refreshTokenRequest *managerRequest.RefreshTokenRequest) *core.ServiceResponse {
		return authHandler.tokenService.RefreshTokens(ginContext.Request.Context(), refreshTokenRequest)
	})
}

// Me mevcut manager; WithAuthID → ManagerService.Me.
func (authHandler *AuthHandler) Me(ginContext *gin.Context) {
	core.WithAuthID(ginContext, string(middleware.ContextKeyManagerID), func(managerID uint) *core.ServiceResponse {
		return authHandler.managerService.Me(ginContext.Request.Context(), managerID)
	})
}

// Logout token core.WithAuthToken ile alınır, TokenService.RevokeAccessToken çağrılır.
func (authHandler *AuthHandler) Logout(ginContext *gin.Context) {
	core.WithAuthToken(ginContext, func(accessToken string) *core.ServiceResponse {
		return authHandler.tokenService.RevokeAccessToken(ginContext.Request.Context(), accessToken)
	})
}
