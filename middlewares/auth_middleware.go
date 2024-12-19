package middlewares

import (
	"net/http"
	"to-do-list-api/models"
	"to-do-list-api/pkg"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := pkg.DB
		//Lire le cookie de session
		sessionToken, err := c.Cookie("session_token")
		if err != nil || sessionToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentification requise"})
			c.Abort()
			return
		}

		//Vérifier si le token correspond à une session valide
		var session models.Session
		if err := query.Where("token = ?", sessionToken).First(&session).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session invalide"})
			c.Abort()
			return
		}

		//Vérifier si la session a expiré
		if session.ExpiresAt.Before(pkg.TimeNow()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expirée"})
			c.Abort()
			return
		}

		//Récupérer l'utilisateur associé à la session
		var user models.User
		if err := query.First(&user, session.UserID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Utilisateur introuvable"})
			c.Abort()
			return
		}
		//Ajouter au contexte pour une utilisation ultérieure
		c.Set("currentUser", &user)

		// Continuer vers le prochain middleware ou handler
		c.Next()
	}
}
