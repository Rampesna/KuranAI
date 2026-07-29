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
	dashboardHandler := managerHandler.NewDashboardHandler(services.NewManagerDashboardService())

	mgr := rg.Group("/manager")
	{
		mgr.POST("/auth/login", authHandler.Login)
	}

	mgrProtected := rg.Group("/manager")
	mgrProtected.Use(middleware.AuthMiddleware(tokenService), middleware.RequireManager())
	{
		mgrProtected.GET("/auth/me", authHandler.Me)
		mgrProtected.GET("/dashboard/overview", dashboardHandler.Overview)
		mgrProtected.GET("/users", dashboardHandler.Users)
		mgrProtected.GET("/users/:id", dashboardHandler.UserDetail)
		mgrProtected.GET("/users/:id/conversations", dashboardHandler.Conversations)
		mgrProtected.GET("/conversations/:id/messages", dashboardHandler.ConversationMessages)
		mgrProtected.GET("/purchases", dashboardHandler.Purchases)
		mgrProtected.GET("/questions", dashboardHandler.Questions)
	}
}
