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

// NewTicketRequest represents the data needed to create a new ticket.
type NewTicketRequest struct {
	CustomerName  string `json:"customer_name" binding:"required"`
	CustomerPhone string `json:"customer_phone" binding:"required"`
}

// CreateTicket handles the creation of a new ticket for a given queue.
func CreateTicket(db *storage.PostgresDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		queueIDStr := c.Param("queueId")
		queueID, err := uuid.Parse(queueIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid queue ID format"})
			return
		}

		var req NewTicketRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ticket, err := db.CreateTicket(c.Request.Context(), queueID, req.CustomerName, req.CustomerPhone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket"})
			return
		}

		c.JSON(http.StatusCreated, ticket)
	}
}

// updateTicketStatusHandler is a generic handler for updating a ticket's status.
func updateTicketStatusHandler(db *storage.PostgresDB, status string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ticketIDStr := c.Param("ticketId")
		ticketID, err := uuid.Parse(ticketIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID format"})
			return
		}

		ticket, err := db.UpdateTicketStatus(c.Request.Context(), ticketID, status)
		if err != nil {
			if err.Error() == fmt.Sprintf("ticket with ID %s not found", ticketID.String()) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket status"})
			return
		}

		c.JSON(http.StatusOK, ticket)
	}
}

// GetEstimatedWaitTime handles retrieving the estimated wait time for a given queue.
func GetEstimatedWaitTime(db *storage.PostgresDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		queueIDStr := c.Param("queueId")
		queueID, err := uuid.Parse(queueIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid queue ID format"})
			return
		}

		waitTime, err := db.CalculateEstimatedWaitTime(c.Request.Context(), queueID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate estimated wait time"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"estimated_wait_time_seconds": int(waitTime.Seconds())})
	}
}
