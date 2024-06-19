package main

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSubscription is a mock of ethereum.Subscription interface
type MockSubscription struct {
	mock.Mock
}

func (m *MockSubscription) Unsubscribe() {
	m.Called()
}

func (m *MockSubscription) Err() <-chan error {
	args := m.Called()
	return args.Get(0).(<-chan error)
}

func initTestConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func TestLoadConfig(t *testing.T) {
	initTestConfig()

	rpcURL := viper.GetString("RPC_URL")
	contractAddress := viper.GetString("CONTRACT_ADDRESS")
	dbConnStr := viper.GetString("DB_CONN_STR")

	assert.NotEmpty(t, rpcURL, "RPC_URL should not be empty")
	assert.NotEmpty(t, contractAddress, "CONTRACT_ADDRESS should not be empty")
	assert.NotEmpty(t, dbConnStr, "DB_CONN_STR should not be empty")
}

func TestLoadConfig_Error(t *testing.T) {
	// Clear any previously set environment variables
	viper.Reset()

	rpcURL := viper.GetString("RPC_URL")
	contractAddress := viper.GetString("CONTRACT_ADDRESS")
	dbConnStr := viper.GetString("DB_CONN_STR")

	assert.Empty(t, rpcURL, "RPC_URL should be empty")
	assert.Empty(t, contractAddress, "CONTRACT_ADDRESS should be empty")
	assert.Empty(t, dbConnStr, "DB_CONN_STR should be empty")
}
