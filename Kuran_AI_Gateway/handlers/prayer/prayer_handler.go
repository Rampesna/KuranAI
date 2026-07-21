package prayer

import (
	"github.com/gin-gonic/gin"
	"github.com/rampesna/kuran-ai-gateway/core"
	prayerRequest "github.com/rampesna/kuran-ai-gateway/requests/prayer"
	"github.com/rampesna/kuran-ai-gateway/services"
)

// PrayerHandler namaz vakitleri handler
type PrayerHandler struct {
	prayerService *services.PrayerTimesService
}

// NewPrayerHandler PrayerHandler oluşturur
func NewPrayerHandler(prayerService *services.PrayerTimesService) *PrayerHandler {
	return &PrayerHandler{prayerService: prayerService}
}

// GetPrayerTimes GET /api/v1/prayer-times?lat=...&lng=... (opsiyonel; yoksa İstanbul)
func (h *PrayerHandler) GetPrayerTimes(c *gin.Context) {
	core.WithValidatedQuery(c, func(req *prayerRequest.PrayerTimesRequest) *core.ServiceResponse {
		return h.prayerService.GetPrayerTimes(c.Request.Context(), req)
	})
}
