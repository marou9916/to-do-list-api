package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"to-do-list-api/models"
	"to-do-list-api/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateUser permet de créer un nouvel utilisateur
func CreateUser(c *gin.Context) {
	var user models.User

	query := pkg.DB

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	//Vérifier que username et email sont non nuls
	if user.Username == "" || user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le username et l'email sont requis"})
		return
	}

	//Vérification de l'unicité du username
	var existingUser models.User

	if err := query.Where("username = ?", user.Username).First(&existingUser).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne lors de la vérification du username"})
			return
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Ce username est déjà utilisé"})
		return
	}

	//Vérification du format du user
	if !pkg.ValidateUsernameFormat(user.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format du username invalide"})
		return
	}

	//Vérification du format de l'email
	if !pkg.ValidateEmailFormat(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format d'email invalide"})
		return
	}

	//Vérification de l'unicité de l'email
	if err := query.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne lors de la vérification de l'email"})
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

	c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("User %s créé avec succès", user.Username), "user": user})
}

// GetUser permet de récupérer un utilisateur par son ID
func GetUser(c *gin.Context) {
	id := c.Param("id")
	query := pkg.DB
	var user models.User

	if _, err := strconv.Atoi(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "L'ID doit être un entier valide"})
		return
	}

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

// GetUsers godoc
// @Summary Récupère tous les utilisateurs
// @Description Liste tous les utilisateurs existants
// @Tags Users
// @Produce json
// @Success 200 {object} map[string][]models.User "Liste des utilisateurs"
// @Failure 500 {object} map[string]string{"error": "Description de l'erreur"}
// @Router /users [get]

// GetUsers permet de récupérer tous les utilisateurs
func GetUsers(c *gin.Context) {
	var users []models.User
	query := pkg.DB

	if err := query.Find(&users).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateurs introuvables"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des users"})
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

	if _, err := strconv.Atoi(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "L'ID doit être un entier valide"})
		return
	}

	//Vérification de l'existence de l'utilisateur à mettre à jour dans la base de données
	if err := query.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur introuvable"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne lors de la récupération du user"})
		}
		return
	}

	var updatedUserData models.User

	//Lecture des données envoyées
	if err := c.ShouldBindJSON(&updatedUserData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format des données invalide"})
		return
	}

	// Vérification des champs obligatoires
	if updatedUserData.Username == "" || updatedUserData.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le username et l'email sont requis"})
		return
	}

	//Vérification du format du username
	if !pkg.ValidateUsernameFormat(updatedUserData.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format du username invalide"})
		return
	}

	//Vérification du format de l'email
	if !pkg.ValidateEmailFormat(updatedUserData.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format d'email invalide"})
		return
	}

	//Vérification de l'unicité du username (si modifié)
	if updatedUserData.Username != user.Username {
		log.Printf("Mise à jour du username:%s", updatedUserData.Username)
		log.Printf("username:%s", user.Username)
		fmt.Printf("Param ID: %s, Database User ID: %d\n", id, user.ID)

		var existingUser models.User

		if err := query.Debug().Where("username = ? AND id != ?", updatedUserData.Username, user.ID).First(&existingUser).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne lors de la vérification de l'unicité du username"})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ce username est déjà utilisé"})
			return
		}

		//Vérification de l'unicité de l'email (si modifié)
		if err := query.Where("email = ? AND id != ?", updatedUserData.Email, user.ID).First(&existingUser).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne lors de la vérification de l'unicité du username"})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cet email est déjà utilisé"})
			return
		}

		//Mise à jour des données de l'utilisateur
		user.Username = updatedUserData.Username
		user.Email = updatedUserData.Email

		//Sauvegarder dans la base de données
		if err := query.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne lors de la sauvegarde des mises à jour de l'utilisateur"})
			return
		}

		//Envoyer une réponse au client
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s mis à jour avec succès", user.Username)})
	}
}

// DeleteUser permet de supprimer un utilisateur
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	query := pkg.DB

	if _, err := strconv.Atoi(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "L'id doit être un entier valide"})
		return
	}

	var user models.User

	//Récupérer le user à supprimer
	if err := query.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur introuvable"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération de l'utilisateur à supprimer"})
		}
		return
	}

	// Supprimer l'utilisateur correspondant
	if err := query.Unscoped().Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression de l'utilisateur"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Utilisateur %s supprimé avec succès", user.Username)})
}
