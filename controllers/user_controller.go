package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"to-do-list-api/models"
	"to-do-list-api/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// CreateUser permet de créer un nouvel utilisateur
func CreateUser(c *gin.Context) {
	var user models.User

	query := pkg.DB

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	//Vérification du format de l'email
	if !emailRegex.MatchString(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format d'email invalide"})
		return
	}

	//Vérification de l'unicité de l'email
	var existingUser models.User
	if err := query.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification de l'email"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cet email est déjà utilisé"})
		return
	}

	if err := query.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("User %s créé avec succès", user.Username)})
}

// GetUser permet de récupérer un utilisateur par son ID
func GetUser(c *gin.Context) {
	id := c.Param("id")
	query := pkg.DB
	var user models.User

	if err := query.First(&user, id).Error; err != nil {
		if err := query.First(&user, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur introuvable"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération de l'utilisateur"})
			}
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{"user": user})

}

// GetUsers permet de récupérer tous les utilisateurs
func GetUsers(c *gin.Context) {
	var users []models.User
	query := pkg.DB

	if err := query.Find(&users).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur introuvable"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération de l'utilisateur"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// UpdateUser permet de mettre à jour les informations d'un utilisateur
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	query := pkg.DB

	var user models.User
	//Vérification de l'existence de l'utilisateur à mettre à jour dans la base de données
	if err := query.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User non trouvé"})
		return
	}

	var updatedUserData models.User

	if err := c.ShouldBindJSON(&updatedUserData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	//Vérification du format de l'email
	if !emailRegex.MatchString(updatedUserData.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format d'email invalide"})
		return
	}

	//Vérification de l'unicité de l'email
	if updatedUserData.Email != "" && updatedUserData.Email != user.Email {
		var existingUser models.User
		if err := query.Where("email = ?", updatedUserData.Email).First(&existingUser).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification de l'email"})
			}
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cet email est déjà utilisé"})
			return
		}

	}

	//Mettre à jour les données de l'utilisateur
	user.Username = updatedUserData.Username
	if updatedUserData.Email != "" {
		user.Email = updatedUserData.Email
	}

	//Sauvegarder dans la base de données
	if err := query.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour de l'utilisateur"})
		return
	}

	//Envoyer une réponse au client
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s créé avec succès", user.Username)})
}

// DeleteUser permet de supprimer un utilisateur
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// Supprimer l'utilisateur correspondant
	if err := pkg.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression de l'utilisateur"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Utilisateur supprimé avec succès"})
}
