package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/talhacan/gateway/core"
	"github.com/talhacan/gateway/services"
)

// ContextKey context key type.
type ContextKey string

const (
	ContextKeyClaims   ContextKey = "claims"
	ContextKeyUserID   ContextKey = "userId"
	ContextKeyManagerID ContextKey = "managerId"
)

// AuthMiddleware JWT token doğrular; geçerliyse claims'i context'e ekler.
func AuthMiddleware(tokenService *services.TokenService) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		authHeader := ginContext.GetHeader("Authorization")
		if authHeader == "" {
			ginContext.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Authorization header required", 401, nil))
			ginContext.Abort()
			return
		}

		headerParts := strings.SplitN(authHeader, " ", 2)
		if len(headerParts) != 2 || !strings.EqualFold(headerParts[0], "Bearer") {
			ginContext.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Invalid authorization format", 401, nil))
			ginContext.Abort()
			return
		}

		accessToken := strings.TrimSpace(headerParts[1])
		if accessToken == "" {
			ginContext.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Token required", 401, nil))
			ginContext.Abort()
			return
		}

		tokenClaims, serviceResponse := tokenService.ValidateAccessToken(ginContext.Request.Context(), accessToken)
		if !serviceResponse.IsSuccess {
			ginContext.JSON(serviceResponse.StatusCode, serviceResponse)
			ginContext.Abort()
			return
		}

		ginContext.Set(string(ContextKeyClaims), tokenClaims)
		if tokenClaims.TokenableType == "user" {
			ginContext.Set(string(ContextKeyUserID), tokenClaims.TokenableID)
		} else if tokenClaims.TokenableType == "manager" {
			ginContext.Set(string(ContextKeyManagerID), tokenClaims.TokenableID)
		}
		ginContext.Next()
	}
}

// RequireUser sadece user token'ı kabul eder.
func RequireUser() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		claimsValue, exists := ginContext.Get(string(ContextKeyClaims))
		if !exists {
			ginContext.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Unauthorized", 401, nil))
			ginContext.Abort()
			return
		}
		tokenClaims := claimsValue.(*services.Claims)
		if tokenClaims.TokenableType != "user" {
			ginContext.JSON(http.StatusForbidden, core.NewServiceResponse(false, "User access required", 403, nil))
			ginContext.Abort()
			return
		}
		ginContext.Next()
	}
}

// RequireManager sadece manager token'ı kabul eder.
func RequireManager() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		claimsValue, exists := ginContext.Get(string(ContextKeyClaims))
		if !exists {
			ginContext.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Unauthorized", 401, nil))
			ginContext.Abort()
			return
		}
		tokenClaims := claimsValue.(*services.Claims)
		if tokenClaims.TokenableType != "manager" {
			ginContext.JSON(http.StatusForbidden, core.NewServiceResponse(false, "Manager access required", 403, nil))
			ginContext.Abort()
			return
		}
		ginContext.Next()
	}
}
