package api

import (
	"github.com/gin-gonic/gin"
	"github.com/smartq/smartq/internal/storage"
)

func NewRouter(db *storage.PostgresDB) *gin.Engine {
	router := gin.Default()

	// Serve static files for the staff dashboard
	router.Static("/staff", "./web/staff-dashboard")
	// Serve static files for the public display
	router.Static("/display", "./web/public-display")

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Queue routes
		v1.POST("/queues", CreateQueue(db))
		v1.GET("/queues/:queueId", GetQueue(db))
		v1.GET("/queues/:queueId/tickets", GetTickets(db))
		v1.POST("/queues/:queueId/tickets", CreateTicket(db))
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
