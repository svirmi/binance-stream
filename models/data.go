package models

import (
	"log/slog"
	"sync"
	"time"
)

type KlineTick struct {
	Symbol      string `json:"symbol"`
	Open        string `json:"open"`
	Close       string `json:"close"`
	High        string `json:"high"`
	Low         string `json:"low"`
	Volume      string `json:"volume"` // TODO : check why volume is 0
	Interval    string `json:"interval"`
	TradeCount  int    `json:"tradeCount"`
	TakerVolume string `json:"takerVolume"`
	TakerAmount string `json:"takerAmount"`
	Amount      string `json:"amount"`
	OpenTime    int64  `json:"openTime"`
	CloseTime   int64  `json:"closeTime"`
}

type KLineStream struct {
	KLineChan chan *KlineTick
	Wg        *sync.WaitGroup
	Closer    chan interface{}
	Logger    *slog.Logger
}

type Tick struct {
	StandardSymbol string    `json:"StandardSymbol"`
	ExchangeSymbol string    `json:"ExchangeSymbol"`
	Price          string    `json:"Price"`
	Exchange       string    `json:"Exchange"`
	Time           time.Time `json:"Time"`
}

type BidAskTick struct {
	StandardSymbol string    `json:"StandardSymbol"`
	ExchangeSymbol string    `json:"ExchangeSymbol"`
	BidPrice       string    `json:"Bid"`
	AskPrice       string    `json:"Ask"`
	Exchange       string    `json:"Exchange"`
	Time           time.Time `json:"Time"`
}

type Stream struct {
	TickChan       chan *Tick
	BidAskTickChan chan *BidAskTick
	Wg             *sync.WaitGroup
	Closer         chan interface{}
	Logger         *slog.Logger
}

type SymbolMap map[string]string
