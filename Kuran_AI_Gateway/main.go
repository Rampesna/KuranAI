package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rampesna/kuran-ai-gateway/config"
	"github.com/rampesna/kuran-ai-gateway/database"
	"github.com/rampesna/kuran-ai-gateway/docs"
	publicHandlers "github.com/rampesna/kuran-ai-gateway/handlers/public"
	"github.com/rampesna/kuran-ai-gateway/middleware"
	"github.com/rampesna/kuran-ai-gateway/redis"
	"github.com/rampesna/kuran-ai-gateway/routes"
	"github.com/rampesna/kuran-ai-gateway/services"
)

//go:embed static/privacy-policy.html
var privacyPolicyHTML []byte

//go:embed static/support.html
var supportHTML []byte

//go:embed static/terms.html
var termsHTML []byte

//go:embed static/contact.html
var contactHTML []byte

func main() {
	_ = godotenv.Load()

	// CLI: migrate subcommand (sadece migration çalıştırır ve çıkar)
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		appConfig := config.Load()
		if err := database.Connect(appConfig); err != nil {
			log.Fatalf("migration: %v", err)
		}
		log.Println("migration completed")
		return
	}

	appConfig := config.Load()
	if len(appConfig.JWT.AccessSecret) < 32 || len(appConfig.JWT.RefreshSecret) < 32 {
		log.Fatal("JWT_ACCESS_SECRET and JWT_REFRESH_SECRET must each contain at least 32 characters")
	}

	if err := database.Connect(appConfig); err != nil {
		log.Fatalf("database: %v", err)
	}
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatalf("database.DB: %v", err)
	}
	defer sqlDB.Close()

	if err := redis.Connect(appConfig); err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer redis.Client.Close()

	tokenService := services.NewTokenService(appConfig)
	userService := services.NewUserService(tokenService)
	passkeyService, err := services.NewPasskeyService(appConfig, tokenService)
	if err != nil {
		log.Fatalf("passkey: %v", err)
	}

	ginEngine := gin.Default()
	corsOrigins := []string{
		"http://localhost:7203", "http://localhost:5173", "http://127.0.0.1:7203", "http://127.0.0.1:5173",
	}
	if corsEnv := os.Getenv("CORS_ORIGINS"); corsEnv != "" {
		for _, o := range strings.Split(corsEnv, ",") {
			if o = strings.TrimSpace(o); o != "" {
				corsOrigins = append(corsOrigins, o)
			}
		}
	}
	ginEngine.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Passkey-Transaction", "X-Device-Info", "X-Device-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	ginEngine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	ginEngine.GET("/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
			return
		}
		if err := redis.Client.Ping(ctx).Err(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	ginEngine.GET("/privacy-policy", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", privacyPolicyHTML)
	})
	ginEngine.GET("/support", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", supportHTML)
	})
	ginEngine.GET("/terms", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", termsHTML)
	})
	ginEngine.GET("/contact", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", contactHTML)
	})

	// OpenAPI spec (Scalar için)
	ginEngine.GET("/openapi.json", func(c *gin.Context) {
		specBytes, err := docs.OpenAPI3()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build spec"})
			return
		}
		if publicURL := strings.TrimSpace(os.Getenv("PUBLIC_API_URL")); publicURL != "" {
			var spec map[string]interface{}
			if json.Unmarshal(specBytes, &spec) == nil {
				spec["servers"] = []map[string]interface{}{{"url": publicURL, "description": "API base"}}
				if out, err := json.Marshal(spec); err == nil {
					specBytes = out
				}
			}
		}
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.Data(http.StatusOK, "application/json", specBytes)
	})

	// Scalar API Reference UI
	ginEngine.GET("/docs", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, docs.ScalarHTML())
	})

	// Avatar uploads (static files)
	ginEngine.Static("/uploads", "./uploads")

	// Public support form API (rate limited: 5 requests per 10 minutes per IP)
	publicAPI := ginEngine.Group("/api/public")
	publicAPI.Use(middleware.RateLimiter(5, 10*time.Minute))
	{
		publicAPI.POST("/support", publicHandlers.SubmitSupportRequest)
	}

	subSvc := services.NewSubscriptionService()
	usageSvc := services.NewUsageService(subSvc)
	if err := subSvc.SeedPackages(context.Background()); err != nil {
		log.Printf("warning: subscription seed: %v", err)
	}

	apiGroup := ginEngine.Group("/api")
	prayerService := services.NewPrayerTimesService()
	verseSvc := services.NewDailyVerseService()
	managerService := services.NewManagerService(tokenService)
	routes.RegisterUserRoutes(apiGroup, userService, tokenService, passkeyService)
	routes.RegisterManagerRoutes(apiGroup, managerService, tokenService)
	routes.RegisterPrayerRoutes(apiGroup, prayerService)
	routes.RegisterChatRoutes(apiGroup, tokenService, appConfig.OpenAI.APIKey, usageSvc)
	routes.RegisterSubscriptionRoutes(apiGroup, tokenService, subSvc, usageSvc, appConfig.AppleIAP.SharedSecret, appConfig.AppleIAP.BundleID)
	routes.RegisterVerseRoutes(apiGroup, verseSvc, tokenService)

	addr := ":" + appConfig.Server.Port
	log.Printf("server listening on %s", addr)
	if err := ginEngine.Run(addr); err != nil {
		log.Fatalf("run: %v", err)
	}
}
