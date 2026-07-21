package dto

// DeviceInfo request'ten (header/client) alınan cihaz bilgisi.
type DeviceInfo struct {
	UserAgent  string
	IPAddress  string
	DeviceInfo string
	DeviceType string
}

// AuthResponse login/refresh cevabı.
type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    string `json:"expiresAt"`
	TokenType    string `json:"tokenType"`
}
