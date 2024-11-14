package config

import (
	"os"
)

// https://github.com/MartinEllegard/tibber-harvester/blob/main/config/config.go ???

func Load() string {
	env := os.Getenv("APP_ENV")

	if env == "" {
		return "http::addr=localhost:9000;auto_flush_rows=100;auto_flush_interval=1000;"
	}

	return "http::addr=questdb:9000;auto_flush_rows=100;auto_flush_interval=1000;"
}
