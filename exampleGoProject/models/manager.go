package models

import (
	"time"

	"gorm.io/gorm"
)

// Manager yönetici modeli. Sütun sırası: id, email, username, password, name, status, created_at, updated_at, deleted_at.
type Manager struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Username string `gorm:"uniqueIndex" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Name     string `gorm:"not null" json:"name"`
	Status   int    `gorm:"default:1" json:"status"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName GORM tablo adı.
func (Manager) TableName() string {
	return "managers"
}
