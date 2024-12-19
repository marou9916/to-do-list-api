package controllers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"to-do-list-api/models"
	"to-do-list-api/pkg"

	"github.com/gin-gonic/gin"
)

// Liste des statuts valides
var validStatUses = map[string]bool{
	"to-do":       true,
	"in-progress": true,
	"done":        true,
}

// GetTasks permet de récupérer la liste des tâches
func GetTasks(c *gin.Context) {
	var tasks []models.Task
	status := c.Query("status") //paramètre de filtrage

	query := pkg.DB
	if status != "" {
		if !validStatUses[status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Statut invalide. Options : 'to-do', 'in-progress', 'done'"})
			return
		}
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des tâches"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})

}

// CreateTask permet de créer une tâche
func CreateTask(c *gin.Context) {
	var task models.Task
	query := pkg.DB

	// Lier les données de la requête au modèle Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format des données invalides"})
		return
	}

	//Vérifier que le user associé est existant
	var user models.User

	if err := query.First(&user, task.UserID).Error; err != nil {
		log.Printf("Id user associé %d\n", task.UserID)
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Utilisateur associé introuvable"})
		return
	}

	//Vérifier que le titre est saisi
	if task.Title == "" || task.Title == " " {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le titre est requis"})
		return
	}

	// Vérifier que le statut est valide
	if task.Status == "" || task.Status == " " || !validStatUses[task.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Statut invalide. Options : 'to-do', 'in-progress', 'done'"})
		return
	}

	// Enregistrer la tâche dans la base de données
	if err := query.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la tâche"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("Tâche %s créée et associée au user %s avec succès", task.Title, user.Username)})

}

// UpdateTask permet de mettre à jour une tâche
func UpdateTask(c *gin.Context) {
	titleRegex := regexp.MustCompile(`^[\p{L}0-9\s]{4,}$`) // p{L} : pour autoriser tout caractère alphabétique (accentué ou non) dans le titre
	query := pkg.DB

	// Récupérer la tâche à partir du middleware
	taskFromContext, exists := c.Get("task")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de récupérer la tâche depuis le contexte"})
		return
	}

	task, ok := taskFromContext.(*models.Task)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne lors de la récupération de la tâche"})
		return
	}

	var updatedTask models.Task

	//Récupérer les données du corps
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	//Empêcher la modification de l'id du user associé
	if updatedTask.UserID != task.UserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pour associer à un autre user, veuillez créer une nouvelle tâche"})
		return
	}

	//Mettre à jour les champs de la tâche (vérifier Status si modifié)
	if updatedTask.Status != task.Status {
		if !validStatUses[updatedTask.Status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Statut invalide. Options : 'to-do', 'in-progress', 'done'"})
			return
		}
		task.Status = updatedTask.Status
	}
	if updatedTask.Title != task.Title {
		if !titleRegex.MatchString(updatedTask.Title) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Le titre doit comporter au moins 4 caractères alphanumériques."})
			return
		}

		updatedTask.Title = strings.TrimSpace(updatedTask.Title) //Nettoyer les espaces en excès avant de mettre à jour
		task.Title = updatedTask.Title
	}

	if err := query.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour de la tâche"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tâche mis à jour avec succès", "task": task})
}

// DeleteTask permet de supprimer une tâche
func DeleteTask(c *gin.Context) {
	// Récupérer la tâche depuis le contexte
	task, _ := c.Get("task")
	castedTask := task.(*models.Task)

	//Supprimer la tâche
	if err := pkg.DB.Unscoped().Delete(castedTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression de la tâche"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tâche supprimée avec succès"})
}
