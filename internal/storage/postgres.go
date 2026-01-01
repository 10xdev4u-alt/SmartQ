package storage

import (
	"context"
	"fmt"
	"log"
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
