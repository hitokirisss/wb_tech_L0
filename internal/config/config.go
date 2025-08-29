package config

import (
	"fmt"
	"os"
)

type Config struct {
	HTTPAddr string

	PGUser string
	PGPass string
	PGHost string
	PGPort string
	PGDB   string

	CachePreload int
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func New() *Config {
	return &Config{
		HTTPAddr:     getenv("HTTP_ADDR", ":8081"),
		PGUser:       getenv("POSTGRES_USER", "orders_user"),
		PGPass:       getenv("POSTGRES_PASSWORD", "123123"),
		PGHost:       getenv("POSTGRES_HOST", "localhost"),
		PGPort:       getenv("POSTGRES_PORT", "55432"),
		PGDB:         getenv("POSTGRES_DB", "ordersdb"),
		CachePreload: 100, // предзагрузка кеша последних заказов
	}
}

func (c *Config) PGURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.PGUser, c.PGPass, c.PGHost, c.PGPort, c.PGDB)
}
