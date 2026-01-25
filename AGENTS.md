# Project Overview
Go-based REST and GraphQL API services orchestrated by Hasura GraphQL Engine and
PostgreSQL for local development and integration testing.

## Repository Structure
- graphql-api-app/: Go GraphQL service (entry in cmd/api-app, logic in internal/).
- rest-api-app/: Go REST service (entry in cmd/api-app, logic in internal/).
- hasura/: Hasura metadata and database migrations.
- docker-compose.yml: Local multi-service orchestration.
- README.md / AGENTS.md: Developer-facing runbooks and notes.
- docs/hasura-console-checklist.md: Console steps to validate Remote Schema and Action wiring.

## Build & Dev Commands
```bash
# Start/stop the full stack (from repo root)
docker-compose up -d
docker-compose down

# Logs
# All services
docker-compose logs -f
# Specific services
docker-compose logs -f hasura
docker-compose logs -f rest-api-app
docker-compose logs -f graphql-api-app

# REST API (from rest-api-app/)
go run ./cmd/api-app

docker build -t rest-api-app .

# GraphQL API (from graphql-api-app/)
go run ./cmd/api-app

docker build -t graphql-api-app .

# Hasura metadata/migrations (from hasura/)
hasura migrate apply --database-name default
hasura metadata apply
```

## Local Services (Ports & Endpoints)
- REST API: http://localhost:8080 (GET /Customer/getCustomerInfo, GET /healthz)
- GraphQL API: http://localhost:8081 (container listens on :8080, endpoint /graphql)
- Hasura: http://localhost:8082 (admin secret: `admin-secret` for local only; JWT secret is set in docker-compose.yml and must not be reused in non-local environments)
- PostgreSQL: localhost:5432 (user/password: postgres/postgres, db: postgresdb; local only)
- metadata-autofix helper: built from hasura/scripts and runs after Hasura is healthy to reconcile metadata.

## Hasura Automation
- `metadata-autofix` service (docker-compose) runs after the Hasura health check and triggers `replace_metadata` only if Actions are missing.
  - Image: `hasura/scripts/Dockerfile.metadata-autofix` (Python 3.10 + curl + PyYAML preinstalled)
  - Script: `hasura/scripts/metadata-autofix.sh`
  - Environment: `HASURA_ENDPOINT` (default `http://hasura:8080`), `HASURA_ADMIN_SECRET` (default `admin-secret`), optional `LOG_DIR` for log/output files
  - Flow: export current metadata → detect Actions from `actions.yaml` → if missing, run `replace_metadata` using local `hasura/metadata` → re-check

## Operational Checks
- Use docs/hasura-console-checklist.md to verify Remote Schema `graphql_api` and Action `getCustomerInfo` after `docker-compose up -d`.

## Code Style & Conventions
- Go version: 1.25.6 (module declarations).
- Package layout: cmd/ for entrypoints, internal/ for app packages.
- Naming: exported identifiers use PascalCase; unexported use camelCase.
- GoDoc comments: exported types/functions start with the symbol name.
- TODO: No explicit lint/format config found (e.g., .golangci.yml, .editorconfig).
- TODO: Confirm any project-specific style rules beyond standard Go conventions.

## Agent Guardrails (Security)
- Always:
  - Read files and inspect code.
  - Run a single-service test or lint command.
- Ask First:
  - Install or upgrade packages.
  - Delete or move files.
  - Run full builds or docker-compose up --build.
  - Push changes to git remotes.
- Never:
  - Edit .env files or commit secrets.
  - Hardcode passwords, tokens, or admin secrets.
  - Modify vendored dependencies.

## Definition of Done
- [ ] Code changes are minimal and scoped to the request.
- [ ] Local run command executed for the touched service.
- [ ] TODO: Confirm test command (no *_test.go found).
- [ ] TODO: Confirm lint/format command for Go sources.
- [ ] Documentation updated if behavior or commands changed.
