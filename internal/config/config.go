package config

import (
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/api/option"
	gmaps "googlemaps.github.io/maps"
	"log"
	"svipp-server/internal/auth"
	"svipp-server/internal/database"
	"svipp-server/internal/env"
	"svipp-server/internal/maps"
	"svipp-server/internal/sms"
	"svipp-server/sql"
)

type Services struct {
	DB          *database.Queries
	DBPool      *pgxpool.Pool
	MapsClient  *maps.MapsService
	SmsClient   *sms.TwilioClient
	FirebaseApp *firebase.App
	JwtService  *auth.JWTService
}

type Config struct {
	HTTPPort  int
	JwtSecret []byte
	IsProd    bool
}

func New() (*Config, *Services, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	services, err := InitializeServices(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize services: %w", err)
	}

	return cfg, services, nil

}

func loadConfig() (*Config, error) {
	cfg := &Config{}
	cfg.IsProd = env.GetBool("IS_PRODUCTION", false)
	log.Printf("Svipp server starting up in production-mode: %t", cfg.IsProd)
	err := env.LoadEnv(cfg.IsProd)
	if err != nil {
		return nil, err
	}
	cfg.HTTPPort = env.GetInt("PORT", 8080)
	cfg.JwtSecret = []byte(env.GetString("JWT_SECRET", "nVe2NeA2ByJDrDeDqOjGw0RBQS4WQkA53TY14DQl8/Q="))

	return cfg, nil
}

func InitializeServices(cfg *Config) (*Services, error) {
	services := &Services{}

	services.JwtService = auth.NewJWTService(&cfg.JwtSecret)

	// Initialize Maps client
	mapsClient, err := gmaps.NewClient(gmaps.WithAPIKey(env.GetString("GOOGLE_MAPS_API_KEY", "")))
	if err != nil {
		return nil, fmt.Errorf("failed to create google maps client: %w", err)
	}
	services.MapsClient = maps.NewMapsService(mapsClient)

	// Initialize SMS client
	services.SmsClient = sms.NewTwilioClient(
		env.GetString("TWILIO_ACCOUNT_SID", ""),
		env.GetString("TWILIO_AUTH_TOKEN", ""),
		env.GetString("TWILIO_MESSAGING_SERVICE_SID", ""),
	)

	// Initialize Firebase
	sa := option.WithCredentialsFile("transport-91700-firebase-adminsdk-b09hi-816702cf95.json")
	firebaseApp, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase: %w", err)
	}
	services.FirebaseApp = firebaseApp

	// Initialize Database
	dbUrl := env.GetString("DATABASE_URL", "postgres://svipp@localhost:5432/svipp?sslmode=disable")
	dbPool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	services.DBPool = dbPool
	services.DB = database.New(dbPool)

	if env.GetBool("DB_AUTOMIGRATE", true) {
		if err := sql.RunMigrations(dbPool); err != nil {
			return nil, fmt.Errorf("failed to run database migrations: %w", err)
		}
	}

	return services, nil
}
