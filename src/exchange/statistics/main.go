package statistics

import (
	types "exchange/transactionTypes"
	"time"
)

type OHLCV struct {
	ID       int64
	Interval int32
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Time     time.Time
	Ticker   string
}

type TickerStatistics struct {
	statistics OHLCV
	id         int64
	first      bool
}

var s map[string]*TickerStatistics

func init() {
	s = make(map[string]*TickerStatistics)
	s["RIM7"] = &TickerStatistics{
		id:    1,
		first: true,
	}
	s["SPFB.Si-6.17"] = &TickerStatistics{
		id:    1,
		first: true,
	}
}

func GenStatistics(line types.Order) {
	ticker := line.Ticker
	tickerStat := s[ticker]
	if tickerStat.first {
		tickerStat.statistics.ID = tickerStat.id
		tickerStat.statistics.Time = line.Time
		tickerStat.statistics.Interval = 1
		tickerStat.statistics.High = line.Price
		tickerStat.statistics.Low = line.Price
		tickerStat.statistics.Ticker = line.Ticker
		tickerStat.statistics.Open = line.Price
		tickerStat.first = false
	}
	tickerStat.statistics.Close = line.Price
	if tickerStat.statistics.High <= line.Price {
		tickerStat.statistics.High = line.Price
	}
	if tickerStat.statistics.Low >= line.Price {
		tickerStat.statistics.Low = line.Price
	}

}

func reload(ticker string) {
	s[ticker].id += 1
	s[ticker].first = true
}

func GetStatistics(ticker string) OHLCV {
	reload(ticker)
	return s[ticker].statistics
}
