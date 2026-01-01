# SmartQ: Smart Queue Ticketing System

SmartQ is a modern, efficient queue management system designed for small clinics, salons, and other businesses that rely on orderly customer flow. It replaces the chaotic "who's next?" shout with a simple, digital solution.

## MVP Feature Details

### 1. Customer Onboarding
- A unique QR code is physically displayed at the business premises.
- Scanning the QR code opens a simple, mobile-friendly web page.
- The page has a form with fields for "Name" and "Phone Number".
- On submission, the system validates the phone number format.
- A new ticket is created with a unique, sequential ticket number (e.g., A-101, A-102).
- The customer is shown a confirmation page with their ticket number and their position in the queue (e.g., "You are number 5 in the queue").

### 2. Real-time Queue Display
- A web page designed for a tablet or TV screen.
- It displays an ordered list of ticket numbers currently in the queue.
- The ticket number currently being served is clearly highlighted.
- The display updates in near real-time without requiring a manual refresh (e.g., using WebSockets or polling).

### 3. Staff Queue Management
- A simple, password-protected web page for staff.
- Displays the current queue.
- A "Call Next" button that:
    - Moves the top customer from "waiting" to "serving" status.
    - Triggers an SMS notification to that customer.
- A "Mark as Served" button that:
    - Removes the currently serving customer from the queue.
- A "Cancel" button to remove a waiting customer from the queue.

### 4. SMS Notifications
- When a staff member clicks "Call Next", the customer with that ticket receives an SMS notification (e.g., "It's your turn! Please proceed to the counter.").
- (Optional) When a customer moves to the second position in the queue, they receive a "You're next in line!" preparation SMS.

## Future Goals

- **Priority Queueing:** Implement rules to give certain customers priority.
- **Staff Dashboard:** A web-based dashboard for staff to manage the queue (e.g., call next, remove customer).
- **Wait-Time Estimation:** Provide customers with an estimated wait time based on historical data.
- **Analytics:** Track queue metrics to help businesses optimize their operations.
- **CLI & GUI:** A command-line and graphical user interface for administration and management.

## Technology Stack

- **Backend:** Go
- **Database:** (To be determined: likely PostgreSQL or SQLite for self-hosting)
- **Frontend (Dashboard/Display):** (To be determined: likely a modern JavaScript framework)
- **Notifications:** SMS Gateway (e.g., Twilio)
- **Deployment:** Docker for self-hostability
