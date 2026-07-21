package user

import (
	"github.com/gin-gonic/gin"
	"github.com/rampesna/kuran-ai-gateway/core"
	"github.com/rampesna/kuran-ai-gateway/middleware"
	userRequest "github.com/rampesna/kuran-ai-gateway/requests/user"
	"github.com/rampesna/kuran-ai-gateway/services"
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

// Register kullanıcı kaydı.
func (h *AuthHandler) Register(c *gin.Context) {
	core.WithValidatedJSON(c, func(req *userRequest.RegisterRequest) *core.ServiceResponse {
		return h.userService.Register(c.Request.Context(), req)
	})
}

// Login email/username + password ile giriş.
func (h *AuthHandler) Login(c *gin.Context) {
	core.WithValidatedJSONAndDeviceInfo(c, func(req *userRequest.LoginRequest) *core.ServiceResponse {
		return h.userService.Login(c.Request.Context(), req)
	})
}

// LoginWithGoogle Google OAuth ile giriş.
func (h *AuthHandler) LoginWithGoogle(c *gin.Context) {
	core.WithValidatedJSONAndDeviceInfo(c, func(req *userRequest.GoogleLoginRequest) *core.ServiceResponse {
		return h.userService.LoginWithGoogle(c.Request.Context(), req)
	})
}

// LoginWithApple Apple Sign in ile giriş.
func (h *AuthHandler) LoginWithApple(c *gin.Context) {
	core.WithValidatedJSONAndDeviceInfo(c, func(req *userRequest.AppleLoginRequest) *core.ServiceResponse {
		return h.userService.LoginWithApple(c.Request.Context(), req)
	})
}

// RefreshToken token yenileme.
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	core.WithValidatedJSONAndDeviceInfo(c, func(req *userRequest.RefreshTokenRequest) *core.ServiceResponse {
		return h.tokenService.RefreshTokens(c.Request.Context(), req)
	})
}

// Me mevcut kullanıcı.
func (h *AuthHandler) Me(c *gin.Context) {
	core.WithAuthID(c, string(middleware.ContextKeyUserID), func(userID uint) *core.ServiceResponse {
		return h.userService.Me(c.Request.Context(), userID)
	})
}

// Logout çıkış.
func (h *AuthHandler) Logout(c *gin.Context) {
	core.WithAuthToken(c, func(accessToken string) *core.ServiceResponse {
		return h.tokenService.RevokeAccessToken(c.Request.Context(), accessToken)
	})
}
