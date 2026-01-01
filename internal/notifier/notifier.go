package notifier

import (
	"encoding/json"
	"log"
)

// Notifier is responsible for sending real-time updates to connected WebSocket clients.
type Notifier struct {
	hub *Hub
}

// NewNotifier creates a new Notifier instance.
func NewNotifier(hub *Hub) *Notifier {
	return &Notifier{
		hub: hub,
	}
}

// SendTicketUpdate sends a ticket update message to all connected WebSocket clients.
func (n *Notifier) SendTicketUpdate(ticket interface{}) {
	message, err := json.Marshal(map[string]interface{}{
		"type": "ticket_update",
		"data": ticket,
	})
	if err != nil {
		log.Printf("Error marshalling ticket update: %v", err)
		return
	}
	n.hub.broadcast <- message
}

// SendQueueUpdate sends a queue update message to all connected WebSocket clients.
func (n *Notifier) SendQueueUpdate(queue interface{}) {
	message, err := json.Marshal(map[string]interface{}{
		"type": "queue_update",
		"data": queue,
	})
	if err != nil {
		log.Printf("Error marshalling queue update: %v", err)
		return
	}
	n.hub.broadcast <- message
}