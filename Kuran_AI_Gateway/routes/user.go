package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	userauth "github.com/rampesna/kuran-ai-gateway/handlers/user"
	"github.com/rampesna/kuran-ai-gateway/middleware"
	"github.com/rampesna/kuran-ai-gateway/services"
)

// RegisterUserRoutes user API rotalarını kaydeder.
func RegisterUserRoutes(rg *gin.RouterGroup, userService *services.UserService, tokenService *services.TokenService, passkeyService *services.PasskeyService) {
	authHandler := userauth.NewAuthHandler(userService, tokenService)
	passkeyHandler := userauth.NewPasskeyHandler(passkeyService)
	profileHandler := userauth.NewProfileHandler(userService)

	user := rg.Group("/user")
	{
		user.POST("/auth/register", authHandler.Register)
		user.POST("/auth/login", authHandler.Login)
		user.POST("/auth/login/google", authHandler.LoginWithGoogle)
		user.POST("/auth/login/apple", authHandler.LoginWithApple)
		user.POST(
			"/auth/login/apple/mobile/start",
			middleware.RateLimiter(10, 10*time.Minute),
			authHandler.StartAppleMobileAuth,
		)
		user.POST(
			"/auth/login/mobile/finish",
			middleware.RateLimiter(20, 10*time.Minute),
			authHandler.FinishMobileAuth,
		)
		user.POST("/auth/passkey/login/options", passkeyHandler.BeginLogin)
		user.POST("/auth/passkey/login/finish", passkeyHandler.FinishLogin)
		user.POST("/auth/refresh", authHandler.RefreshToken)
	}

	userProtected := rg.Group("/user")
	userProtected.Use(middleware.AuthMiddleware(tokenService), middleware.RequireUser())
	{
		userProtected.GET("/auth/me", authHandler.Me)
		userProtected.POST("/auth/logout", authHandler.Logout)
		userProtected.GET("/auth/passkey", passkeyHandler.Status)
		userProtected.POST("/auth/passkey/register/options", passkeyHandler.BeginRegistration)
		userProtected.POST("/auth/passkey/register/finish", passkeyHandler.FinishRegistration)
		userProtected.DELETE("/auth/passkey", passkeyHandler.DeleteAll)
		userProtected.PATCH("/profile/username", profileHandler.UpdateUsername)
		userProtected.PATCH("/profile/password", profileHandler.UpdatePassword)
		userProtected.POST("/profile/avatar", profileHandler.UploadAvatar)
		userProtected.DELETE("/profile/avatar", profileHandler.RemoveAvatar)
		userProtected.DELETE("/profile/account", profileHandler.DeleteAccount)
	}
}
