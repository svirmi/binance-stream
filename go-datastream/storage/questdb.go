package storage

import (
	"context"
	"log"
	"log/slog"

	qdb "github.com/questdb/go-questdb-client/v3"
	"github.com/svirmi/binance-stream/models"
)

type QuestDbConnection struct {
	Sender         qdb.LineSender
	TickChan       chan *models.Tick
	BidAskTickChan chan *models.BidAskTick
	KlineChan      chan *models.KlineTick
	Context        context.Context
	logger         *slog.Logger
}

func NewQuestDbConnection(url string, logger *slog.Logger) (*QuestDbConnection, error) {

	ctx := context.TODO()

	// Connect to QuestDB
	sender, err := qdb.LineSenderFromConf(ctx, url)
	if err != nil {
		log.Fatal(err)
	}

	return &QuestDbConnection{
		Sender:  sender,
		Context: ctx,
		logger:  logger,
	}, nil
}

func (questDb *QuestDbConnection) Close() {
	questDb.Sender.Close(questDb.Context)
	questDb.logger.Info("questDb connection closed")
}
