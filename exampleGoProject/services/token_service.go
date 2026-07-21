package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/talhacan/gateway/config"
	"github.com/talhacan/gateway/core"
	"github.com/talhacan/gateway/database"
	"github.com/talhacan/gateway/dto"
	"github.com/talhacan/gateway/models"
	"github.com/talhacan/gateway/redis"
)

// TokenService JWT token üretimi, doğrulama ve Redis/DB kaydı (UserService/ManagerService tarafından kullanılır).
type TokenService struct {
	appConfig *config.Config
}

// NewTokenService TokenService oluşturur.
func NewTokenService(appConfig *config.Config) *TokenService {
	return &TokenService{appConfig: appConfig}
}

// Claims JWT claims.
type Claims struct {
	TokenableID   uint   `json:"tokenableId"`
	TokenableType string `json:"tokenableType"`
	TokenID       uint   `json:"tokenId"`
	jwt.RegisteredClaims
}

// IssueTokens access + refresh token oluşturur, DB ve Redis'e yazar. UserService/ManagerService Login içinden çağrılır.
func (tokenService *TokenService) IssueTokens(ctx context.Context, tokenableID uint, tokenableType models.TokenableType, device *dto.DeviceInfo) *core.ServiceResponse {
	accessToken, accessExp, generateErr := tokenService.generateToken(tokenableID, tokenableType, 0, tokenService.appConfig.JWT.AccessSecret, time.Duration(tokenService.appConfig.JWT.AccessExpireMin)*time.Minute)
	if generateErr != nil {
		return core.NewServiceResponse(
			false,
			"Failed to generate access token",
			500,
			nil,
		)
	}
	userAgent, ipAddr, deviceInfo, deviceType := "", "", "", ""
	if device != nil {
		userAgent, ipAddr, deviceInfo, deviceType = device.UserAgent, device.IPAddress, device.DeviceInfo, device.DeviceType
	}
	personalAccessToken := &models.PersonalAccessToken{
		TokenableID:   tokenableID,
		TokenableType: tokenableType,
		AccessToken:   accessToken,
		ExpiresAt:     accessExp,
		UserAgent:     userAgent,
		IPAddress:     ipAddr,
		DeviceInfo:    deviceInfo,
		DeviceType:    deviceType,
	}

	refreshToken, _, refreshGenerateErr := tokenService.generateToken(tokenableID, tokenableType, 0, tokenService.appConfig.JWT.RefreshSecret, time.Duration(tokenService.appConfig.JWT.RefreshExpireDay)*24*time.Hour)
	if refreshGenerateErr != nil {
		return core.NewServiceResponse(
			false,
			"Failed to generate refresh token",
			500,
			nil,
		)
	}
	personalAccessToken.RefreshToken = refreshToken

	if createErr := database.DB.WithContext(ctx).Create(personalAccessToken).Error; createErr != nil {
		return core.NewServiceResponse(
			false,
			"Failed to save token",
			500,
			nil,
		)
	}

	accessKey := redis.KeyPrefixAccess + accessToken
	refreshKey := redis.KeyPrefixRefresh + refreshToken
	accessTTL := time.Until(accessExp)
	refreshTTL := time.Duration(tokenService.appConfig.JWT.RefreshExpireDay) * 24 * time.Hour
	if accessTTL > 0 {
		redis.Client.Set(ctx, accessKey, fmt.Sprintf("%d|%s", personalAccessToken.ID, string(tokenableType)), accessTTL)
	}
	if refreshTTL > 0 {
		redis.Client.Set(ctx, refreshKey, fmt.Sprintf("%d|%d|%s", personalAccessToken.ID, tokenableID, tokenableType), refreshTTL)
	}

	return core.NewServiceResponse(
		true,
		"Login successful",
		201,
		&dto.AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    accessExp.Format(time.RFC3339),
			TokenType:    "Bearer",
		},
	)
}

