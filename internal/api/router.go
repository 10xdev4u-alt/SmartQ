package api

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Queue routes
		queues := v1.Group("/queues")
		{
			queues.POST("/", CreateQueue)
			// Other queue routes will go here
		}
	}

	return router
}
