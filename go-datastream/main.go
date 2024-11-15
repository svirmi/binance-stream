package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	qdb "github.com/questdb/go-questdb-client/v3"

	"github.com/svirmi/binance-stream/config"
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

	logger.Info("application started, environment:", slog.String("env", config.Env), "log file:", slog.String("logfile", ""))

	cnfgString := config.QuestDBAddr

	ctx := context.TODO()
	// Connect to QuestDB running locally.
	sender, err := qdb.LineSenderFromConf(ctx, cnfgString)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure to close the sender on exit to release resources.
	defer sender.Close(ctx)

	// Send a few ILP messages.
	ts = time.Now()

	err = sender.
		Table("trades_go").
		Symbol("pair", "USDGBP").
		Symbol("type", "change").
		Float64Column("traded_price", 0.83).
		Float64Column("limit_price", 0.84).
		Int64Column("qty", 100).
		At(ctx, ts)
	if err != nil {
		log.Fatal(err)
	}

	ts = time.Now()

	err = sender.
		Table("trades_go").
		Symbol("pair", "GBPJPY").
		Symbol("type", "sell").
		Float64Column("traded_price", 135.97).
		Float64Column("limit_price", 0.84).
		Int64Column("qty", 400).
		At(ctx, ts)
	if err != nil {
		log.Fatal(err)
	}

	ts = time.Now()

	err = sender.
		Table("trades_go").
		Symbol("pair", "BTCUSDT").
		Symbol("type", "change").
		Float64Column("traded_price", 10.83).
		Float64Column("limit_price", 509.84).
		Int64Column("qty", 1400).
		At(ctx, ts)
	if err != nil {
		log.Fatal(err)
	}

	// Make sure that the messages are sent over the network.
	err = sender.Flush(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(ts)

	logger.Info("Hello from data stream!")
}
