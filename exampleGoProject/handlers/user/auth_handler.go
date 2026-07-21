package user

import (
	"github.com/gin-gonic/gin"
	"github.com/talhacan/gateway/core"
	"github.com/talhacan/gateway/middleware"
	userRequest "github.com/talhacan/gateway/requests/user"
	"github.com/talhacan/gateway/services"
)

// AuthHandler user auth handler; UserService + TokenService kullanır.
type AuthHandler struct {
	userService  *services.UserService
	tokenService *services.TokenService
}

// NewAuthHandler AuthHandler oluşturur.
func NewAuthHandler(userService *services.UserService, tokenService *services.TokenService) *AuthHandler {
	return &AuthHandler{userService: userService, tokenService: tokenService}
}

// Register kullanıcı kaydı; request validate edilir, UserService.Register çağrılır.
func (authHandler *AuthHandler) Register(ginContext *gin.Context) {
	core.WithValidatedJSON(ginContext, func(registerRequest *userRequest.RegisterRequest) *core.ServiceResponse {
		return authHandler.userService.Register(ginContext.Request.Context(), registerRequest)
	})
}

// Login user login; DeviceInfo core.WithValidatedJSONAndDeviceInfo ile request'e set edilir; ToLoginDTO() sadece request kullanır.
func (authHandler *AuthHandler) Login(ginContext *gin.Context) {
	core.WithValidatedJSONAndDeviceInfo(ginContext, func(loginRequest *userRequest.LoginRequest) *core.ServiceResponse {
		return authHandler.userService.Login(ginContext.Request.Context(), loginRequest)
	})
}

// RefreshToken TokenService.RefreshTokens; request → ToRefreshTokenDTO() (request üzerinden).
func (authHandler *AuthHandler) RefreshToken(ginContext *gin.Context) {
	core.WithValidatedJSONAndDeviceInfo(ginContext, func(refreshTokenRequest *userRequest.RefreshTokenRequest) *core.ServiceResponse {
		return authHandler.tokenService.RefreshTokens(ginContext.Request.Context(), refreshTokenRequest)
	})
}

// Me mevcut user; WithAuthID → UserService.Me.
func (authHandler *AuthHandler) Me(ginContext *gin.Context) {
	core.WithAuthID(ginContext, string(middleware.ContextKeyUserID), func(userID uint) *core.ServiceResponse {
		return authHandler.userService.Me(ginContext.Request.Context(), userID)
	})
}

// Logout token core.WithAuthToken ile alınır, TokenService.RevokeAccessToken çağrılır.
func (authHandler *AuthHandler) Logout(ginContext *gin.Context) {
	core.WithAuthToken(ginContext, func(accessToken string) *core.ServiceResponse {
		return authHandler.tokenService.RevokeAccessToken(ginContext.Request.Context(), accessToken)
	})
}
