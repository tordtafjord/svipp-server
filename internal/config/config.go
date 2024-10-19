package config

import (
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/api/option"
	gmaps "googlemaps.github.io/maps"
	"log"
	"strconv"
	"svipp-server/internal/auth"
	"svipp-server/internal/database"
	"svipp-server/pkg/maps"
	"svipp-server/pkg/sms"
	"svipp-server/sql"
)

type Services struct {
	DB          *database.Queries
	DBPool      *pgxpool.Pool
	MapsClient  *maps.MapsService
	SmsClient   *sms.TwilioClient
	FirebaseApp *firebase.App
	AuthService *auth.Service
}

type Config struct {
	HTTPPort int
	IsProd   bool
	Domain   string
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
	cfg.IsProd = getEnvBool("IS_PRODUCTION", false)
	log.Printf("Svipp server starting up in production-mode: %t", cfg.IsProd)
	err := loadEnv(cfg.IsProd)
	if err != nil {
		return nil, err
	}
	cfg.HTTPPort = getEnvInt("PORT", 8080)

	if cfg.IsProd {
		cfg.Domain = "svipp.app"
	} else {
		cfg.Domain = "localhost:" + strconv.Itoa(cfg.HTTPPort)
	}

	return cfg, nil
}

func InitializeServices(cfg *Config) (*Services, error) {
	services := &Services{}

	// Initialize Maps client
	mapsClient, err := gmaps.NewClient(gmaps.WithAPIKey(getEnvString("GOOGLE_MAPS_API_KEY", "")))
	if err != nil {
		return nil, fmt.Errorf("failed to create google maps client: %w", err)
	}
	services.MapsClient = maps.NewMapsService(mapsClient)

	// Initialize SMS client
	services.SmsClient = sms.NewTwilioClient(
		getEnvString("TWILIO_ACCOUNT_SID", ""),
		getEnvString("TWILIO_AUTH_TOKEN", ""),
		getEnvString("TWILIO_MESSAGING_SERVICE_SID", ""),
		cfg.IsProd,
	)

	// Initialize Firebase
	sa := option.WithCredentialsFile("transport-91700-firebase-adminsdk-b09hi-816702cf95.json")
	firebaseApp, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase: %w", err)
	}
	services.FirebaseApp = firebaseApp

	// Initialize Database
	dbUrl := getEnvString("DATABASE_URL", "postgres://svipp@localhost:5432/svipp?sslmode=disable")
	dbPool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	services.DBPool = dbPool
	services.DB = database.New(dbPool)

	if getEnvBool("DB_AUTOMIGRATE", true) {
		if err := sql.RunMigrations(dbPool); err != nil {
			return nil, fmt.Errorf("failed to run database migrations: %w", err)
		}
	}

	services.AuthService = auth.NewAuthService(services.DB)

	return services, nil
}
