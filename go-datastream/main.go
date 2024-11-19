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
	"github.com/svirmi/binance-stream/exchanges"
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

	// start goroutine to continuously read the KlineChan channel in QuestDB client struct and publish data to redis
	go questDb.PublishKlineTick()

	klineStream := &models.KLineStream{
		KLineChan: make(chan *models.KlineTick),
		Wg:        &sync.WaitGroup{},
		Closer:    make(chan interface{}),
		Logger:    logger, // do I really need it here?
	}

	pairs := exchanges.Pairs()

	nPairs := len(pairs)

	if nPairs == 0 {
		logger.Error("no pairs to process")
		log.Fatal("no pairs to process")
	}

	logger.Info("pairs to process", slog.Int("pairs", nPairs))

	for _, symbol := range pairs { // https://github.com/binance/binance-spot-api-docs/blob/master/web-socket-streams.md
		klineStream.Wg.Add(1)
		go exchanges.StartKlinestreams(klineStream, symbol)
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

loop:

	for {
		<-interrupt
		break loop
	}

	close(klineStream.Closer)

	klineStream.Wg.Wait()

	close(klineStream.KLineChan)

	logger.Info("finishing program ...")

	// Make sure that the messages are sent over the network.
	err = questDb.Sender.Flush(questDb.Context)
	if err != nil {
		log.Fatal(err)
	}
}
