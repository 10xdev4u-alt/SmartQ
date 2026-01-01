const API_BASE_URL = 'http://localhost:8080/api/v1';
const WS_BASE_URL = 'ws://localhost:8080/ws'; // WebSocket URL
let currentQueueId = ''; // This should be set dynamically, e.g., from URL or config

document.addEventListener('DOMContentLoaded', () => {
    // For now, we'll hardcode a queue ID for testing.
    // In a real app, this would come from a configuration or URL parameter.
    currentQueueId = 'f400a87d-45c7-459c-b76a-aa7b7a68c822'; // Replace with a valid queue ID from your DB

    if (currentQueueId) {
        fetchQueueTickets(currentQueueId);
        setupWebSocket(); // Setup WebSocket instead of polling
    } else {
        document.getElementById('serving-ticket').textContent = 'Queue not selected.';
        document.getElementById('waiting-tickets').innerHTML = '<li>Queue not selected.</li>';
    }
});

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
        document.getElementById('serving-ticket').textContent = 'Error loading.';
        document.getElementById('waiting-tickets').innerHTML = '<li>Error loading tickets.</li>';
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
    const servingTicketSpan = document.getElementById('serving-ticket');
    const waitingTicketsList = document.getElementById('waiting-tickets');
    waitingTicketsList.innerHTML = ''; // Clear existing tickets

    let servingTicket = null;
    const waiting = [];

    tickets.forEach(ticket => {
        if (ticket.status === 'serving') {
            servingTicket = ticket;
        } else if (ticket.status === 'waiting') {
            waiting.push(ticket);
        }
    });

    if (servingTicket) {
        servingTicketSpan.textContent = servingTicket.ticket_number;
    } else {
        servingTicketSpan.textContent = '---';
    }

    if (waiting.length === 0) {
        waitingTicketsList.innerHTML = '<li>No one waiting.</li>';
    } else {
        waiting.forEach(ticket => {
            const listItem = document.createElement('li');
            listItem.textContent = ticket.ticket_number;
            waitingTicketsList.appendChild(listItem);
        });
    }
}
