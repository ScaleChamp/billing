package models

import "time"

type Environment struct {
	DatabaseUrl       string        `env:"DATABASE_URL,required"`
	DBConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"1m"`
	DBMaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" envDefault:"1"`
	DBMaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" envDefault:"1"`
}
