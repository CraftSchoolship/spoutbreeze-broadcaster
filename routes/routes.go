package routes

import (
	"github.com/gin-gonic/gin"
	"spoutbreeze/controllers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	broadcasterGroup := router.Group("/broadcaster")
	{
		broadcasterGroup.POST("/joinBBB", controllers.JoinBBB)
	}

	return router
}
