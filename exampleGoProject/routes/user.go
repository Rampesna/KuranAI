package routes

import (
	"github.com/gin-gonic/gin"
	userauth "github.com/talhacan/gateway/handlers/user"
	"github.com/talhacan/gateway/middleware"
	"github.com/talhacan/gateway/services"
)

// RegisterUserRoutes user API rotalarını kaydeder: register, login, refresh, logout, me.
func RegisterUserRoutes(routerGroup *gin.RouterGroup, userService *services.UserService, tokenService *services.TokenService) {
	authHandler := userauth.NewAuthHandler(userService, tokenService)

	// Public
	user := routerGroup.Group("/user")
	{
		user.POST("/auth/register", authHandler.Register)
		user.POST("/auth/login", authHandler.Login)
		user.POST("/auth/refresh", authHandler.RefreshToken)
	}

	// Protected (JWT + RequireUser)
	userProtected := routerGroup.Group("/user")
	userProtected.Use(middleware.AuthMiddleware(tokenService), middleware.RequireUser())
	{
		userProtected.GET("/auth/me", authHandler.Me)
		userProtected.POST("/auth/logout", authHandler.Logout)

		providerHandler := userauth.NewProviderHandler()
		userProtected.GET("/providers", providerHandler.GetAllProviders)
		userProtected.GET("/providers/:id/models", providerHandler.GetModelsByProviderID)

		chatHandler := userauth.NewChatHandler()
		userProtected.POST("/chat/sendMessage", chatHandler.SendMessage)
	}
}
