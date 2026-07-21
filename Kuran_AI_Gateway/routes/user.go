package routes

import (
	"github.com/gin-gonic/gin"
	userauth "github.com/rampesna/kuran-ai-gateway/handlers/user"
	"github.com/rampesna/kuran-ai-gateway/middleware"
	"github.com/rampesna/kuran-ai-gateway/services"
)

// RegisterUserRoutes user API rotalarını kaydeder.
func RegisterUserRoutes(rg *gin.RouterGroup, userService *services.UserService, tokenService *services.TokenService) {
	authHandler := userauth.NewAuthHandler(userService, tokenService)
	profileHandler := userauth.NewProfileHandler(userService)

	user := rg.Group("/user")
	{
		user.POST("/auth/register", authHandler.Register)
		user.POST("/auth/login", authHandler.Login)
		user.POST("/auth/login/google", authHandler.LoginWithGoogle)
		user.POST("/auth/login/apple", authHandler.LoginWithApple)
		user.POST("/auth/refresh", authHandler.RefreshToken)
	}

	userProtected := rg.Group("/user")
	userProtected.Use(middleware.AuthMiddleware(tokenService), middleware.RequireUser())
	{
		userProtected.GET("/auth/me", authHandler.Me)
		userProtected.POST("/auth/logout", authHandler.Logout)
		userProtected.PATCH("/profile/username", profileHandler.UpdateUsername)
		userProtected.PATCH("/profile/password", profileHandler.UpdatePassword)
		userProtected.POST("/profile/avatar", profileHandler.UploadAvatar)
		userProtected.DELETE("/profile/avatar", profileHandler.RemoveAvatar)
		userProtected.DELETE("/profile/account", profileHandler.DeleteAccount)
	}
}
