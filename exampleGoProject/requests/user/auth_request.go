package user

import "github.com/talhacan/gateway/dto"

// RegisterRequest kullanıcı kaydı için request (public endpoint).
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
}

// LoginRequest user auth login için request (DeviceInfo core.WithValidatedJSONAndDeviceInfo tarafından set edilir).
type LoginRequest struct {
	Credential string          `json:"credential" binding:"required"`
	Password   string          `json:"password" binding:"required"`
	DeviceInfo *dto.DeviceInfo  // core tarafından set edilir
}

// SetDeviceInfo core.DeviceInfoSetter implementasyonu.
func (r *LoginRequest) SetDeviceInfo(d *dto.DeviceInfo) { r.DeviceInfo = d }

// RefreshTokenRequest refresh token için request (TokenService.RefreshTokenRequest interface'i implement eder).
type RefreshTokenRequest struct {
	RefreshToken string         `json:"refreshToken" binding:"required"`
	DeviceInfo   *dto.DeviceInfo // core tarafından set edilir
}

// SetDeviceInfo core.DeviceInfoSetter implementasyonu.
func (r *RefreshTokenRequest) SetDeviceInfo(d *dto.DeviceInfo) { r.DeviceInfo = d }

// GetRefreshToken services.RefreshTokenRequest interface.
func (r *RefreshTokenRequest) GetRefreshToken() string { return r.RefreshToken }

// GetDeviceInfo services.RefreshTokenRequest interface.
func (r *RefreshTokenRequest) GetDeviceInfo() *dto.DeviceInfo { return r.DeviceInfo }
