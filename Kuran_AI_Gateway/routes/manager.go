package routes

import (
	"github.com/gin-gonic/gin"
	managerHandler "github.com/rampesna/kuran-ai-gateway/handlers/manager"
	"github.com/rampesna/kuran-ai-gateway/middleware"
	"github.com/rampesna/kuran-ai-gateway/services"
)

// RegisterManagerRoutes manager API rotalarını kaydeder.
func RegisterManagerRoutes(rg *gin.RouterGroup, managerService *services.ManagerService, tokenService *services.TokenService) {
	authHandler := managerHandler.NewAuthHandler(managerService)

	mgr := rg.Group("/manager")
	{
		mgr.POST("/auth/login", authHandler.Login)
	}

	mgrProtected := rg.Group("/manager")
	mgrProtected.Use(middleware.AuthMiddleware(tokenService), middleware.RequireManager())
	{
		mgrProtected.GET("/auth/me", authHandler.Me)
	}
}
