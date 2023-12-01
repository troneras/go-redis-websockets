package config

import (
	"flag"
	"os"

	"github.com/troneras/gorews/data"
)

func DefaultConfig() *Config {
	return &Config{
		APIBindAddr: "0.0.0.0:8888",
		Domain:      "domain.com",
		Brand:       "brand",
		WebPath:     "",
		ApiPath:     "",
		Sha1Secret:  "secret",
	}
}

type Config struct {
	MessageChan chan *data.Message
	APIBindAddr string
	Domain      string
	Brand       string
	WebPath     string
	ApiPath     string
	Sha1Secret  string
}

var cfg = DefaultConfig()

func Configure() *Config {
	ReadEnv()
	return cfg
}

func ReadEnv() {
	flag.StringVar(&cfg.APIBindAddr, "api-bind-addr", getEnv("WS_API_BIND_ADDR", "0.0.0.0:8888"), "HTTP bind interface and port for API, e.g. 0.0.0.0:9027 or just :9027")
	flag.StringVar(&cfg.Domain, "api-domain", getEnv("WS_API_DOMAIN", "localhost"), "Domain where to send the requests for open/close events API, e.g. domain.com")
	flag.StringVar(&cfg.Brand, "api-brand", getEnv("WS_API_BRAND", "brand1"), "Brand name for API, e.g. brand")
	flag.StringVar(&cfg.WebPath, "api-web-path", getEnv("WS_API_WEB_PATH", ""), "Web path for API, e.g. /api/v1")
	flag.StringVar(&cfg.Sha1Secret, "ws-secret", getEnv("WS_SECRET", "secret"), "Sha1 secret for API, e.g. secret")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