func (tokenService *TokenService) generateToken(tokenableID uint, tokenableType models.TokenableType, tokenID uint, secret string, ttl time.Duration) (string, time.Time, error) {
	exp := time.Now().Add(ttl)
	claims := &Claims{
		TokenableID:   tokenableID,
		TokenableType: string(tokenableType),
		TokenID:       tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    tokenService.appConfig.JWT.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, signErr := token.SignedString([]byte(secret))
	if signErr != nil {
		return "", time.Time{}, signErr
	}
	return signed, exp, nil
}

// RevokeAccessToken access token'ı Redis'ten siler (logout); bir daha kullanılamaz.
func (tokenService *TokenService) RevokeAccessToken(ctx context.Context, accessToken string) *core.ServiceResponse {
	accessTokenRedisKey := redis.KeyPrefixAccess + accessToken
	redis.Client.Del(ctx, accessTokenRedisKey)
	return core.NewServiceResponse(
		true,
		"Logged out successfully",
		200,
		nil,
	)
}

// ValidateAccessToken token'ı Redis + JWT ile doğrular; geçerliyse claims döner (middleware kullanır).
func (tokenService *TokenService) ValidateAccessToken(ctx context.Context, accessToken string) (*Claims, *core.ServiceResponse) {
	accessTokenRedisKey := redis.KeyPrefixAccess + accessToken
	if _, redisGetErr := redis.Client.Get(ctx, accessTokenRedisKey).Result(); redisGetErr != nil {
		return nil, core.NewServiceResponse(
			false,
			"Invalid or expired token",
			401,
			nil,
		)
	}
	claims, parseErr := tokenService.parseAccessToken(accessToken)
	if parseErr != nil {
		return nil, core.NewServiceResponse(
			false,
			"Invalid token",
			401,
			nil,
		)
	}
	return claims, core.NewServiceResponse(
		true,
		"OK",
		200,
		nil,
	)
}

func (tokenService *TokenService) parseAccessToken(tokenString string) (*Claims, error) {
	jwtToken, parseErr := jwt.ParseWithClaims(tokenString, &Claims{}, func(jwtToken *jwt.Token) (interface{}, error) {
		return []byte(tokenService.appConfig.JWT.AccessSecret), nil
	})
	if parseErr != nil {
		return nil, parseErr
	}
	claims, ok := jwtToken.Claims.(*Claims)
	if !ok || !jwtToken.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// RefreshTokenRequest Login/Refresh handler'dan gelen request; user veya manager paketinden implement edilir.
type RefreshTokenRequest interface {
	GetRefreshToken() string
	GetDeviceInfo() *dto.DeviceInfo
}

// RefreshTokens refresh token ile yeni access + refresh üretir (handler request ile çağrılır).
func (tokenService *TokenService) RefreshTokens(ctx context.Context, req RefreshTokenRequest) *core.ServiceResponse {
	refreshToken := req.GetRefreshToken()
	refreshTokenRedisKey := redis.KeyPrefixRefresh + refreshToken
	if _, redisGetErr := redis.Client.Get(ctx, refreshTokenRedisKey).Result(); redisGetErr != nil {
		return core.NewServiceResponse(
			false,
			"Invalid or expired refresh token",
			401,
			nil,
		)
	}
	jwtToken, parseErr := jwt.ParseWithClaims(refreshToken, &Claims{}, func(jwtToken *jwt.Token) (interface{}, error) {
		return []byte(tokenService.appConfig.JWT.RefreshSecret), nil
	})
	if parseErr != nil || !jwtToken.Valid {
		return core.NewServiceResponse(
			false,
			"Invalid refresh token",
			401,
			nil,
		)
	}
	claims := jwtToken.Claims.(*Claims)
	redis.Client.Del(ctx, refreshTokenRedisKey)
	return tokenService.IssueTokens(ctx, claims.TokenableID, models.TokenableType(claims.TokenableType), req.GetDeviceInfo())
}
