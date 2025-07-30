package db

import "fmt"

type PostgresConfig struct {
	User     string `env:"user"`
	Password string `env:"password"`
	DbName   string `env:"dbname"`
	Host     string `env:"host"`
	Port     int    `env:"port"`
}

func (cfg *PostgresConfig) GetConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DbName,
	)
}
