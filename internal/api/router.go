package api

import (
	"github.com/gin-gonic/gin"
	"github.com/smartq/smartq/internal/storage"
)

func NewRouter(db *storage.PostgresDB) *gin.Engine {
	router := gin.Default()

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Queue routes
		v1.POST("/queues", CreateQueue(db))
		v1.GET("/queues/:queueId", GetQueue(db))
		// Other queue routes will go here


	}

	return router
}
