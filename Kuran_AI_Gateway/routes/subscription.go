package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rampesna/kuran-ai-gateway/handlers/subscription"
	"github.com/rampesna/kuran-ai-gateway/middleware"
	"github.com/rampesna/kuran-ai-gateway/services"
)

// RegisterSubscriptionRoutes abonelik API rotalarını kaydeder.
func RegisterSubscriptionRoutes(
	rg *gin.RouterGroup,
	tokenService *services.TokenService,
	subSvc *services.SubscriptionService,
	usageSvc *services.UsageService,
	appleSecret,
	appleBundleID string,
	googlePlay services.GooglePlayPurchaseVerifier,
) {
	appleSvc := services.NewAppleReceiptService(appleSecret, appleBundleID)
	subHandler := subscription.NewSubscriptionHandler(subSvc, usageSvc)
	purchaseHandler := subscription.NewPurchaseHandler(subSvc, appleSvc)
	googlePurchaseHandler := subscription.NewGooglePurchaseHandler(subSvc, googlePlay)

	subProtected := rg.Group("/v1/subscription")
	subProtected.Use(middleware.AuthMiddleware(tokenService), middleware.RequireUser())
	{
		subProtected.GET("/me", subHandler.Me)
		subProtected.GET("/packages", subHandler.Packages)
		subProtected.POST("/purchase", purchaseHandler.Purchase)
		subProtected.POST("/google/sync", googlePurchaseHandler.Sync)
	}
}
