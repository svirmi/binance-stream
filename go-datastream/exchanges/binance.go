package exchanges

import (
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
