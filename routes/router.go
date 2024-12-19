package routes

import (
	"net/http"
	"to-do-list-api/controllers"
	"to-do-list-api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	//Routes pour l'authentification
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", controllers.Register)
		authRoutes.POST("/login", controllers.Login)
		authRoutes.POST("/logout", middlewares.AuthRequired(), controllers.Logout)
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
		taskRoutes.GET("/", controllers.GetTasks)
		taskRoutes.POST("/", controllers.CreateTask)
		taskRoutes.PUT("/:id", middlewares.AuthorizeTaskOwnerShip(), controllers.UpdateTask)
		taskRoutes.DELETE("/:id", middlewares.AuthorizeTaskOwnerShip(), controllers.DeleteTask)
	}

	return router
}
