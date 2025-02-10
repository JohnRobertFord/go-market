package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
)

type Stringer interface {
	String() string
	Config
}

type Config struct {
	Bind        string `json:"bind" env:"RUN_ADDRESS"`
	Accrual     string `json:"" env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseURI string `json:"databaseURI" env:"DATABASE_URI"`
}

func (c *Config) String() string {
	return fmt.Sprintf("[Config] Host: %s, ACCRUAL_SYSTEM_ADDRESS:%v, DATABASE_URI:%s",
		c.Bind,
		c.Accrual,
		c.DatabaseURI)
}

func InitConfig() (*Config, error) {
	var cfg Config
	var envCfg Config

	flag.StringVar(&cfg.Bind, "a", ":8081", "адрес и порт запуска сервиса: env RUN_ADDRESS")
	flag.StringVar(&cfg.Accrual, "r", ":8082", "адрес системы расчёта начислений: env ACCRUAL_SYSTEM_ADDRESS")
	flag.StringVar(&cfg.DatabaseURI, "d", "host=localhost user=postgres_user password=postgres_password dbname=pos_db sslmode=disable", "адрес подключения к базе данных: env DATABASE_URI")

	flag.Parse()

	if err := env.Parse(&envCfg); err != nil {
		return nil, fmt.Errorf("cant load config: %e", err)
	}

	if os.Getenv("RUN_ADDRESS") != "" {
		cfg.Bind = envCfg.Bind
	}
	if os.Getenv("ACCRUAL_SYSTEM_ADDRESS") != "" {
		cfg.Accrual = envCfg.Accrual
	}
	if os.Getenv("DATABASE_URI") != "" {
		cfg.DatabaseURI = envCfg.DatabaseURI
	}

	return &cfg, nil
}
