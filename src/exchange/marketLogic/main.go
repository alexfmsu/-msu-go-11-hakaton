package marketLogic

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	stat "exchange/statistics"
	types "exchange/transactionTypes"
)

func Init() {
	go OrdersEmulate("../src/exchange/marketData/RIM7_170505_170505.txt")
	go OrdersEmulate("../src/exchange/marketData/SPFB.Si-6.17_170505_170505.txt")
}

func isNow(t time.Time) bool {
	now := time.Now()
	return t.Unix() == now.Unix()
}

func beforeNow(t time.Time) bool {
	now := time.Now()
	return t.Unix() < now.Unix()
}

func OrdersEmulate(filePath string) {
	log.Println("Open file: ", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	countLines := 0

	prevPrice := 0.0
	for scanner.Scan() {
		countLines++

		line := scanner.Text()
		lineData := strings.Split(line, ",")
		parseTicker := lineData[0]

		if countLines == 1 {
			tmp, err := strconv.ParseInt(lineData[3], 10, 64)
			if err != nil {
				log.Fatal("Not int")
			}
			StartTs = tmp
		}

		parseTime, err := getTime(lineData[3])
		if err != nil {
			continue
		}
		parsePrice, err := strconv.ParseFloat(lineData[4], 64)
		if err != nil {
			continue
		}
		parseAmount, err := strconv.Atoi(lineData[5])
		if err != nil {
			continue
		}

		var tp types.OrderType
		if countLines == 1 {
			tp = types.Buy
			prevPrice = parsePrice
		} else if parsePrice > prevPrice {
			tp = types.Buy
			prevPrice = parsePrice
		} else if parsePrice < prevPrice {
			tp = types.Sell
			prevPrice = parsePrice
		} else {
			tp = types.Buy
		}

		order := types.Order{
			Type:    tp,
			Ticker:  parseTicker,
			Time:    parseTime,
			Price:   parsePrice,
			Amount:  parseAmount,
			Instant: true,
		}

		// log.Printf("Order: %v\n", order)

		stat.GenStatistics(order)

		if beforeNow(order.Time) {
			log.Println("Order before now")
			continue
		}
		if !isNow(order.Time) {
			//log.Println("Sleep for ", time.Until(order.Time))
			time.Sleep(time.Until(order.Time))
		}
		if isNow(order.Time) {
			if order.Type == types.Buy {
				order.Type = types.Sell
			} else {
				order.Type = types.Buy
			}
			OrdersChannelIn <- &order
		} else {
			log.Println("Error in sleeping time")
		}
	}
}
