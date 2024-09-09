package config

import (
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"google.golang.org/api/option"
	"svipp-server/internal/database"
	"svipp-server/internal/env"
)

type Config struct {
	BaseURL  string
	HTTPPort int
	DB       struct {
		DBQ         *database.Queries
		URL         string
		Automigrate bool
		DBPool      *pgxpool.Pool
	}
	JWT struct {
		SecretKey []byte
	}
	Maps struct {
		APIKey string
	}
	Pricing struct {
		CostPerMin    float64
		PricingFactor float64
	}
	FirebaseApp *firebase.App
	IsProd      bool
}

func TestConfig() (*Config, error) {
	cfg := &Config{}
	cfg.IsProd = env.GetBool("IS_PRODUCTION", false)
	err := env.LoadEnv(cfg.IsProd)
	if err != nil {
		return nil, err
	}

	cfg.BaseURL = env.GetString("BASE_URL", "http://localhost:8080")
	cfg.HTTPPort = env.GetInt("PORT", 8080)
	cfg.DB.URL = env.GetString("TEST_DATABASE_URL", "postgres://svipp@localhost:5432/svipp-test?sslmode=disable")
	cfg.DB.Automigrate = env.GetBool("DB_AUTOMIGRATE", true)
	cfg.JWT.SecretKey = []byte(env.GetString("JWT_SECRET", "nVe2NeA2ByJDrDeDqOjGw0RBQS4WQkA53TY14DQl8/Q="))
	cfg.Maps.APIKey = env.GetString("GOOGLE_MAPS_API_KEY", "")
	cfg.Pricing.PricingFactor = env.GetFloat("REVENUE_FACTOR", 1.2)
	cfg.Pricing.CostPerMin = env.GetFloat("DRIVER_COST_PER_MIN", 5.0)

	// Initialize Firebase Admin SDK
	sa := option.WithCredentialsFile("transport-91700-firebase-adminsdk-b09hi-816702cf95.json")
	cfg.FirebaseApp, err = firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		return nil, err
	}

	err = cfg.initDB()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func New() (*Config, error) {
	cfg := &Config{}
	cfg.IsProd = env.GetBool("IS_PRODUCTION", false)
	err := env.LoadEnv(cfg.IsProd)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Svipp server running in production: %t\n\n", cfg.IsProd)

	cfg.BaseURL = env.GetString("BASE_URL", "http://localhost:8080")
	cfg.HTTPPort = env.GetInt("PORT", 8080)
	cfg.DB.URL = env.GetString("DATABASE_URL", "postgres://svipp@localhost:5432/svipp?sslmode=disable")
	cfg.DB.Automigrate = env.GetBool("DB_AUTOMIGRATE", true)
	cfg.JWT.SecretKey = []byte(env.GetString("JWT_SECRET", "nVe2NeA2ByJDrDeDqOjGw0RBQS4WQkA53TY14DQl8/Q="))
	cfg.Maps.APIKey = env.GetString("GOOGLE_MAPS_API_KEY", "")
	cfg.Pricing.PricingFactor = env.GetFloat("REVENUE_FACTOR", 1.2)
	cfg.Pricing.CostPerMin = env.GetFloat("DRIVER_COST_PER_MIN", 5.0)

	// Initialize Firebase Admin SDK
	sa := option.WithCredentialsFile("transport-91700-firebase-adminsdk-b09hi-816702cf95.json")
	cfg.FirebaseApp, err = firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		return nil, err
	}

	err = cfg.initDB()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) initDB() error {
	var err error
	c.DB.DBPool, err = pgxpool.New(context.Background(), c.DB.URL)
	if err != nil {
		return err
	}

	c.DB.DBQ = database.New(c.DB.DBPool)

	if c.DB.Automigrate {
		if err := c.runMigrations(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) runMigrations() error {
	goose.SetBaseFS(nil)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(c.DB.DBPool)
	defer db.Close()

	return goose.Up(db, "sql/migrations")
}
