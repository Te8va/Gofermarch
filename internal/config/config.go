package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
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
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error parsing environment variables: %v", err)
	}

	runAddrFlag := flag.String("a", "", "Address and port to run the service")
	dbURIFlag := flag.String("d", "", "Database connection URI")
	accrualAddrFlag := flag.String("r", "", "Accrual system address")

	flag.Parse()

	if *runAddrFlag != "" {
		cfg.RunAddress = *runAddrFlag
	}
	if *dbURIFlag != "" {
		cfg.DatabaseURI = *dbURIFlag
	}
	if *accrualAddrFlag != "" {
		cfg.AccrualSystemAddress = *accrualAddrFlag
	}

	return &cfg
}
