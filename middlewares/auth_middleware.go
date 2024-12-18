package middlewares

import (
	"net/http"
	"strconv"
	"to-do-list-api/models"
	"to-do-list-api/pkg"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken, err := c.Cookie("session_token")
		if err != nil || sessionToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Non autorisé"})
			c.Abort()
			return
		}
		// Vérifier si l'utilisateur existe
		userID, err := strconv.Atoi(sessionToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session invalide"})
			c.Abort()
			return
		}

		var user models.User
		query := pkg.DB
		if err := query.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session invalide"})
			c.Abort()
			return
		}

		// Ajouter l'utilisateur au contexte pour une utilisation ultérieure
		c.Set("currentUser", user)

		// Continuer vers le prochain middleware ou handler
		c.Next()
	}
}
