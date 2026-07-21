package subscription

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rampesna/kuran-ai-gateway/core"
	"github.com/rampesna/kuran-ai-gateway/middleware"
	"github.com/rampesna/kuran-ai-gateway/services"
)

// SubscriptionHandler abonelik API handler.
type SubscriptionHandler struct {
	subSvc *services.SubscriptionService
	usageSvc *services.UsageService
}

// NewSubscriptionHandler SubscriptionHandler oluşturur.
func NewSubscriptionHandler(subSvc *services.SubscriptionService, usageSvc *services.UsageService) *SubscriptionHandler {
	return &SubscriptionHandler{subSvc: subSvc, usageSvc: usageSvc}
}

// Me GET /api/v1/subscription/me - kullanıcının abonelik durumu ve günlük kullanımı.
func (h *SubscriptionHandler) Me(c *gin.Context) {
	userIDVal, ok := c.Get(string(middleware.ContextKeyUserID))
	if !ok {
		c.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	userID := userIDVal.(uint)

	sub, pkg, _ := h.subSvc.GetActiveSubscription(c.Request.Context(), userID)
	msgUsed, msgLimit, ttsUsed, ttsLimit, pkgID, hasPackage, err := h.usageSvc.GetDailyUsage(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.NewServiceResponse(false, "usage fetch failed", 500, nil))
		return
	}

	data := gin.H{
		"subscription": gin.H{
			"hasPackage": hasPackage,
			"packageId":  pkgID,
			"package":    nil,
			"expiresAt":  nil,
			"isMonthly":  nil,
		},
		"usage": gin.H{
			"messageUsed":  msgUsed,
			"messageLimit": msgLimit,
			"ttsUsed":      ttsUsed,
			"ttsLimit":     ttsLimit,
		},
	}

	if sub != nil && pkg != nil {
		data["subscription"] = gin.H{
			"hasPackage": true,
			"packageId":  pkg.ID,
			"package": gin.H{
				"id":                 pkg.ID,
				"name":               pkg.Name,
				"nameEn":             pkg.NameEn,
				"slug":               pkg.Slug,
				"dailyQuestionLimit": pkg.DailyQuestionLimit,
				"dailyTTSLimit":      pkg.DailyTTSLimit,
				"monthlyPriceTL":     pkg.MonthlyPriceTL,
				"yearlyPriceTL":      pkg.YearlyPriceTL,
			},
			"expiresAt": sub.ExpiresAt,
			"isMonthly": sub.IsMonthly,
		}
	}

	c.JSON(http.StatusOK, core.NewServiceResponse(true, "ok", 200, data))
}

// Packages GET /api/v1/subscription/packages - tüm paketler.
func (h *SubscriptionHandler) Packages(c *gin.Context) {
	pkgs, err := h.subSvc.GetPackages(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.NewServiceResponse(false, "packages fetch failed", 500, nil))
		return
	}
	list := make([]gin.H, 0, len(pkgs))
	for _, p := range pkgs {
		list = append(list, gin.H{
			"id":                 p.ID,
			"name":               p.Name,
			"nameEn":             p.NameEn,
			"slug":               p.Slug,
			"dailyQuestionLimit": p.DailyQuestionLimit,
			"dailyTTSLimit":      p.DailyTTSLimit,
			"monthlyPriceTL":     p.MonthlyPriceTL,
			"yearlyPriceTL":      p.YearlyPriceTL,
			"sortOrder":          p.SortOrder,
		})
	}
	c.JSON(http.StatusOK, core.NewServiceResponse(true, "ok", 200, gin.H{"packages": list}))
}
