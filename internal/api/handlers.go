package api

import (
	"fmt" // Import fmt package for error formatting
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid" // Import uuid package
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

// GetQueue handles retrieving a queue by its ID.
func GetQueue(db *storage.PostgresDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		queueIDStr := c.Param("queueId")
		queueID, err := uuid.Parse(queueIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid queue ID format"})
			return
		}

		queue, err := db.GetQueueByID(c.Request.Context(), queueID)
		if err != nil {
			if err.Error() == fmt.Sprintf("queue with ID %s not found", queueID.String()) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve queue"})
			return
		}

		c.JSON(http.StatusOK, queue)
	}
}

// GetTickets handles retrieving all tickets for a given queue ID.
func GetTickets(db *storage.PostgresDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		queueIDStr := c.Param("queueId")
		queueID, err := uuid.Parse(queueIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid queue ID format"})
			return
		}

		tickets, err := db.GetTicketsByQueueID(c.Request.Context(), queueID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tickets"})
			return
		}

		c.JSON(http.StatusOK, tickets)
	}
}
