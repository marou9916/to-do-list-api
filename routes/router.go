package routes

import (
	"net/http"
	"to-do-list-api/controllers"
	"to-do-list-api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Aucune confiance envers les proxies
	err := router.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}

	//Routes pour l'authentification
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", controllers.Register)
		authRoutes.POST("/login", controllers.Login)
		authRoutes.POST("/logout", middlewares.AuthRequired(), controllers.Logout)
		authRoutes.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Inscription ? Connexion ? Ou déconnexion ?"})
		})
	}

	//Routes pour les utilisateurs
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/", controllers.GetUsers)
		userRoutes.GET("/:id", controllers.GetUser)
		userRoutes.PUT("/:id", controllers.UpdateUser)
		userRoutes.POST("/", controllers.CreateUser)
		userRoutes.DELETE("/:id", controllers.DeleteUser)
		userRoutes.DELETE("/", func(c *gin.Context) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "L'ID est requis pour cette opération"})
		})
	}

	//Routes pour les tâches
	taskRoutes := router.Group("/tasks")
	taskRoutes.Use(middlewares.AuthRequired())
	{
		taskRoutes.GET("/", middlewares.AuthorizeTaskOwnerShip(), controllers.GetTasks)
		taskRoutes.POST("/", middlewares.AuthorizeTaskOwnerShip(), controllers.CreateTask)
		taskRoutes.PUT("/:id",middlewares.AuthorizeTaskOwnerShip(), controllers.UpdateTask)
		taskRoutes.DELETE("/:id",middlewares.AuthorizeTaskOwnerShip(), controllers.DeleteTask)
	}

	return router
}
