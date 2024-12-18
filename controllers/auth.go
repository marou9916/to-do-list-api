package controllers

import (
	"net/http"
	"time"
	"to-do-list-api/models"
	"to-do-list-api/pkg"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Register permet d'enregistrer un nouvel utilisateur

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required, email"`
		Password string `json:"password" binding:"required, min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du hachage du mot de passe"})
		return
	}

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := pkg.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de l'utilisateur"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Inscription réalisée avec succès"})
}

// Login permet de se connecter
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required, email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	if err := pkg.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email incorrect"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne"})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Mot de passe incorrect"})
		return
	}

	//Générer un token de session unique
	sessionToken := pkg.GenerateToken()
	expiration := time.Now().Add(24 * time.Hour)

	//Créer une nouvelle session
	session := models.Session{
		Token:     sessionToken,
		UserID:    user.ID,
		ExpiresAt: expiration,
	}
	if err := pkg.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Erreur lors de la création de la session"})
		return
	}

	//Configurer un cookie sécurisé
	c.SetCookie("session_token", session.Token, int(24*time.Hour.Seconds()), "/", "", true, false)

	c.JSON(http.StatusOK, gin.H{"message": "Connexion réussie"})
}

// Logout permet de se déconnecter
func Logout(c *gin.Context) {
	c.SetCookie("session_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Déconnexion réussie"})
}
