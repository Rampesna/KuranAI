package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rampesna/kuran-ai-gateway/config"
	"github.com/rampesna/kuran-ai-gateway/core"
	"github.com/rampesna/kuran-ai-gateway/database"
	"github.com/rampesna/kuran-ai-gateway/dto"
	"github.com/rampesna/kuran-ai-gateway/models"
	"github.com/rampesna/kuran-ai-gateway/redis"
	goredis "github.com/redis/go-redis/v9"
)

// TokenService JWT token üretimi, doğrulama ve Redis/DB kaydı.
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
	DeviceType    string `json:"-"`
	jwt.RegisteredClaims
}

// IssueTokens access + refresh token oluşturur, DB ve Redis'e yazar.
func (s *TokenService) IssueTokens(ctx context.Context, tokenableID uint, tokenableType models.TokenableType, device *dto.DeviceInfo) *core.ServiceResponse {
	accessToken, accessExp, err := s.generateToken(tokenableID, tokenableType, 0, s.appConfig.JWT.AccessSecret, time.Duration(s.appConfig.JWT.AccessExpireMin)*time.Minute)
	if err != nil {
		return core.NewServiceResponse(false, "Failed to generate access token", 500, nil)
	}
	userAgent, ipAddr, deviceInfo, deviceType := "", "", "", ""
	if device != nil {
		userAgent, ipAddr, deviceInfo, deviceType = device.UserAgent, device.IPAddress, device.DeviceInfo, device.DeviceType
	}
	token := &models.JWTToken{
		TokenableID:   tokenableID,
		TokenableType: tokenableType,
		AccessToken:   accessToken,
		ExpiresAt:     accessExp,
		UserAgent:     userAgent,
		IPAddress:     ipAddr,
		DeviceInfo:    deviceInfo,
		DeviceType:    deviceType,
	}
	refreshToken, _, err := s.generateToken(tokenableID, tokenableType, 0, s.appConfig.JWT.RefreshSecret, time.Duration(s.appConfig.JWT.RefreshExpireDay)*24*time.Hour)
	if err != nil {
		return core.NewServiceResponse(false, "Failed to generate refresh token", 500, nil)
	}
	token.RefreshToken = refreshToken

	if createErr := database.DB.WithContext(ctx).Create(token).Error; createErr != nil {
		return core.NewServiceResponse(false, "Failed to save token", 500, nil)
	}

	accessKey := redis.KeyPrefixAccess + accessToken
	refreshKey := redis.KeyPrefixRefresh + refreshToken
	accessTTL := time.Until(accessExp)
	refreshTTL := time.Duration(s.appConfig.JWT.RefreshExpireDay) * 24 * time.Hour
	if accessTTL > 0 {
		redis.Client.Set(ctx, accessKey, fmt.Sprintf("%d|%s|%s", token.ID, string(tokenableType), normalizeDeviceType(deviceType)), accessTTL)
	}
	if refreshTTL > 0 {
		redis.Client.Set(ctx, refreshKey, fmt.Sprintf("%d|%d|%s", token.ID, tokenableID, tokenableType), refreshTTL)
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

func (s *TokenService) generateToken(tokenableID uint, tokenableType models.TokenableType, tokenID uint, secret string, ttl time.Duration) (string, time.Time, error) {
	exp := time.Now().Add(ttl)
	claims := &Claims{
		TokenableID:   tokenableID,
		TokenableType: string(tokenableType),
		TokenID:       tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.appConfig.JWT.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, exp, nil
}

// RevokeAccessToken access token'ı Redis'ten siler (logout).
func (s *TokenService) RevokeAccessToken(ctx context.Context, accessToken string) *core.ServiceResponse {
	redis.Client.Del(ctx, redis.KeyPrefixAccess+accessToken)
	return core.NewServiceResponse(true, "Logged out successfully", 200, nil)
}

// ValidateAccessToken token'ı Redis + JWT ile doğrular.
func (s *TokenService) ValidateAccessToken(ctx context.Context, accessToken string) (*Claims, *core.ServiceResponse) {
	redisValue, err := redis.Client.Get(ctx, redis.KeyPrefixAccess+accessToken).Result()
	if err != nil {
		return nil, core.NewServiceResponse(false, "Invalid or expired token", 401, nil)
	}
	claims, err := s.parseAccessToken(accessToken)
	if err != nil {
		return nil, core.NewServiceResponse(false, "Invalid token", 401, nil)
	}
	parts := strings.Split(redisValue, "|")
	if len(parts) > 0 {
		if tokenID, parseErr := strconv.ParseUint(parts[0], 10, 64); parseErr == nil {
			claims.TokenID = uint(tokenID)
		}
	}
	if len(parts) > 2 {
		claims.DeviceType = normalizeDeviceType(parts[2])
	}
	return claims, core.NewServiceResponse(true, "OK", 200, nil)
}

// TrackUserPresence records authenticated user activity without persisting
// request-level events in PostgreSQL. Members older than two minutes are pruned.
func (s *TokenService) TrackUserPresence(ctx context.Context, claims *Claims) {
	if claims == nil || claims.TokenableType != string(models.TokenableTypeUser) {
		return
	}
	platform := normalizeDeviceType(claims.DeviceType)
	member := fmt.Sprintf("%d|%s|%d", claims.TokenableID, platform, claims.TokenID)
	now := time.Now().Unix()
	pipe := redis.Client.TxPipeline()
	pipe.ZAdd(ctx, redis.KeyUserPresence, goredis.Z{Score: float64(now), Member: member})
	pipe.ZRemRangeByScore(ctx, redis.KeyUserPresence, "-inf", strconv.FormatInt(now-120, 10))
	pipe.Expire(ctx, redis.KeyUserPresence, 10*time.Minute)
	_, _ = pipe.Exec(ctx)
}

func normalizeDeviceType(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	switch {
	case strings.Contains(value, "ios"), strings.Contains(value, "iphone"), strings.Contains(value, "ipad"):
		return "ios"
	case strings.Contains(value, "android"):
		return "android"
	case strings.Contains(value, "web"), strings.Contains(value, "browser"):
		return "web"
	default:
		return "unknown"
	}
}

func (s *TokenService) parseAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.appConfig.JWT.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// RefreshTokenRequest Login/Refresh handler'dan gelen request interface.
type RefreshTokenRequest interface {
	GetRefreshToken() string
	GetDeviceInfo() *dto.DeviceInfo
}

// RefreshTokens refresh token ile yeni access + refresh üretir.
func (s *TokenService) RefreshTokens(ctx context.Context, req RefreshTokenRequest) *core.ServiceResponse {
	refreshToken := req.GetRefreshToken()
	refreshKey := redis.KeyPrefixRefresh + refreshToken
	if _, err := redis.Client.Get(ctx, refreshKey).Result(); err != nil {
		return core.NewServiceResponse(false, "Invalid or expired refresh token", 401, nil)
	}
	token, err := jwt.ParseWithClaims(refreshToken, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.appConfig.JWT.RefreshSecret), nil
	})
	if err != nil || !token.Valid {
		return core.NewServiceResponse(false, "Invalid refresh token", 401, nil)
	}
	claims := token.Claims.(*Claims)
	redis.Client.Del(ctx, refreshKey)
	return s.IssueTokens(ctx, claims.TokenableID, models.TokenableType(claims.TokenableType), req.GetDeviceInfo())
}
