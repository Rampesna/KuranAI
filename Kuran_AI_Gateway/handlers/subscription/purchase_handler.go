package subscription

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rampesna/kuran-ai-gateway/core"
	"github.com/rampesna/kuran-ai-gateway/middleware"
	"github.com/rampesna/kuran-ai-gateway/requests/subscription"
	"github.com/rampesna/kuran-ai-gateway/services"
)

// productID -> (packageSlug, isMonthly)
var productIDMap = map[string]struct {
	Slug      string
	IsMonthly bool
}{
	"kuranai.bronze.monthly":  {"bronze", true},
	"kuranai.bronze.yearly":   {"bronze", false},
	"kuranai.silver.monthly":  {"silver", true},
	"kuranai.silver.yearly":   {"silver", false},
	"kuranai.gold.monthly":    {"gold", true},
	"kuranai.gold.yearly":     {"gold", false},
	"kuranai.diamond.monthly": {"diamond", true},
	"kuranai.diamond.yearly":  {"diamond", false},
}

// PurchaseHandler satın alma işlemi.
type PurchaseHandler struct {
	subSvc   *services.SubscriptionService
	appleSvc *services.AppleReceiptService
}

// NewPurchaseHandler PurchaseHandler oluşturur.
func NewPurchaseHandler(subSvc *services.SubscriptionService, appleSvc *services.AppleReceiptService) *PurchaseHandler {
	return &PurchaseHandler{subSvc: subSvc, appleSvc: appleSvc}
}

// Purchase POST /api/v1/subscription/purchase - Apple receipt ile abonelik oluşturur.
func (h *PurchaseHandler) Purchase(c *gin.Context) {
	userIDVal, ok := c.Get(string(middleware.ContextKeyUserID))
	if !ok {
		c.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	userID := userIDVal.(uint)

	var req subscription.PurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.ReceiptData == "" {
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "receiptData required", 400, nil))
		return
	}

	receiptData := strings.TrimSpace(req.ReceiptData)

	result, err := h.appleSvc.Verify(receiptData)
	if err != nil {
		log.Printf("[Purchase] Apple verify failed: %v", err)
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "Geçersiz ödeme bilgisi: "+err.Error(), 400, nil))
		return
	}

	mapping, ok := productIDMap[result.ProductID]
	if !ok {
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "Bilinmeyen ürün: "+result.ProductID, 400, nil))
		return
	}
	if result.ExpiresAtUnix <= time.Now().UTC().Unix() {
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "Abonelik süresi dolmuş", 400, nil))
		return
	}

	pkgs, err := h.subSvc.GetPackages(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, core.NewServiceResponse(false, "Paketler alınamadı", 500, nil))
		return
	}
	var pkgID uint
	for _, p := range pkgs {
		if p.Slug == mapping.Slug {
			pkgID = p.ID
			break
		}
	}
	if pkgID == 0 {
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "Paket bulunamadı", 400, nil))
		return
	}

	var expiresAt *time.Time
	if result.ExpiresAtUnix > 0 {
		t := time.Unix(result.ExpiresAtUnix, 0).UTC()
		expiresAt = &t
	}
	sub, err := h.subSvc.CreateSubscription(c.Request.Context(), userID, pkgID, mapping.IsMonthly, result.OriginalTxID, result.ProductID, expiresAt)
	if err != nil {
		log.Printf("[Purchase] CreateSubscription error: %v", err)
		c.JSON(http.StatusInternalServerError, core.NewServiceResponse(false, "Abonelik oluşturulamadı", 500, nil))
		return
	}

	c.JSON(http.StatusOK, core.NewServiceResponse(true, "Abonelik aktif", 200, gin.H{
		"subscription": gin.H{
			"id":        sub.ID,
			"packageId": sub.PackageID,
			"expiresAt": sub.ExpiresAt,
		},
	}))
}
