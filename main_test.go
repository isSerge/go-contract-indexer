package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Set up environment variables for testing
	os.Setenv("RPC_URL", "http://localhost:8545")
	os.Setenv("CONTRACT_ADDRESS", "0x1234567890abcdef1234567890abcdef12345678")
	os.Setenv("DB_CONN_STR", "postgres://user:password@localhost:5432/dbname?sslmode=disable")

	rpcURL, contractAddress, dbConnStr, err := loadConfig()

	assert.NoError(t, err)
	assert.Equal(t, "http://localhost:8545", rpcURL)
	assert.Equal(t, "0x1234567890abcdef1234567890abcdef12345678", contractAddress)
	assert.Equal(t, "postgres://user:password@localhost:5432/dbname?sslmode=disable", dbConnStr)
}

func TestLoadConfig_Error(t *testing.T) {
	// Set environment variables to empty strings to test error case
	os.Setenv("RPC_URL", "")
	os.Setenv("CONTRACT_ADDRESS", "")
	os.Setenv("DB_CONN_STR", "")

	_, _, _, err := loadConfig()

	assert.Error(t, err)
	assert.Equal(t, "RPC_URL, CONTRACT_ADDRESS, or DB_CONN_STR is not set in the .env file", err.Error())
}
