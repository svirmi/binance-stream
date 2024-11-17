package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/svirmi/binance-stream/config"
	"github.com/svirmi/binance-stream/models"
	"github.com/svirmi/binance-stream/storage"
	"github.com/svirmi/binance-stream/utils"
)

func main() {

	var ts time.Time

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM)
	signal.Notify(interrupt, syscall.SIGINT)
	signal.Notify(interrupt, os.Interrupt)

	config := config.Load()

	logger := utils.SetupLogger(&config)

	logger.Info("application started", slog.String("env", config.Env))

	questDb, err := storage.NewQuestDbConnection(config.QuestDBAddr, logger)

	if err != nil {
		logger.Error(err.Error())
		log.Fatal(err)
	}

	defer questDb.Close()

	// start goroutine to continuously read the KlineChan channel in RedisClient struct and publish data to redis
	go questDb.PublishKlineTick()

	klineStream := &models.KLineStream{
		KLineChan: make(chan *models.KlineTick),
		Wg:        &sync.WaitGroup{},
		Closer:    make(chan interface{}),
		Logger:    logger, // do I really need it?
	}

	klineStream.Logger.Info("message from klineStream logger")

	// Send a few ILP messages.
	ts = time.Now()

	err = questDb.Sender.
		Table("trades_go").
		Symbol("pair", "USDGBP").
		Symbol("type", "change").
		Float64Column("traded_price", 0.83).
		Float64Column("limit_price", 0.84).
		Int64Column("qty", 100).
		At(questDb.Context, ts)
	if err != nil {
		log.Fatal(err)
	}

	ts = time.Now()

	err = questDb.Sender.
		Table("trades_go").
		Symbol("pair", "GBPJPY").
		Symbol("type", "sell").
		Float64Column("traded_price", 135.97).
		Float64Column("limit_price", 0.84).
		Int64Column("qty", 400).
		At(questDb.Context, ts)
	if err != nil {
		log.Fatal(err)
	}

	ts = time.Now()

	err = questDb.Sender.
		Table("trades_go").
		Symbol("pair", "BTCUSDT").
		Symbol("type", "change").
		Float64Column("traded_price", 10.83).
		Float64Column("limit_price", 509.84).
		Int64Column("qty", 1400).
		At(questDb.Context, ts)
	if err != nil {
		log.Fatal(err)
	}

	// Make sure that the messages are sent over the network.
	err = questDb.Sender.Flush(questDb.Context)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(ts)

	logger.Info("Hello from data stream!")
}
