package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/talhacan/gateway/config"
	"github.com/talhacan/gateway/database"
	"github.com/talhacan/gateway/docs"
	"github.com/talhacan/gateway/redis"
	"github.com/talhacan/gateway/routes"
	"github.com/talhacan/gateway/services"
)

func main() {
	_ = godotenv.Load()

	appConfig := config.Load()

	if databaseConnectErr := database.Connect(appConfig); databaseConnectErr != nil {
		log.Fatalf("database: %v", databaseConnectErr)
	}
	sqlDB, dbInstanceErr := database.DB.DB()
	if dbInstanceErr != nil {
		log.Fatalf("database.DB: %v", dbInstanceErr)
	}
	defer sqlDB.Close()

	if redisConnectErr := redis.Connect(appConfig); redisConnectErr != nil {
		log.Fatalf("redis: %v", redisConnectErr)
	}
	defer redis.Client.Close()

	// Servisler
	tokenService := services.NewTokenService(appConfig)
	userService := services.NewUserService(tokenService)
	managerService := services.NewManagerService(tokenService)

	ginEngine := gin.Default()
	// CORS: localhost (dev) + production frontend; ek origin'ler CORS_ORIGINS env ile (virgülle ayrılmış)
	corsOrigins := []string{
		"http://localhost:7203", "http://localhost:5173", "http://127.0.0.1:7203", "http://127.0.0.1:5173",
		"https://talhacan.com", "https://www.talhacan.com", "http://talhacan.com", "http://www.talhacan.com",
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
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	ginEngine.GET("/health", func(ginContext *gin.Context) {
		ginContext.JSON(200, gin.H{"status": "ok"})
	})
	// OpenAPI 3 spec (Scalar dokümantasyonu için); PUBLIC_API_URL varsa servers[0].url olarak kullanılır (canlıda base URL için)
	ginEngine.GET("/openapi.json", func(ginContext *gin.Context) {
		specBytes, specErr := docs.OpenAPI3()
		if specErr != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build spec"})
			return
		}
		if publicURL := strings.TrimSpace(os.Getenv("PUBLIC_API_URL")); publicURL != "" {
			var spec map[string]interface{}
			if jsonErr := json.Unmarshal(specBytes, &spec); jsonErr == nil {
				spec["servers"] = []map[string]interface{}{{"url": publicURL, "description": "API base"}}
				if out, err := json.Marshal(spec); err == nil {
					specBytes = out
				}
			}
		}
		ginContext.Header("Content-Type", "application/json; charset=utf-8")
		ginContext.Data(http.StatusOK, "application/json", specBytes)
	})
	// Scalar API Reference UI
	ginEngine.GET("/docs", func(ginContext *gin.Context) {
		ginContext.Header("Content-Type", "text/html; charset=utf-8")
		ginContext.String(http.StatusOK, docs.ScalarHTML())
	})

	apiGroup := ginEngine.Group("/api")
	routes.RegisterUserRoutes(apiGroup, userService, tokenService)
	routes.RegisterManagerRoutes(apiGroup, managerService, userService, tokenService)

	serverAddress := ":" + appConfig.Server.Port
	log.Printf("server listening on %s", serverAddress)
	if runErr := ginEngine.Run(serverAddress); runErr != nil {
		log.Fatalf("run: %v", runErr)
	}
}
