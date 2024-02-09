package config

// Config is the struct that holds the configuration of the application
type Config struct {
	Redis Redis
}

// Redis is the struct that holds the configuration of the Redis connection
type Redis struct {
	Addr     string `envconfig:"REDIS_ADDR" default:"localhost:6379"`
	Password string `envconfig:"REDIS_PASS" default:""`
	DB       int    `envconfig:"REDIS_DB" default:"0"`
}
