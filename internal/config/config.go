package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	Env      string `env:"env"`
	Server   Server
	Postgres Postgres
	Redis    Redis
}

type Server struct {
	Host string `env:"SERVER_HOST" env-default:"localhost"`
	Port string `env:"SERVER_PORT" env-default:"3000"`
}

type Postgres struct {
	User     string `env:"DB_USER" env-default:"postgres"`
	Password string `env:"DB_PASSWORD"`
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	Name     string `env:"DB_NAME" env-default:"postgres"`
	SSLMode  bool   `env:"SSL_MODE" env-default:"false"`
}

type Redis struct {
	Name     int    `env:"REDIS_NAME" env-default:"0"`
	Host     string `env:"REDIS_HOST" env-default:"localhost"`
	Port     string `env:"REDIS_PORT" env-default:"6379"`
	User     string `env:"REDIS_USER" env-default:""`
	Password string `env:"REDIS_PASSWORD" env-default:""`
}

func Load() *Config {
	var cfg Config

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	err = cleanenv.ReadEnv(&cfg)

	if err != nil {
		log.Fatalf("error while read config: %v", err)
	}

	return &cfg
}
