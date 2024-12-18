package pkg

import (
	"log"
	"to-do-list-api/migrations"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	if DB, err = gorm.Open(sqlite.Open("./todo.db"), &gorm.Config{}); err != nil {
		log.Fatal("Échec de la connexion à la base de données :", err)
		return
	}
	log.Println("Base de données connectée avec succès !")

	DB.Exec("PRAGMA foreign_keys = ON;")

	if err = migrations.Migrate(DB); err != nil {
		log.Fatal("Échec de la migration des modèles :", err)
		return
	}
	log.Println("Migration des modèles effectuée avec succès !")
}
