package config

import (
	"log"
	"os"
)

func Load() {
	configFile := os.Getenv("CONFIG_FILE")

	if configFile == "" {
		log.Fatal("CONFIG_FILE variable is not set")
	}
}
