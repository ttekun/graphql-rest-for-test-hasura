# Hasura Console Checklist (Remote Schema / Action)

## Prerequisites

- [ ] Run `docker-compose up -d`
- [ ] Open Hasura Console: http://localhost:8082
- [ ] If prompted for the admin secret, enter `admin-secret`

## 1) Basic service health checks

- [ ] REST API is reachable (example): http://localhost:8080/health
- [ ] GraphQL API endpoint is reachable (example): http://localhost:8081/graphql
- [ ] Hasura health endpoint is reachable (example): http://localhost:8082/healthz

## 2) Verify Remote Schema is registered (graphql-api-app)

### In Hasura Console

- [ ] Go to **Remote Schemas**
- [ ] Confirm `graphql_api` exists
- [ ] Open `graphql_api` details and confirm the URL is:
  - [ ] `http://graphql-api-app:8080/graphql`

### Runtime check (via Hasura GraphiQL)

- [ ] Go to **API** (GraphiQL)
- [ ] Run a query that belongs to the remote schema (use autocomplete to find available fields)

Example (if available in the schema):

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

- [ ] The query succeeds and returns data (no remote schema / upstream errors)

## 3) Verify Action is registered (rest-api-app)

### In Hasura Console

- [ ] Go to **Actions**
- [ ] Confirm action `getCustomerInfo` exists
- [ ] Open `getCustomerInfo` details and confirm the handler is:
  - [ ] `http://rest-api-app:8080/Customer/getCustomerInfo`

> Note: `getCustomerInfo` is registered as a **query** action (not a mutation).

### Runtime check (via Hasura GraphiQL)

- [ ] Go to **API** (GraphiQL)
- [ ] Use autocomplete to locate `getCustomerInfo` under the `Query` root
- [ ] Execute it and confirm it returns the expected fields

Example (adjust root type based on where it appears):

```graphql
query {
  getCustomerInfo {
    loginId
    name
    email
    address
  }
}
```

- [ ] The action call succeeds and returns data (no handler / upstream errors)

## 4) Quick troubleshooting

- [ ] If `graphql_api` remote schema is unhealthy, check `graphql-api-app` container logs and that `/graphql` responds inside the Docker network
- [ ] If `getCustomerInfo` action fails, check `rest-api-app` container is running (note: Hasura does not `depends_on` `rest-api-app` in `docker-compose.yml`)
- [ ] If calling Hasura APIs directly (not via Console), include header:
  - [ ] `x-hasura-admin-secret: admin-secret`
