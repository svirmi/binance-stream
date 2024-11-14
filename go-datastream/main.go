package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/svirmi/binance-stream/config"
)

func main() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM)
	signal.Notify(interrupt, syscall.SIGINT)
	signal.Notify(interrupt, os.Interrupt)

	config.Load()

	log.Println("Hello from data stream!")

}
