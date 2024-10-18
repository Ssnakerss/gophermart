package flags

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type AppConfig struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`

	ENV string `env:"ENV"`
}

func NewAppConfig() *AppConfig {
	cfg := &AppConfig{}
	readConfig(cfg)
	return cfg
}

func readConfig(cfg *AppConfig) {

	//приоритет настроек
	//1. Env
	//2. Command line
	//3. Deafault
	env.Parse(cfg)
	if cfg.RunAddress == "" {
		flag.StringVar(&cfg.RunAddress, "a", "0.0.0.0:8080", "run address") // TODO: read from environment
	}
	if cfg.DatabaseURI == "" {
		flag.StringVar(&cfg.DatabaseURI, "d", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "database uri")
	}
	if cfg.AccrualSystemAddress == "" {
		flag.StringVar(&cfg.AccrualSystemAddress, "r", "http://localhost:8081", "accrual system address")
	}
	flag.StringVar(&cfg.ENV, "e", "PROD", "environment DEV - PROD default PROD")
	flag.Parse()
}
