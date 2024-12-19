package middlewares

import (
	"net/http"
	"strconv"
	"to-do-list-api/models"
	"to-do-list-api/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthorizeTaskOwnerShip génère un middleware vérifiant les permissions d'exécution d'une opération sur une tâche
func AuthorizeTaskOwnerShip() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Récupérer l'utilisateur authentifié
		authentifiedUser, exists := c.Get("currentUser")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
			return
		}

		user, ok := authentifiedUser.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de récupérer l'utilisateur"})
			return
		}

		// Récupérer l'ID de la tâche depuis les paramètres
		taskID := c.Param("id")
		if _, err := strconv.Atoi(taskID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "L'ID doit être un entier valide"})
			return
		}
		// Vérifier que la tâche existe et appartient à l'utilisateur
		var task models.Task
		if err := pkg.DB.First(&task, taskID).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne lors de la vérification de l'appartenance de la tâche à l'utilisateur"})
				c.Abort()
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "Tâche non trouvée"})
				c.Abort()
			}
			return
		}

		if task.UserID != user.ID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Action non autorisée"})
			c.Abort()
			return
		}

		// Ajouter la tâche au contexte pour un usage ultérieur dans le contrôleur
		c.Set("task", &task) //le pointeur ici évite de dupliquer l'objet en mémoire

		// Continuer la chaîne des middlewares
		c.Next()
	}
}
