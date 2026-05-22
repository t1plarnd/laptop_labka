# laptop_labka

Lab work: REST API server in Go for managing laptops (CRUD). Backed by PostgreSQL.

## Requirements

- Docker + Docker Compose
- (optional) Go 1.25+ for local builds and tests

## Run

```
docker compose up --build
```

App on `:8080`, Postgres on `:5432`. Schema auto-created from `init.sql`.

Healthcheck:

```
curl http://localhost:8080/health
```

Stop and wipe data:

```
docker compose down -v
```

## Environment

| Var          | Example                                                  |
|--------------|----------------------------------------------------------|
| DATABASE_URL | `postgres://laptop:laptop@db:5432/laptops?sslmode=disable` |

## Tests

```
go test ./...
```

Tests use the in-memory store, no Postgres required.

## API

| Method | Path              | Description       |
|--------|-------------------|-------------------|
| POST   | /laptops          | create laptop     |
| GET    | /laptops          | list laptops      |
| GET    | /laptops/{id}     | get laptop by id  |
| PUT    | /laptops/{id}     | full update       |
| PATCH  | /laptops/{id}     | partial update    |
| DELETE | /laptops/{id}     | delete laptop     |

## Bruno collection

Requests for all endpoints are in `bruno/`. Open the folder in Bruno and pick the `local` environment.
