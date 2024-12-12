package models

import (
	"gorm.io/gorm"
)

// Task représente une tâche dans le système

	type Task struct {
		gorm.Model
		Title  string `gorm:"not null" json:"title"`
		Status string `gorm:"check:status IN ('to-do','in-progress','done')" json:"status"`
		UserID uint   `gorm:"not null" json:"user_id"` // Clé étrangère
		User   User   `gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID; references:ID" json:"-"`
	}
