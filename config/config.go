package config

import (
	"fmt"
	"net/url"
	"time"
)

// Config is the struct that holds the configuration of the application
type Config struct {
	Redis    Redis
	Postgres Postgres
	InfluxDB InfluxDB
}

// Redis is the struct that holds the configuration of the Redis connection
type Redis struct {
	Addr string `envconfig:"REDIS_ADDR" default:"localhost:6379"`
	Pass string `envconfig:"REDIS_PASS" default:""`
	DB   int    `envconfig:"REDIS_DB" default:"0"`
}

type Postgres struct {
	Host             string        `envconfig:"POSTGRES_HOST" default:"localhost"`
	Port             string        `envconfig:"POSTGRES_PORT" default:"5432"`
	User             string        `envconfig:"POSTGRES_USER" default:""`
	Pass             string        `envconfig:"POSTGRES_PASS" default:""`
	DB               string        `envconfig:"POSTGRES_DB" default:""`
	SSLMode          string        `envconfig:"POSTGRES_SSL_MODE" default:"disable"`
	StatementTimeout time.Duration `envconfig:"POSTGRES_STATEMENT_TIMEOUT" default:"500s"`
}

func (p Postgres) URL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&timezone=UTC&statement_timeout=%d",
		p.User,
		url.QueryEscape(p.Pass),
		p.Host,
		p.Port,
		p.DB,
		p.SSLMode,
		p.StatementTimeout.Milliseconds())
}

type InfluxDB struct {
	URL   string `envconfig:"INFLUXDB_URL" default:"http://localhost:8086"`
	Token string `envconfig:"INFLUXDB_TOKEN" default:""`
}
