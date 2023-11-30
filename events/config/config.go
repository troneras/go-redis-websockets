package config

import (
	"flag"
	"os"
)

func DefaultConfig() *Config {
	return &Config{
		BasicUser: "admin",
		BasicPass: "admin",
	}
}

type Config struct {
	BasicUser string
	BasicPass string
}

var cfg = DefaultConfig()

func Configure() *Config {
	ReadEnv()
	return cfg
}

func ReadEnv() {
	flag.StringVar(&cfg.BasicUser, "basic-user", getEnv("BASIC_USER", "admin"), "Basic auth username")
	flag.StringVar(&cfg.BasicPass, "basic-pass", getEnv("BASIC_PASS", "celeron"), "Basic auth password")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
