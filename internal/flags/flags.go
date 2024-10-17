package flags

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type AppConfig struct {
	RUN_ADDRESS            string `env:"RUN_ADDRESS"`
	DATABASE_URI           string `env:"DATABASE_URI"`
	ACCRUAL_SYSTEM_ADDRESS string `env:"ACCRUAL_SYSTEM_ADDRESS"`

	ENV string `env:"ENV"`
}

func NewAppConfig() *AppConfig {
	cfg := &AppConfig{}
	readConfig(cfg)
	return cfg
}

func readConfig(cfg *AppConfig) {
	flag.StringVar(&cfg.RUN_ADDRESS, "a", "0.0.0.0:8080", "run address") // TODO: read from environment
	flag.StringVar(&cfg.DATABASE_URI, "d", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "database uri")
	flag.StringVar(&cfg.ACCRUAL_SYSTEM_ADDRESS, "r", "0.0.0.0:8080", "accrual system address")
	flag.StringVar(&cfg.ENV, "e", "PROD", "environment DEV - PROD default PROD")
	flag.Parse()
	env.Parse(cfg)
}
