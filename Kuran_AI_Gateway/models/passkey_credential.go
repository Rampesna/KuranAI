package models

import (
	"time"

	"gorm.io/datatypes"
)

// PasskeyCredential stores a WebAuthn public-key credential. Private keys
// remain inside the user's authenticator and are never sent to the server.
type PasskeyCredential struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	UserID         uint           `gorm:"not null;index" json:"userId"`
	CredentialID   []byte         `gorm:"type:bytea;not null;uniqueIndex" json:"-"`
	CredentialData datatypes.JSON `gorm:"type:jsonb;not null" json:"-"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	LastUsedAt     *time.Time     `json:"lastUsedAt,omitempty"`
}

func (PasskeyCredential) TableName() string {
	return "passkey_credentials"
}
