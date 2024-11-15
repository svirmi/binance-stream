package config

import (
	"os"
)

type Config struct {
	Env         string
	Logfile     string
	QuestDBAddr string
}

// https://github.com/MartinEllegard/tibber-harvester/blob/main/config/config.go ???

func Load() Config {
	env := os.Getenv("APP_ENV")

	if env == "" {
		return Config{
			Env:         "virtualmachine",
			Logfile:     "app.log",
			QuestDBAddr: "http::addr=localhost:9000;auto_flush_rows=100;auto_flush_interval=1000;",
		}
	}

	return Config{
		Env:         env, // docker
		Logfile:     "app/app.log",
		QuestDBAddr: "http::addr=questdb:9000;auto_flush_rows=100;auto_flush_interval=1000;",
	}
}
