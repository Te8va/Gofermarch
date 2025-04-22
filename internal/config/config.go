package config

import (
	"flag"
	"log"
	"log/slog"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	PostgresUser         string `env:"POSTGRES_USER"         envDefault:"gophermart"`
	PostgresPassword     string `env:"POSTGRES_PASSWORD"     envDefault:"gophermart"`
	PostgresDB           string `env:"POSTGRES_DATABASE"     envDefault:"gophermart"`
	PostgresPort         int    `env:"POSTGRES_PORT"         envDefault:"5432"`
	MigrationsPath       string `env:"MIGRATIONS_PATH"       envDefault:"migrations"`
	LogFilePath          string `env:"LOG_FILE_PATH"         envDefault:"logfile.log"`
	JWTKey               string `env:"JWT_KEY"               envDefault:"supermegasecret"`
	DatabaseURI          string `env:"POSTGRES_CONN"         envDefault:"postgres://gophermart:gophermart@localhost:5432/gophermart?sslmode=disable"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://localhost:8081"`
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("Error loading .env file")
	}

	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error parsing environment variables: %v", err)
	}

	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "Service address")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "Database DSN")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", cfg.AccrualSystemAddress, "Accrual system address")
	flag.Parse()

	if err := env.Parse(cfg); err != nil {
		slog.Error("Failed to parse environment variables", "error", err)
	}

	return &cfg
}
