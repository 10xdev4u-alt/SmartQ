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
