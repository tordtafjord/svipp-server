# README



## Getting started

Before running the application you will need a working PostgreSQL installation and a valid DSN (data source name) for connecting to the database.


Note that this DSN must be in the format `user:pass@localhost:port/db` and **not** be prefixed with `postgres://`.

Make sure that you're in the root of the project directory, fetch the dependencies with `go mod tidy`, then run the application using `go run ./cmd/api`:

```
go mod tidy 
go run ./cmd/api
```

## HOT RELOAD RUN
```
air
```


### Connecting to remote database
ssh -L 8181:srv-captain--postgres-postgis:5432 tordrt@37.27.87.98 -p 4646

### Goose migrations
goose postgres "host=localhost port=5432 user=svipp dbname=svipp sslmode=disable" up
#### Goose migrations on remote:
goose -dir "sql/schema" postgres "postgres://transport:password@localhost:8181/transport?sslmode=disable" up


## Tailwind CSS Generation
```
npm run build-css
```

## DB code generation
```
sqlc generate
```