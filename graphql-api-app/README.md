# GraphQL API App (Go)

This directory contains the Go implementation of the GraphQL API application.

## Prerequisites

- Go (version 1.25.6 or later)
- Docker

## Getting Started

### Running Locally

1. Navigate to the `graphql-api-app` directory.
2. Run the application:
   ```bash
   go run ./cmd/api-app
   ```
   The server will start on port `8080` and expose the GraphQL endpoint at `http://localhost:8080/graphql`.

### Building and Running with Docker

1. Navigate to the `graphql-api-app` directory.
2. Build the Docker image:
   ```bash
   docker build -t graphql-api-app .
   ```
3. Run the Docker container:
   ```bash
   docker run -d -p 8081:8080 graphql-api-app
   ```
   The container listens on port `8080`; the example maps it to host port `8081` to align with `docker-compose.yml`.

## GraphQL Endpoint

- `POST /graphql`

## Example Query

```graphql
query {
  queryProducts {
    id
    name
    price
    availability
    comment
  }
}
```

## Example Response

```json
{
  "data": {
    "queryProducts": [
      {
        "id": "1",
        "name": "Product 1",
        "price": 9.99,
        "availability": true,
        "comment": "Great product"
      },
      {
        "id": "2",
        "name": "Product 2",
        "price": 24.99,
        "availability": false,
        "comment": "Out of stock"
      }
    ]
  }
}
```
