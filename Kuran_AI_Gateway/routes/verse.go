package routes

import (
	"github.com/gin-gonic/gin"
	verseHandler "github.com/rampesna/kuran-ai-gateway/handlers/verse"
	"github.com/rampesna/kuran-ai-gateway/middleware"
	"github.com/rampesna/kuran-ai-gateway/services"
)

// RegisterVerseRoutes günün ayeti API rotalarını kaydeder.
func RegisterVerseRoutes(rg *gin.RouterGroup, verseService *services.DailyVerseService, tokenService *services.TokenService) {
	handler := verseHandler.NewVerseHandler(verseService)

	// Public: Bugünün ayetini çek
	rg.GET("/verse-of-day", handler.GetToday)

	// Manager: CRUD
	managerGroup := rg.Group("/manager/verse-of-day")
	managerGroup.Use(middleware.AuthMiddleware(tokenService), middleware.RequireManager())
	{
		managerGroup.GET("", handler.List)
		managerGroup.GET("/:id", handler.GetByID)
		managerGroup.POST("", handler.Create)
		managerGroup.PUT("/:id", handler.Update)
		managerGroup.DELETE("/:id", handler.Delete)
	}
}
