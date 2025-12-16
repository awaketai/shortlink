package config

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Port     string
	LogLevel string
}

func LoadConfig() (Config, error) {
	config := Config{
		Server: ServerConfig{
			Port:     "8080",
			LogLevel: "info",
		},
	}
	return config, nil
}
