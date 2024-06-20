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

## Potential Features and Improvements

### Historical Data Fetching

- **Fetch Historical Events**: Fetch and index historical events from the
  contract.
- **Start Block Configuration**: Provide configuration to specify the start
  block for fetching historical events.

### Event Parsing

- **Additional Event Types**: Extend support to parse and handle other ERC-20
  events or events from other types of smart contracts.

### Database Enhancements

- **Batch Inserts**: Implement batch inserts for handling a high volume of
  events more efficiently.
- **Database Indexing**: Add indexing to the database tables to improve query
  performance.
- **Migration and Schema Version Management**: Use tools like `golang-migrate`
  to manage database schema changes and versioning.

### Logging and Monitoring

- **Structured Logging**: Implement structured logging with log levels for
  better traceability and debugging.
- **Monitoring and Alerting**: Integrate with monitoring tools like Prometheus
  and Grafana to visualize metrics and set up alerts.

### Testing and Quality Assurance

- **Unit Tests**: Increase the coverage of unit tests across all components.
- **Integration Tests**: Implement comprehensive integration tests to ensure all
  components work together as expected.

### API Endpoints

- Provide RESTful API endpoints for accessing event data and application status.

### Deployment

- **Docker and CI/CD**: Update CI/CD pipelines to push Docker images to a
  container registry.

### License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE)
file for details.
