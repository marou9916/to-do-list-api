package controllers

import (
	"net/http"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	//Vérifier que le user associé est existant
	var user models.User

	if err := query.First(&user, task.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Utilisateur associé introuvable"})
		return
	}

	// Vérifier que le statut est valide
	if task.Status == "" || task.Status == " " || !validStatUses[task.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Statut invalide"})
		return
	}

	// Enregistrer la tâche dans la base de données
	if err := query.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la tâche"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Tâche créée avec succès", "task": task})

}

// UpdateTask permet de mettre à jour une tâche
func UpdateTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	query := pkg.DB

	//Récupérer la tâche à mettre à jour
	if err := query.First(&task, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tâche avec l'ID " + id + " introuvable"})
		return
	}

	var updatedTask models.Task

	//Récupérer les données du corps
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	//Mettre à jour les champs de la tâche (vérifier Status si modifié)
	if updatedTask.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Statut invalide"})
		return
	} else if updatedTask.Status != "" {
		if !validStatUses[updatedTask.Status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Statut invalide"})
			return
		}
		task.Status = updatedTask.Status
	}

	task.Title = updatedTask.Title
	if err := query.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour de la tâche"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tâche mis à jour avec succès", "task": task})
}

// DeleteTask permet de supprimer une tâche
func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	query := pkg.DB

	if err := query.Delete(&models.Task{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression de la tâche"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tâche supprimée avec succès"})
}
