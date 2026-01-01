# Architecture

The SmartQ system is designed as a set of cooperating services, making it scalable, maintainable, and easy to deploy.

## Technology Choices

- **Backend Framework:** [Gin](https://gin-gonic.com/) - A high-performance, minimalist web framework for Go.
- **Real-time Communication:** [Gorilla WebSocket](https://github.com/gorilla/websocket) - For pushing real-time updates to the Public Display and Staff Dashboard.
- **Database:**
    - **Primary:** PostgreSQL - A powerful, open-source object-relational database system.
    - **Self-hosted/Embedded Option:** SQLite - For simple, single-file database deployment.
    - **Driver:** [pgx](https://github.com/jackc/pgx) for PostgreSQL.
- **Frontend:** Simple HTML, CSS, and JavaScript for the MVP. This keeps the initial version lightweight and avoids framework overhead. We can upgrade to a framework like React or Vue.js later if needed.
- **Deployment:** Docker and Docker Compose for containerization and service orchestration.

## Components

1.  **Queue Service (Go):** The core of the system, built with Gin. It manages the queue, customer data, and business logic. It exposes:
    - A REST API for the customer onboarding and staff management actions.
    - A WebSocket endpoint for real-time updates.

2.  **Database (PostgreSQL/SQLite):** Stores queue information, customer tickets, and historical data for wait-time estimation. The storage layer in our Go application will be designed to abstract away the specific database implementation.

3.  **Public Display App (Web):** A simple HTML/JavaScript page that connects to the Queue Service's WebSocket endpoint to receive and display real-time queue updates.

4.  **Customer Onboarding App (Web):** A lightweight HTML/CSS/JS single-page application that allows customers to scan a QR code and submit their details to the Queue Service's REST API.

5.  **Notification Service (Go):** A module within the main Go application that listens for queue events (e.g., "customer called") and dispatches SMS notifications via a third-party gateway like Twilio.

6.  **Staff Dashboard (Web):** A password-protected HTML/CSS/JS single-page application that interacts with the REST API (for actions like "call next") and connects to the WebSocket endpoint for real-time queue visualization.

7.  **CLI (Go):** A command-line interface for administrative tasks.

## Data Flow (MVP)

1.  A customer scans a QR code, which leads to the **Customer Onboarding App**.
2.  The customer enters their phone number and name.
3.  The app sends a `POST` request to the **Queue Service** REST API.
4.  The **Queue Service** validates the data, adds the customer to the queue in the **Database**, and assigns a ticket number.
5.  The **Public Display App** and **Staff Dashboard** receive a real-time queue update via their WebSocket connection.
6.  A staff member clicks "Call Next" on the **Staff Dashboard**.
7.  The dashboard sends a `POST` request to the **Queue Service** REST API.
8.  The **Queue Service** updates the customer's status in the database, which triggers two actions:
    a. It broadcasts a WebSocket message to all connected clients (Display and Dashboard) to show the customer is now being served.
    b. It triggers the **Notification Service** to send an SMS to the customer.
