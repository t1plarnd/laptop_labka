# laptop_labka

Lab work: REST API server in Go for managing laptops (CRUD).

## Requirements

- Go 1.22+

## Run

```
go run .
```

Server listens on `:8080`. Healthcheck:

```
curl http://localhost:8080/health
```

## Tests

```
go test ./...
```

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
