package config

import (
	"flag"
	"os"
)

func DefaultConfig() *Config {
	return &Config{
		UseTLS:  false,
		SSLCert: "",
		SSLKey:  "",
	}
}

type Config struct {
	UseTLS  bool
	SSLCert string
	SSLKey  string
}

var cfg = DefaultConfig()

func Configure() *Config {
	ReadEnv()
	return cfg
}

func ReadEnv() {
	flag.BoolVar(&cfg.UseTLS, "use-tls", getEnvBool("HTTP_USE_TLS", false), "Use TLS")
	flag.StringVar(&cfg.SSLCert, "ssl-cert", getEnv("HTTP_SSL_CERT", ""), "SSL certificate")
	flag.StringVar(&cfg.SSLKey, "ssl-key", getEnv("HTTP_SSL_KEY", ""), "SSL key")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		return value == "true"
	}
	return fallback
}
