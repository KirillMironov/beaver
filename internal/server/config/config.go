package config

type Config struct {
	ServerAddress string
	DataDir       string
}

func Load() (Config, error) {
	return Config{
		ServerAddress: ":8080",
		DataDir:       "/tmp/beaver",
	}, nil
}
