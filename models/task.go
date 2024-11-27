package models

import (
	"gorm.io/gorm"
)

// Task représente une tâche dans le système
type Task struct {
	gorm.Model
	Title  string `gorm:"not null" json:"title"`
	Status string `gorm:"default:pending" json:"status"`
	UserID uint   `gorm:"not null" json:"user_id"` // Clé étrangère
	User   User   `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
}
