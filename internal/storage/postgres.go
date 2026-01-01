package storage

import (
	"context"
	"fmt"
	"log"
	"strconv" // Add strconv import
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/smartq/smartq/internal/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5" // Import pgx for pgx.ErrNoRows
)

// Queue represents a queue in the database.
type Queue struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// Ticket represents a ticket in the database.
type Ticket struct {
	ID           uuid.UUID `json:"id"`
	QueueID      uuid.UUID `json:"queue_id"`
	CustomerName string    `json:"customer_name"`
	CustomerPhone string   `json:"customer_phone"`
	TicketNumber string    `json:"ticket_number"`
	Status       string    `json:"status"`
	Position     int       `json:"position"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TicketHistory represents a status change event for a ticket.
type TicketHistory struct {
	ID        uuid.UUID `json:"id"`
	TicketID  uuid.UUID `json:"ticket_id"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"` // Redundant but for consistency
}

type PostgresDB struct {
	pool *pgxpool.Pool
}

func NewPostgresDB(cfg *config.Config) (*PostgresDB, error) {
	connConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), connConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Ping the database to verify the connection
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL!")
	return &PostgresDB{pool: pool}, nil
}

func (db *PostgresDB) Close() {
	db.pool.Close()
	log.Println("PostgreSQL connection pool closed.")
}

// CreateQueue inserts a new queue into the database.
func (db *PostgresDB) CreateQueue(ctx context.Context, name string) (*Queue, error) {
	queue := &Queue{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
	}

	query := `INSERT INTO queues (id, name, created_at) VALUES ($1, $2, $3) RETURNING id, name, created_at`
	err := db.pool.QueryRow(ctx, query, queue.ID, queue.Name, queue.CreatedAt).Scan(&queue.ID, &queue.Name, &queue.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert queue: %w", err)
	}

	return queue, nil
}

// GetQueueByID retrieves a queue from the database by its ID.
func (db *PostgresDB) GetQueueByID(ctx context.Context, id uuid.UUID) (*Queue, error) {
	queue := &Queue{}
	query := `SELECT id, name, created_at FROM queues WHERE id = $1`
	err := db.pool.QueryRow(ctx, query, id).Scan(&queue.ID, &queue.Name, &queue.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("queue with ID %s not found", id.String())
		}
		return nil, fmt.Errorf("failed to get queue by ID: %w", err)
	}
	return queue, nil
}

