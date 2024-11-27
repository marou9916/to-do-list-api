package main

import (
	"log"
	"to-do-list-api/pkg"
)

func main() {
	// Initialiser la base de données
	pkg.InitDatabase()

	log.Println("API To-Do List prête à démarrer !")
}
