package routes

import (
	"to-do-list-api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	//Routes pour les utilisateurs
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/", controllers.CreateUser)
		userRoutes.DELETE("/:id", controllers.DeleteUser)
	}

	//Routes pour les tâches
	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.GET("/", controllers.GetTasks)
		taskRoutes.POST("/", controllers.CreateTask)
		taskRoutes.PUT("/:id", controllers.UpdateTask)
		taskRoutes.DELETE("/:id", controllers.DeleteTask)
	}

	return router
}
