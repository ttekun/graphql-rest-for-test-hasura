# REST API App (Go)

This directory contains the Go implementation of the REST API application.

## Prerequisites

- Go (version 1.25.6 or later)
- Docker

## Getting Started

### Running Locally

1.  Navigate to the `rest-api-app` directory.
2.  Run the application:
    ```bash
    go run ./cmd/api-app
    ```
    The server will start on port `8080`.

### Building and Running with Docker

1.  Navigate to the `rest-api-app` directory.
2.  Build the Docker image:
    ```bash
    docker build -t rest-api-app .
    ```
3.  Run the Docker container:
    ```bash
    docker run -d -p 8080:8080 rest-api-app
    ```

## API Endpoints

-   `GET /Customer/getCustomerInfo`: Retrieves a list of customers.
-   `GET /healthz`: Health check endpoint. Returns `{"status": "ok"}`.
