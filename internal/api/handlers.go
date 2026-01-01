package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewQueue represents the data needed to create a new queue.
type NewQueue struct {
	Name string `json:"name" binding:"required"`
}

// CreateQueue handles the creation of a new queue.
func CreateQueue(c *gin.Context) {
	var newQueue NewQueue
	if err := c.ShouldBindJSON(&newQueue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real implementation, we would save the new queue to the database here.
	// For now, we'll just return a dummy response.

	c.JSON(http.StatusCreated, gin.H{
		"id":         "a1b2c3d4-e5f6-7890-1234-567890abcdef", // Dummy UUID
		"name":       newQueue.Name,
		"created_at": "2026-01-01T10:00:00Z",              // Dummy timestamp
	})
}
