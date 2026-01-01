package api

import (
	"github.com/gin-gonic/gin"
	"github.com/smartq/smartq/internal/notifier" // Import notifier
	"github.com/smartq/smartq/internal/storage"
)

// NewRouter sets up the Gin router and its routes.
func NewRouter(db *storage.PostgresDB, hub *notifier.Hub) *gin.Engine { // Accept hub
	router := gin.Default()

	// Serve static files for the staff dashboard
	router.Static("/staff", "./web/staff-dashboard")
	// Serve static files for the public display
	router.Static("/display", "./web/public-display")

	// WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		notifier.ServeWs(hub, c)
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Queue routes
		v1.POST("/queues", CreateQueue(db))
		v1.GET("/queues", GetQueues(db))
		v1.GET("/queues/:queueId", GetQueue(db))
		v1.GET("/queues/:queueId/tickets", GetTickets(db))
		v1.POST("/queues/:queueId/tickets", CreateTicket(db))
		v1.GET("/queues/:queueId/estimated-wait-time", GetEstimatedWaitTime(db))
		// Other queue routes will go here

		// Ticket routes
		tickets := v1.Group("/tickets")
		{
			tickets.POST("/:ticketId/call", updateTicketStatusHandler(db, "serving"))
			tickets.POST("/:ticketId/serve", updateTicketStatusHandler(db, "served"))
			tickets.POST("/:ticketId/cancel", updateTicketStatusHandler(db, "cancelled"))
		}
	}

	return router
}
