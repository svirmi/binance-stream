package exchanges

import (
	"log/slog"
	"strings"
	"time"

	spotws "github.com/linstohu/nexapi/binance/spot/websocketmarket"
	spotWsTypes "github.com/linstohu/nexapi/binance/spot/websocketmarket/types"

	"github.com/svirmi/binance-stream/models"
	"github.com/svirmi/binance-stream/utils"
)

func Pairs() []string {
	quoteAsset := "USDT"

	tPairs := utils.TradablePairs(quoteAsset)

	var symbols []string

	for _, tPair := range tPairs {
		symbols = append(symbols, tPair.Symbol)
	}

	pairs := utils.RemoveDuplicates(symbols)

	return pairs
}

func StartKlinestreams(stream *models.KLineStream, symbol string) {
	client, err := spotws.NewSpotMarketStreamClient(&spotws.SpotMarketStreamCfg{
		Debug:         false,
		BaseURL:       spotws.SpotMarketStreamBaseURL,
		Logger:        stream.Logger,
		AutoReconnect: true,
	})
	if err != nil {
		panic(err)
	}

	defer func() {
		stream.Logger.Info("closing connection for " + symbol)
		stream.Wg.Done()
	}()

	err = client.Open()
	if err != nil {
		panic(err)
	}

	topic, err := client.GetKlineTopic(&spotws.KlineTopicParam{
		Symbol:   strings.ToLower(symbol),
		Interval: "1s",
	})

	if err != nil {
		slog.Error("error client.GetKlineTopic()")
	}

	client.AddListener(topic, func(e any) {

		message := e.(*spotWsTypes.Kline)

		kline := &models.KlineTick{
			Symbol:   message.Symbol,
			Open:     message.Kline.OpenPrice,
			Close:    message.Kline.ClosePrice,
			High:     message.Kline.HighPrice,
			Low:      message.Kline.LowPrice,
			Volume:   message.Kline.BaseAssetVolume, // or QuotAssetVolume ???
			Interval: message.Kline.Interval,
		}

		// fmt.Println("volume", kline.Volume, " symbol", kline.Symbol)

		stream.KLineChan <- kline
	})

	s := append([]string{}, topic)

	client.Subscribe(s)

	stream.Logger.Info("subscribed to kline " + topic)

	<-stream.Closer

	client.UnSubscribe(s)
	client.Close()

	time.Sleep(800 * time.Millisecond)

	if !client.IsConnected() {
		stream.Logger.Info("disconnected without errors")
	} else {
		stream.Logger.Info("still connected after 800 ms delay")
	}
}
