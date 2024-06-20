package db

import (
	"database/sql"
	"log"
	"math/big"
	"time"

	// import postgres driver
	_ "github.com/lib/pq"
)

// Interface defines the methods that our database needs to implement
type Interface interface {
	SaveEvent(blockNumber uint64, txHash, eventType string, from, to, owner, spender *string, value *big.Int) error
	Close() error
}

// DB is a struct that holds the database connection
type DB struct {
	conn *sql.DB
}

// InitDB initializes the database connection with retry mechanism and returns the DB instance
func InitDB(connStr string) Interface {
	conn, err := connectWithRetry(connStr, 10, 2*time.Second)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	log.Println("Database connected successfully")

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS erc20_events (
		id SERIAL PRIMARY KEY,
		block_number BIGINT NOT NULL,
		tx_hash VARCHAR(66) NOT NULL,
		event_type VARCHAR(50) NOT NULL,
		from_address VARCHAR(42),
		to_address VARCHAR(42),
		owner_address VARCHAR(42),
		spender_address VARCHAR(42),
		value NUMERIC,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = conn.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Println("Table erc20_events exists or created successfully")

	return &DB{conn: conn}
}

// connectWithRetry attempts to connect to the database with retries.
func connectWithRetry(dsn string, retries int, delay time.Duration) (*sql.DB, error) {
	var db *sql.DB
	var err error
	for i := 0; i < retries; i++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Failed to connect to the database: %v. Retrying...", err)
			time.Sleep(delay)
			continue
		}

		err = db.Ping()
		if err == nil {
			return db, nil
		}

		log.Printf("Failed to ping the database: %v. Retrying...", err)
		time.Sleep(delay)
	}

	return nil, err
}

// SaveEvent saves an indexed event to the database
func (db *DB) SaveEvent(blockNumber uint64, txHash, eventType string, from, to, owner, spender *string, value *big.Int) error {
	var valueStr *string
	if value != nil {
		v := value.String()
		valueStr = &v
	}
	_, err := db.conn.Exec(`
		INSERT INTO erc20_events (block_number, tx_hash, event_type, from_address, to_address, owner_address, spender_address, value)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, blockNumber, txHash, eventType, from, to, owner, spender, valueStr)
	return err
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}
