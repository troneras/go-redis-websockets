package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

func DefaultConfig() *Config {
	return &Config{
		LogLevel: logrus.InfoLevel,
		LogFile:  os.Stdout,
	}
}

type Config struct {
	LogLevel logrus.Level
	LogFile  *os.File
}

var cfg = DefaultConfig()

func Configure() *Config {
	ReadEnv()
	return cfg
}

func ReadEnv() {
	log_level, err := logrus.ParseLevel(getEnv("LOG_LEVEL", "info"))
	if err != nil {
		log_level = logrus.InfoLevel
	}
	cfg.LogLevel = log_level

	log_output := getEnv("LOG_FILE", "stdout")
	switch log_output {
	case "stdout":
		cfg.LogFile = os.Stdout
	case "stderr":
		cfg.LogFile = os.Stderr
	default:
		f, err := os.OpenFile(log_output, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			f = os.Stderr
		}
		cfg.LogFile = f
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
