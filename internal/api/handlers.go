package api

import (
	"net/http"
	"time" // Import time package

	"github.com/gin-gonic/gin"
	"github.com/smartq/smartq/internal/storage"
)

// NewQueue represents the data needed to create a new queue.
type NewQueue struct {
	Name string `json:"name" binding:"required"`
}

// CreateQueue handles the creation of a new queue.
// It now returns a gin.HandlerFunc closure to capture the *storage.PostgresDB instance.
func CreateQueue(db *storage.PostgresDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newQueue NewQueue
		if err := c.ShouldBindJSON(&newQueue); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Use the storage layer to create the queue
		queue, err := db.CreateQueue(c.Request.Context(), newQueue.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create queue"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":         queue.ID,
			"name":       queue.Name,
			"created_at": queue.CreatedAt.Format(time.RFC3339), // Format time for JSON
		})
	}
}
