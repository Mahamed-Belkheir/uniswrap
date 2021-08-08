package uniswrap

import "os"

type Config struct {
	Address string
}

func GetConfig() Config {
	return Config{
		Address: getEnv("ADDRESS", "127.0.0.1:8000"),
	}
}

func getEnv(name, defaultValue string) string {
	env := os.Getenv(name)
	if env == "" {
		return defaultValue
	}
	return env
}
