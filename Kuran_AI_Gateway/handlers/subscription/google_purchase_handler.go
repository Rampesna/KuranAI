package subscription

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rampesna/kuran-ai-gateway/core"
	"github.com/rampesna/kuran-ai-gateway/middleware"
	subscriptionRequest "github.com/rampesna/kuran-ai-gateway/requests/subscription"
	"github.com/rampesna/kuran-ai-gateway/services"
)

type GooglePurchaseHandler struct {
	subSvc     *services.SubscriptionService
	googlePlay services.GooglePlayPurchaseVerifier
}

func NewGooglePurchaseHandler(
	subSvc *services.SubscriptionService,
	googlePlay services.GooglePlayPurchaseVerifier,
) *GooglePurchaseHandler {
	return &GooglePurchaseHandler{subSvc: subSvc, googlePlay: googlePlay}
}

// Sync POST /api/v1/subscription/google/sync verifies a Play purchase on the
// server before granting any entitlement.
func (h *GooglePurchaseHandler) Sync(c *gin.Context) {
	userIDValue, ok := c.Get(string(middleware.ContextKeyUserID))
	if !ok {
		c.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	if h.googlePlay == nil || !h.googlePlay.Ready() {
		c.JSON(http.StatusServiceUnavailable, core.NewServiceResponse(false, "Google Play doğrulaması yapılandırılmadı", 503, nil))
		return
	}
	userID := userIDValue.(uint)

	var request subscriptionRequest.GooglePurchaseRequest
	if err := c.ShouldBindJSON(&request); err != nil ||
		strings.TrimSpace(request.ProductID) == "" ||
		strings.TrimSpace(request.PurchaseToken) == "" {
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "Google Play satın alma bilgisi gerekli", 400, nil))
		return
	}

	mapping, knownProduct := productIDMap[strings.TrimSpace(request.ProductID)]
	if !knownProduct {
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "Bilinmeyen ürün", 400, nil))
		return
	}
	verified, err := h.googlePlay.VerifySubscription(
		c.Request.Context(),
		request.PackageName,
		request.ProductID,
		request.PurchaseToken,
	)
	if err != nil {
		log.Printf("[GooglePurchase] verification failed: %v", err)
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "Google Play satın alma doğrulanamadı", 400, nil))
		return
	}

	packages, err := h.subSvc.GetPackages(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.NewServiceResponse(false, "Paketler alınamadı", 500, nil))
		return
	}
	var packageID uint
	for _, item := range packages {
		if item.Slug == mapping.Slug {
			packageID = item.ID
			break
		}
	}
	if packageID == 0 {
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "Paket bulunamadı", 400, nil))
		return
	}

	saved, err := h.subSvc.CreateGoogleSubscription(
		c.Request.Context(),
		userID,
		packageID,
		mapping.IsMonthly,
		verified,
	)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Abonelik oluşturulamadı"
		if errors.Is(err, services.ErrGoogleAccountConflict) {
			status = http.StatusConflict
			message = "Bu abonelik başka bir hesaba bağlı"
		} else if errors.Is(err, services.ErrInvalidGooglePurchase) {
			status = http.StatusBadRequest
			message = "Google Play satın alma doğrulanamadı"
		}
		c.JSON(status, core.NewServiceResponse(false, message, status, nil))
		return
	}

	c.JSON(http.StatusOK, core.NewServiceResponse(true, "Abonelik aktif", 200, gin.H{
		"subscription": gin.H{
			"id":           saved.ID,
			"packageId":    saved.PackageID,
			"expiresAt":    saved.ExpiresAt,
			"acknowledged": verified.Acknowledged,
		},
	}))
}
