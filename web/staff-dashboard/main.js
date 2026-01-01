const API_BASE_URL = 'http://localhost:8080/api/v1';
const WS_BASE_URL = 'ws://localhost:8080/ws'; // WebSocket URL
let currentQueueId = ''; // This should be set dynamically, e.g., from URL or config

document.addEventListener('DOMContentLoaded', () => {
    // For now, we'll hardcode a queue ID for testing.
    // In a real app, this would come from a login or selection process.
    currentQueueId = 'f400a87d-45c7-459c-b76a-aa7b7a68c822'; // Replace with a valid queue ID from your DB

    if (currentQueueId) {
        fetchQueueDetails(currentQueueId);
        fetchQueueTickets(currentQueueId);
        setupWebSocket(); // Setup WebSocket instead of polling
    } else {
        document.getElementById('queue-info').textContent = 'Please select a queue.';
    }
});

async function fetchQueueDetails(queueId) {
    try {
        const response = await fetch(`${API_BASE_URL}/queues/${queueId}`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const queue = await response.json();
        document.getElementById('queue-info').innerHTML = `
            <p><strong>Queue Name:</strong> ${queue.name}</p>
            <p><strong>Queue ID:</strong> ${queue.id}</p>
        `;
    } catch (error) {
        console.error('Error fetching queue details:', error);
        document.getElementById('queue-info').textContent = 'Failed to load queue details.';
    }
}

async function fetchQueueTickets(queueId) {
    try {
        const response = await fetch(`${API_BASE_URL}/queues/${queueId}/tickets`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const tickets = await response.json();
        renderTickets(tickets);
    } catch (error) {
        console.error('Error fetching queue tickets:', error);
        document.getElementById('tickets').innerHTML = '<li>Failed to load tickets.</li>';
    }
}

function setupWebSocket() {
    const socket = new WebSocket(WS_BASE_URL);

    socket.onopen = (event) => {
        console.log('WebSocket connected:', event);
    };

    socket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        console.log('WebSocket message received:', message);

        if (message.type === 'ticket_update') {
            // Re-fetch all tickets for the current queue to ensure consistency
            // A more optimized approach would be to update individual tickets in the DOM
            // but for simplicity, re-fetching is sufficient for now.
            fetchQueueTickets(currentQueueId);
        } else if (message.type === 'queue_update') {
            // If a queue update is received, and it's relevant to our current queue,
            // we might want to re-fetch queue details or tickets.
            // For now, we'll just log it.
            console.log('Queue update received:', message.data);
            fetchQueueDetails(currentQueueId); // Re-fetch queue details on update
        }
    };

    socket.onclose = (event) => {
        console.log('WebSocket disconnected:', event);
        // Attempt to reconnect after a delay
        setTimeout(setupWebSocket, 3000);
    };

    socket.onerror = (error) => {
        console.error('WebSocket error:', error);
    };
}

function renderTickets(tickets) {
    const ticketsList = document.getElementById('tickets');
    ticketsList.innerHTML = ''; // Clear existing tickets

    if (tickets.length === 0) {
        ticketsList.innerHTML = '<li>No tickets in queue.</li>';
        return;
    }

    tickets.forEach(ticket => {
        const listItem = document.createElement('li');
        listItem.className = ticket.status; // Add status as class for styling
        listItem.innerHTML = `
            <div>
                <strong>${ticket.ticket_number}</strong> - ${ticket.customer_name} (${ticket.customer_phone})
                <br>
                Status: ${ticket.status}
            </div>
            <div class="ticket-actions">
                ${ticket.status === 'waiting' ? `<button onclick="callTicket('${ticket.id}')">Call</button>` : ''}
                ${ticket.status === 'serving' ? `<button onclick="serveTicket('${ticket.id}')">Serve</button>` : ''}
                ${ticket.status === 'waiting' || ticket.status === 'serving' ? `<button class="cancel" onclick="cancelTicket('${ticket.id}')">Cancel</button>` : ''}
            </div>
        `;
        ticketsList.appendChild(listItem);
    });
}

async function callTicket(ticketId) {
    await updateTicketStatus(ticketId, 'call');
}

async function serveTicket(ticketId) {
    await updateTicketStatus(ticketId, 'serve');
}

async function cancelTicket(ticketId) {
    await updateTicketStatus(ticketId, 'cancel');
}

async function updateTicketStatus(ticketId, action) {
    try {
        const response = await fetch(`${API_BASE_URL}/tickets/${ticketId}/${action}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
        });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        // No need to refresh tickets after update, WebSocket will handle it
        // fetchQueueTickets(currentQueueId);
    } catch (error) {
        console.error(`Error updating ticket status (${action}):`, error);
        alert(`Failed to ${action} ticket.`);
    }
}
