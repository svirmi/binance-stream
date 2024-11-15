package utils

import (
	"log"
	"log/slog"
	"os"

	"github.com/svirmi/binance-stream/config"
)

const (
	envDocker = "docker"
	envLocal  = "virtualmachine"
	envProd   = "prod"
)

func SetupLogger(cfg *config.Config) *slog.Logger {
	var logger *slog.Logger

	// for now all settingsa are the same
	switch cfg.Env {
	case envDocker:
		logfile, err := os.Create(cfg.Logfile)

		if err != nil {
			log.Fatal(err)
		}

		logger = slog.New(slog.NewTextHandler(logfile, nil))
	case envLocal:
		logfile, err := os.Create(cfg.Logfile)

		if err != nil {
			log.Fatal(err)
		}

		logger = slog.New(slog.NewTextHandler(logfile, nil))
	default:
		logfile, err := os.Create(cfg.Logfile)

		if err != nil {
			log.Fatal(err)
		}

		logger = slog.New(slog.NewTextHandler(logfile, nil))
	}

	return logger
}
