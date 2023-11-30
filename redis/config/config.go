package config

import (
	"flag"
	"os"
)

func DefaultConfig() *Config {
	return &Config{
		RedisAddr: "localhost:6379",
	}
}

type Config struct {
	RedisAddr string
}

var cfg = DefaultConfig()

func Configure() *Config {
	ReadEnv()
	return cfg
}

func ReadEnv() {
	flag.StringVar(&cfg.RedisAddr, "redis-addr", getEnv("REDIS_ADDR", "localhost:6379"), "Redis server address, e.g., localhost:6379")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
