const API_BASE_URL = 'http://localhost:8080/api/v1';
let currentQueueId = ''; // This should be set dynamically, e.g., from URL or config

document.addEventListener('DOMContentLoaded', () => {
    // For now, we'll hardcode a queue ID for testing.
    // In a real app, this would come from a configuration or URL parameter.
    currentQueueId = 'f400a87d-45c7-459c-b76a-aa7b7a68c822'; // Replace with a valid queue ID from your DB

    if (currentQueueId) {
        fetchQueueTickets(currentQueueId);
        // Set up polling for updates
        setInterval(() => fetchQueueTickets(currentQueueId), 3000); // Poll every 3 seconds
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
