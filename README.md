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
