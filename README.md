# laptop_labka

Lab work: REST API server in Go for managing laptops (CRUD). Backed by PostgreSQL.

## Requirements

- Docker + Docker Compose
- (optional) Go 1.25+ for local builds and tests

## Run

```
cp .env.example .env
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

Compose reads values from `.env` (copy from `.env.example`).

| Var               | Example   | Notes                            |
|-------------------|-----------|----------------------------------|
| POSTGRES_USER     | laptop    | DB user                          |
| POSTGRES_PASSWORD | laptop    | DB password                      |
| POSTGRES_DB       | laptops   | DB name                          |
| POSTGRES_PORT     | 5432      | Host port mapped to Postgres     |
| APP_PORT          | 8080      | Host port mapped to API          |

App container itself reads `DATABASE_URL`, which Compose builds from the values above.

## Tests

```
go test ./...
```

Handler tests use a mock store, no Postgres required.

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