// GetTicketsByQueueID retrieves all tickets for a given queue ID.
func (db *PostgresDB) GetTicketsByQueueID(ctx context.Context, queueID uuid.UUID) ([]*Ticket, error) {
	var tickets []*Ticket
	query := `SELECT id, queue_id, customer_name, customer_phone, ticket_number, status, position, created_at, updated_at
			  FROM tickets
			  WHERE queue_id = $1
			  ORDER BY position ASC, created_at ASC`
	rows, err := db.pool.Query(ctx, query, queueID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tickets: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		ticket := &Ticket{}
		err := rows.Scan(
			&ticket.ID,
			&ticket.QueueID,
			&ticket.CustomerName,
			&ticket.CustomerPhone,
			&ticket.TicketNumber,
			&ticket.Status,
			&ticket.Position,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ticket row: %w", err)
		}
		tickets = append(tickets, ticket)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return tickets, nil
}

// CreateTicket inserts a new ticket into the database.
func (db *PostgresDB) CreateTicket(ctx context.Context, queueID uuid.UUID, customerName, customerPhone string) (*Ticket, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // Rollback on error, commit on success

	// Get the next ticket number and position
	var lastTicketNumber string
	var lastPosition int
	
	// Query for the last ticket number and position for the given queue on the current day
	// This query needs to be robust to handle cases where there are no tickets yet
	// For simplicity, let's assume ticket numbers are sequential per day per queue
	// and positions are also sequential.
	
	// Get last ticket number
	err = tx.QueryRow(ctx, `
		SELECT COALESCE(MAX(ticket_number), 'A-000')
		FROM tickets
		WHERE queue_id = $1 AND created_at::date = NOW()::date`, queueID).Scan(&lastTicketNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get last ticket number: %w", err)
	}

	// Get last position
	err = tx.QueryRow(ctx, `
		SELECT COALESCE(MAX(position), 0)
		FROM tickets
		WHERE queue_id = $1`, queueID).Scan(&lastPosition)
	if err != nil {
		return nil, fmt.Errorf("failed to get last position: %w", err)
	}

	// Generate next ticket number
	prefix := string(lastTicketNumber[0]) // Assuming 'A'
	numStr := lastTicketNumber[2:]        // Assuming '000'
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ticket number: %w", err)
	}
	nextTicketNumber := fmt.Sprintf("%s-%03d", prefix, num+1)
	nextPosition := lastPosition + 1

	ticket := &Ticket{
		ID:           uuid.New(),
		QueueID:      queueID,
		CustomerName: customerName,
		CustomerPhone: customerPhone,
		TicketNumber: nextTicketNumber,
		Status:       "waiting", // Default status
		Position:     nextPosition,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	query := `INSERT INTO tickets (id, queue_id, customer_name, customer_phone, ticket_number, status, position, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, queue_id, customer_name, customer_phone, ticket_number, status, position, created_at, updated_at`
	err = tx.QueryRow(ctx, query,
		ticket.ID,
		ticket.QueueID,
		ticket.CustomerName,
		ticket.CustomerPhone,
		ticket.TicketNumber,
		ticket.Status,
		ticket.Position,
		ticket.CreatedAt,
		ticket.UpdatedAt,
	).Scan(
		&ticket.ID,
		&ticket.QueueID,
		&ticket.CustomerName,
		&ticket.CustomerPhone,
		&ticket.TicketNumber,
		&ticket.Status,
		&ticket.Position,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert ticket: %w", err)
	}

	// Log the initial status change
	if err := LogTicketStatusChange(ctx, tx, ticket.ID, ticket.Status); err != nil {
		return nil, fmt.Errorf("failed to log initial ticket status: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ticket, nil
}

// UpdateTicketStatus updates the status of a ticket.
func (db *PostgresDB) UpdateTicketStatus(ctx context.Context, ticketID uuid.UUID, status string) (*Ticket, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	ticket := &Ticket{}
	query := `UPDATE tickets SET status = $1, updated_at = NOW() WHERE id = $2 RETURNING id, queue_id, customer_name, customer_phone, ticket_number, status, position, created_at, updated_at`
	err = tx.QueryRow(ctx, query, status, ticketID).Scan(
		&ticket.ID,
		&ticket.QueueID,
		&ticket.CustomerName,
		&ticket.CustomerPhone,
		&ticket.TicketNumber,
		&ticket.Status,
		&ticket.Position,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("ticket with ID %s not found", ticketID.String())
		}
		return nil, fmt.Errorf("failed to update ticket status: %w", err)
	}

	// Log the status change
	if err := LogTicketStatusChange(ctx, tx, ticket.ID, ticket.Status); err != nil {
		return nil, fmt.Errorf("failed to log ticket status change: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ticket, nil
}

// LogTicketStatusChange records a ticket's status change in the ticket_history table.
func LogTicketStatusChange(ctx context.Context, tx pgx.Tx, ticketID uuid.UUID, status string) error {
	history := &TicketHistory{
		ID:        uuid.New(),
		TicketID:  ticketID,
		Status:    status,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	query := `INSERT INTO ticket_history (id, ticket_id, status, timestamp, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := tx.Exec(ctx, query, history.ID, history.TicketID, history.Status, history.Timestamp, history.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to log ticket status change: %w", err)
	}
	return nil
}

// CalculateEstimatedWaitTime calculates the estimated wait time for a given queue.
// It does this by looking at the average time tickets spent in 'waiting' status
// for recently served tickets in that queue.
func (db *PostgresDB) CalculateEstimatedWaitTime(ctx context.Context, queueID uuid.UUID) (time.Duration, error) {
	// For simplicity, let's consider the last 10 served tickets to calculate average waiting time.
	// A more sophisticated algorithm might consider time of day, day of week, etc.
	query := `
		SELECT
			EXTRACT(EPOCH FROM (th_served.timestamp - th_waiting.timestamp)) AS waiting_duration
		FROM
			tickets t
		JOIN
			ticket_history th_waiting ON t.id = th_waiting.ticket_id AND th_waiting.status = 'waiting'
		JOIN
			ticket_history th_served ON t.id = th_served.ticket_id AND th_served.status = 'served'
		WHERE
			t.queue_id = $1
			AND t.status = 'served' -- Only consider served tickets for historical data
		ORDER BY
			th_served.timestamp DESC
		LIMIT 10
	`

	rows, err := db.pool.Query(ctx, query, queueID)
	if err != nil {
		return 0, fmt.Errorf("failed to query waiting durations: %w", err)
	}
	defer rows.Close()

	var totalDurationSeconds float64
	var count int
	for rows.Next() {
		var duration float64
		if err := rows.Scan(&duration); err != nil {
			return 0, fmt.Errorf("failed to scan waiting duration: %w", err)
		}
		totalDurationSeconds += duration
		count++
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("error after iterating rows: %w", err)
	}

	if count == 0 {
		return 0, nil // No historical data, so 0 wait time
	}

	averageDurationSeconds := totalDurationSeconds / float64(count)
	return time.Duration(averageDurationSeconds) * time.Second, nil
}
