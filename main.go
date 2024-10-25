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
	dbClient, err := storage.NewMysqlClient()

	if err != nil {
		log.Fatal("error creating storage client", err)
	}

	defer func() {
		dbClient.Close()
	}()

	dbClient.CreateTableKLine()
}
