package main

import (
	"log"
	"to-do-list-api/pkg"
	"to-do-list-api/routes"
)

func main() {
	// Initialiser la base de données
	pkg.InitDatabase()

	//Configurer le routeur
	router := routes.SetupRouter()

	log.Println("API To-Do List prête à démarrer !")
	log.Println("Le serveur fonctionne sur http://localhost:8080")

	//Lancer l'application
	router.Run(":8080")
}
