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

Update `config.yaml` file in the project root:

```yaml
RPC_URL: 'your_rpc_url'
CONTRACT_ADDRESS: 'your_contract_address'
DB_CONN_STR: 'your_db_connection_string'
```

## Docker

This project uses Docker to containerize the application and Docker Compose to
manage multi-container environments.

### Prerequisites

- Docker installed on your local machine.
  [Docker Installation Guide](https://docs.docker.com/get-docker/)
- Docker Compose installed on your local machine.
  [Docker Compose Installation Guide](https://docs.docker.com/compose/install/)

### Building and Running the Docker Containers

1. **Build the Docker images and start the containers:**

```sh
docker-compose up --build
```

This command builds the Docker images defined in the
[docker-compose.yml](docker-compose.yml) file and starts the containers.

2. **Verify the Containers are Running:**

After running the above command, you should see logs from your application
container indicating that it has successfully connected to the PostgreSQL
database and is running.

3. **Stopping the Containers:**

To stop the running containers, use:

```sh
docker-compose down
```

### Docker Compose

Ensure your [docker-compose.yml](docker-compose.yml) file includes the volume
for [config.yml](config.yaml) and sets up the services correctly. This file
defines the services required to run the application.

### License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE)
file for details.
