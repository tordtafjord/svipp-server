package sql

import (
	"embed"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed "migrations"
var migrationFiles embed.FS

func RunMigrations(dbPool *pgxpool.Pool) error {
	goose.SetBaseFS(nil)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(dbPool)
	defer db.Close()

	goose.SetBaseFS(migrationFiles)
	return goose.Up(db, "migrations")
}
