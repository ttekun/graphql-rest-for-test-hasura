# GraphQL-Rest-For-Test-Hasura

This monorepo contains two web applications: a REST API app and a GraphQL app.

## REST API App

The REST API app is developed using C#.

### Prerequisites

- .NET Core SDK

### Getting Started

1. Clone the repository.
2. Navigate to the `rest-api-app` directory.
3. Build the Docker image:
   ```
   docker build -t rest-api-app .
   ```
4. Run the Docker container:
   ```
   docker run -d -p 8080:80 rest-api-app
   ```

### API Endpoint

- GET /getCustomerInfo

### Example Request

```
GET /getCustomerInfo
```

### Example Response

```json
{
  "loginId": "123456",
  "name": "John Doe",
  "email": "johndoe@example.com",
  "address": "123 Main St"
}
```

## GraphQL App

The GraphQL app is developed using Go.

### Prerequisites

- Go

### Getting Started

1. Clone the repository.
2. Navigate to the `graphql-api-app` directory.
3. Build the Docker image:
   ```
   docker build -t graphql-api-app .
   ```
4. Run the Docker container:
   ```
   docker run -d -p 8081:80 graphql-api-app
   ```

### GraphQL Endpoint

- POST /graphql

### Example Query

```graphql
query {
  queryProducts {
    productId
    productName
    price
    availability
    introduction
  }
}
```

### Example Response

```json
{
  "data": {
    "queryProducts": [
      {
        "productId": "1",
        "productName": "Product 1",
        "price": 9.99,
        "availability": true,
        "introduction": "Lorem ipsum dolor sit amet."
      },
      {
        "productId": "2",
        "productName": "Product 2",
        "price": 19.99,
        "availability": false,
        "introduction": "Lorem ipsum dolor sit amet."
      }
    ]
  }
}
```

For more information, please refer to the individual README files in the `rest-api-app` and `graphql-api-app` directories.