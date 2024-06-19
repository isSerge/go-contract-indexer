package db

import (
	"database/sql"
	"log"
	"math/big"

	_ "github.com/lib/pq"
)

// DBInterface defines the methods that our database needs to implement
type DBInterface interface {
	SaveEvent(blockNumber uint64, txHash, eventType string, from, to, owner, spender *string, value *big.Int) error
}

// DB is a struct that holds the database connection
type DB struct {
	conn *sql.DB
}

// InitDB initializes the database connection and returns the DB instance
func InitDB(connStr string) DBInterface {
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	err = conn.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
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
