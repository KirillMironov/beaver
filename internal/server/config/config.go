package config

import "github.com/caarlos0/env/v8"

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	DataDir       string `env:"DATA_DIR,required"`
}

func Load() (config Config, _ error) {
	return config, env.ParseWithOptions(&config, env.Options{Prefix: "BEAVER_"})
}
