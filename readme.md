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

## Setup

### Configuration

Create a `config.yaml` file in the project root:

```yaml
RPC_URL: 'your_rpc_url'
CONTRACT_ADDRESS: 'your_contract_address'
DB_CONN_STR: 'your_db_connection_string'
```
