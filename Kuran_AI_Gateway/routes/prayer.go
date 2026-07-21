package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rampesna/kuran-ai-gateway/handlers/prayer"
	"github.com/rampesna/kuran-ai-gateway/services"
)

// RegisterPrayerRoutes namaz vakitleri API rotalarını kaydeder (auth yok)
func RegisterPrayerRoutes(rg *gin.RouterGroup, prayerService *services.PrayerTimesService) {
	handler := prayer.NewPrayerHandler(prayerService)
	v1 := rg.Group("/v1")
	{
		v1.GET("/prayer-times", handler.GetPrayerTimes)
	}
}
