package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"log/slog"
	"os"
	"runtime/debug"
	"svipp-server/internal/database"
	"svipp-server/internal/env"
	"sync"

	"svipp-server/internal/version"

	"github.com/lmittmann/tint"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

type config struct {
	baseURL  string
	httpPort int
	db       struct {
		dsn         string
		automigrate bool
	}
	jwt struct {
		secretKey string
	}
}

type application struct {
	config config
	db     *database.Queries
	logger *slog.Logger
	wg     sync.WaitGroup
}

func run(logger *slog.Logger) error {
	var cfg config

	cfg.baseURL = env.GetString("BASE_URL", "http://localhost:8080")
	cfg.httpPort = env.GetInt("HTTP_PORT", 8080)
	cfg.db.dsn = env.GetString("DB_DSN", "postgres://svipp@localhost:5432/svipp?sslmode=disable")
	cfg.db.automigrate = env.GetBool("DB_AUTOMIGRATE", true)
	cfg.jwt.secretKey = env.GetString("JWT_SECRET_KEY", "nVe2NeA2ByJDrDeDqOjGw0RBQS4WQkA53TY14DQl8/Q=")

	fmt.Printf("version: %s\n", version.Get())

	dbPool, err := pgxpool.New(context.Background(), cfg.db.dsn)
	if err != nil {
		return err
	}

	defer dbPool.Close()

	if cfg.db.automigrate {
		// Run Goose DB migrations
		goose.SetBaseFS(nil)

		if err = goose.SetDialect("postgres"); err != nil {
			return err
		}

		db := stdlib.OpenDBFromPool(dbPool)
		if err = goose.Up(db, "sql/migrations"); err != nil {
			return err
		}
		if err = db.Close(); err != nil {
			return err
		}
	}

	app := &application{
		config: cfg,
		db:     database.New(dbPool),
		logger: logger,
	}

	return app.serveHTTP()
}
