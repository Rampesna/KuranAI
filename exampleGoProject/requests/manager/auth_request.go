package manager

import "github.com/talhacan/gateway/dto"

// LoginRequest manager auth login için request (DeviceInfo core tarafından set edilir).
type LoginRequest struct {
	Credential string         `json:"credential" binding:"required"`
	Password   string         `json:"password" binding:"required"`
	DeviceInfo *dto.DeviceInfo // core tarafından set edilir
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
