package config

type Config struct {
	DBPath     string
	ServerPort string
}

func Load() *Config {
	return &Config{
		DBPath:     "cats.db",
		ServerPort: "8080",
	}
}
