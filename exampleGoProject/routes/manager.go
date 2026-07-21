package routes

import (
	"github.com/gin-gonic/gin"
	managerauth "github.com/talhacan/gateway/handlers/manager"
	userauth "github.com/talhacan/gateway/handlers/user"
	"github.com/talhacan/gateway/middleware"
	"github.com/talhacan/gateway/services"
)

// RegisterManagerRoutes manager API rotalarını kaydeder: login, refresh, logout, me + user CRUD.
func RegisterManagerRoutes(routerGroup *gin.RouterGroup, managerService *services.ManagerService, userService *services.UserService, tokenService *services.TokenService) {
	authHandler := managerauth.NewAuthHandler(managerService, tokenService)
	userHandler := userauth.NewUserHandler(userService)

	// Public
	manager := routerGroup.Group("/manager")
	{
		manager.POST("/auth/login", authHandler.Login)
		manager.POST("/auth/refresh", authHandler.RefreshToken)
	}

	// Protected (JWT + RequireManager)
	managerProtected := routerGroup.Group("/manager")
	managerProtected.Use(middleware.AuthMiddleware(tokenService), middleware.RequireManager())
	{
		managerProtected.GET("/auth/me", authHandler.Me)
		managerProtected.POST("/auth/logout", authHandler.Logout)
		// User yönetimi (manager sadece user CRUD yapabilir)
		managerProtected.GET("/user/getAll", userHandler.GetAll)
		managerProtected.GET("/user/:id", userHandler.GetByID)
		managerProtected.POST("/user", userHandler.Create)
		managerProtected.PUT("/user/:id", userHandler.Update)
		managerProtected.DELETE("/user/:id", userHandler.Delete)
	}
}
