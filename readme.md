# ERC-20 Contract Indexer

## Overview

A Golang application that indexes ERC-20 contract events and stores them in a
database.

## Features

- Watches for new events from an ERC-20 contract
- Stores event data in a PostgreSQL database
- Provides structured logging and error handling

## Requirements

- Go 1.16+
- Ethereum Node (RPC URL)

## Tools and Libraries

This project uses the following tools and libraries:

- [Logrus](https://github.com/sirupsen/logrus): A structured logger for Go, used
  for logging application events.
- [Viper](https://github.com/spf13/viper): A complete configuration solution for
  Go applications.
- [Go-Ethereum](https://github.com/ethereum/go-ethereum): Official Go
  implementation of the Ethereum protocol, used for interacting with the
  Ethereum blockchain.
- [Testify](https://github.com/stretchr/testify): A toolkit with common
  assertions and mocks for unit tests in Go.
- [GolangCI-Lint](https://github.com/golangci/golangci-lint): A fast Go linters
  runner for ensuring code quality.

## Setup

### Configuration

Create a `config.yaml` file in the project root:

```yaml
RPC_URL: 'your_rpc_url'
CONTRACT_ADDRESS: 'your_contract_address'
DB_CONN_STR: 'your_db_connection_string'
```
