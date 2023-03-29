package config

type Config struct {
	DataDir string
}

func Load() (Config, error) {
	return Config{
		DataDir: "/tmp/beaver",
	}, nil
}
