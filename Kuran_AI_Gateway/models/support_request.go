package models

import (
	"time"

	"gorm.io/gorm"
)

// SupportRequest web sitesi ve uygulama üzerinden gelen destek talepleri.
type SupportRequest struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:128;not null" json:"name"`
	Email     string         `gorm:"size:256;not null" json:"email"`
	Subject   string         `gorm:"size:64;not null;default:'general'" json:"subject"`
	Message   string         `gorm:"type:text;not null" json:"message"`
	IPAddress string         `gorm:"size:64" json:"ipAddress,omitempty"`
	Status    string         `gorm:"size:32;default:'pending'" json:"status"` // pending, read, resolved
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (SupportRequest) TableName() string { return "support_requests" }
