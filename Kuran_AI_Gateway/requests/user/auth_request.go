package user

import "github.com/rampesna/kuran-ai-gateway/dto"

// RegisterRequest kullanıcı kaydı için request.
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
}

// LoginRequest user auth login için request (email/username + password).
type LoginRequest struct {
	Credential string `json:"credential" binding:"required"`
	Password   string `json:"password" binding:"required"`
	DeviceInfo *dto.DeviceInfo
}

func (r *LoginRequest) SetDeviceInfo(d *dto.DeviceInfo) { r.DeviceInfo = d }

// RefreshTokenRequest refresh token için request.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
	DeviceInfo   *dto.DeviceInfo
}

func (r *RefreshTokenRequest) SetDeviceInfo(d *dto.DeviceInfo) { r.DeviceInfo = d }
func (r *RefreshTokenRequest) GetRefreshToken() string         { return r.RefreshToken }
func (r *RefreshTokenRequest) GetDeviceInfo() *dto.DeviceInfo  { return r.DeviceInfo }

// GoogleLoginRequest Google OAuth ile giriş.
type GoogleLoginRequest struct {
	IDToken    string `json:"idToken" binding:"required"`
	DeviceInfo *dto.DeviceInfo
}

func (r *GoogleLoginRequest) SetDeviceInfo(d *dto.DeviceInfo) { r.DeviceInfo = d }

// AppleLoginRequest Apple Sign in ile giriş.
type AppleLoginRequest struct {
	IdentityToken string `json:"identityToken" binding:"required"`
	DeviceInfo    *dto.DeviceInfo
}

func (r *AppleLoginRequest) SetDeviceInfo(d *dto.DeviceInfo) { r.DeviceInfo = d }

// AppleMobileStartRequest web tabanlı Apple girişini Android'e güvenli şekilde devreder.
type AppleMobileStartRequest struct {
	IdentityToken string `json:"identityToken" binding:"required"`
	DeviceInfo    *dto.DeviceInfo
}

func (r *AppleMobileStartRequest) SetDeviceInfo(d *dto.DeviceInfo) { r.DeviceInfo = d }

// MobileAuthExchangeRequest tek kullanımlık mobil giriş kodunu oturum anahtarlarına çevirir.
type MobileAuthExchangeRequest struct {
	Code       string `json:"code" binding:"required,min=32,max=128"`
	DeviceInfo *dto.DeviceInfo
}

func (r *MobileAuthExchangeRequest) SetDeviceInfo(d *dto.DeviceInfo) { r.DeviceInfo = d }
