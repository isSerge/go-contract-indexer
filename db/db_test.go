package db

import (
	"database/sql"
	"log"
	"math/big"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

// InitTestDB initializes an in-memory SQLite database for testing
func InitTestDB() *sql.DB {
	conn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Failed to connect to the SQLite database: %v", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS erc20_events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		block_number BIGINT NOT NULL,
		tx_hash TEXT NOT NULL,
		event_type TEXT NOT NULL,
		from_address TEXT,
		to_address TEXT,
		owner_address TEXT,
		spender_address TEXT,
		value TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = conn.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	return conn
}

func TestInitDB(t *testing.T) {
	conn := InitTestDB()
	assert.NotNil(t, conn, "Database connection should not be nil")
}

func TestSaveEvent(t *testing.T) {
	conn := InitTestDB()
	db := &DB{conn: conn}

	blockNumber := uint64(123456)
	txHash := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	eventType := "Transfer"
	from := "0x1234567890abcdef1234567890abcdef12345678"
	to := "0x1234567890abcdef1234567890abcdef12345679"
	value := big.NewInt(1000)

	err := db.SaveEvent(blockNumber, txHash, eventType, &from, &to, nil, nil, value)
	assert.Nil(t, err, "Error should be nil")

	// Add assertions to check if the event is correctly saved in the database
	var count int
	err = conn.QueryRow(`SELECT COUNT(*) FROM erc20_events WHERE tx_hash = ?`, txHash).Scan(&count)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, 1, count, "Event should be saved in the database")
}
