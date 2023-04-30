package config

import (
	"time"

	"github.com/caarlos0/env/v8"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	DataDir       string `env:"DATA_DIR,required"`

	JWT struct {
		Secret   string        `env:"JWT_SECRET,required"`
		TokenTTL time.Duration `env:"JWT_TOKEN_TTL" envDefault:"1h"`
	}
}

func Load() (config Config, _ error) {
	return config, env.ParseWithOptions(&config, env.Options{Prefix: "BEAVER_"})
}
