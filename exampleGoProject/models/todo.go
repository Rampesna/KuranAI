package models

import (
	"time"

	"gorm.io/gorm"
)

// Todo todo modeli. Sütun sırası: id, user_id, manager_id, title, body, done, created_at, updated_at, deleted_at.
type Todo struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	UserID    *uint  `gorm:"index" json:"userId,omitempty"`
	ManagerID *uint  `gorm:"index" json:"managerId,omitempty"`
	Title     string `gorm:"not null" json:"title"`
	Body      string `json:"body"`
	Done      bool   `gorm:"default:false" json:"done"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName GORM tablo adı.
func (Todo) TableName() string {
	return "todos"
}
