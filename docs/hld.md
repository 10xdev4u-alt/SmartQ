# High-Level Design (HLD)

This document outlines the high-level design of the SmartQ system.

## Database Schema

The database is the backbone of the SmartQ system, responsible for persisting the state of the queues and tickets.

### `queues` table

This table stores the different queues that a business might operate. For the MVP, a business will likely have only one queue, but this design allows for future expansion.

- **`id`**: `UUID`, Primary Key - Unique identifier for the queue.
- **`name`**: `VARCHAR(255)`, Not Null - A human-readable name for the queue (e.g., "Main Counter", "Dr. Smith's Office").
- **`created_at`**: `TIMESTAMP`, Not Null - Timestamp of when the queue was created.

**Example:**
| id                                   | name          | created_at          |
| ------------------------------------ | ------------- | ------------------- |
| `a1b2c3d4-e5f6-7890-1234-567890abcdef` | "Main Counter"| `2026-01-01 10:00:00` |


### `tickets` table

This table stores the individual tickets for each customer in a queue.

- **`id`**: `UUID`, Primary Key - Unique identifier for the ticket.
- **`queue_id`**: `UUID`, Foreign Key to `queues.id` - The queue this ticket belongs to.
- **`customer_name`**: `VARCHAR(255)`, Not Null - The name of the customer.
- **`customer_phone`**: `VARCHAR(50)`, Not Null - The phone number of the customer, used for SMS notifications.
- **`ticket_number`**: `VARCHAR(10)`, Not Null - The public-facing ticket number (e.g., "A-101"). This should be unique per queue for a given day.
- **`status`**: `VARCHAR(20)`, Not Null - The current status of the ticket.
    - `waiting`: The customer is in the queue.
    - `serving`: The customer is currently being served.
    - `served`: The customer has been served.
    - `cancelled`: The customer has left the queue.
- **`position`**: `INTEGER`, Not Null - The position of the ticket in the queue. This will be managed by the application logic.
- **`created_at`**: `TIMESTAMP`, Not Null - Timestamp of when the ticket was created.
- **`updated_at`**: `TIMESTAMP`, Not Null - Timestamp of the last update to the ticket.

**Example:**
| id                                   | queue_id                             | customer_name | customer_phone | ticket_number | status  | position | created_at          | updated_at          |
| ------------------------------------ | ------------------------------------ | ------------- | -------------- | ------------- | ------- | -------- | ------------------- | ------------------- |
| `b2c3d4e5-f6a7-8901-2345-67890abcdef1` | `a1b2c3d4-e5f6-7890-1234-567890abcdef` | "John Doe"    | "+15551234567" | "A-101"       | `waiting` | 1        | `2026-01-01 10:05:00` | `2026-01-01 10:05:00` |
| `c3d4e5f6-a7b8-9012-3456-7890abcdef2` | `a1b2c3d4-e5f6-7890-1234-567890abcdef` | "Jane Smith"  | "+15557654321" | "A-102"       | `waiting` | 2        | `2026-01-01 10:06:00` | `2026-01-01 10:06:00` |
