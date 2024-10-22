package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/svirmi/binance-stream/storage"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading config, %v", err)
	}
}

func main() {
	storage.NewMysqlClient()
}
