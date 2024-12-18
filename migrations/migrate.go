package migrations

import (
	"to-do-list-api/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Task{},
		&models.Session{},
	)
}
