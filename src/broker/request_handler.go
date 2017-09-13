package main

import (
	bk "exchange/brokerInterface"
	"net/http"
	"strconv"
)

// отправить на биржу заявку на покупку или продажу тикера; отменить ранее отправленную заявку
// type: buy - купить, sell - продать, cancel - отмена (должен включать order_id)
// http://localhost:7777/request?token=rfBd56ti2SMtYvSgD5xAV0YU9&amount=1&price=1&ticker=Rim7&type=buy
func handleRequest(w http.ResponseWriter, r *http.Request, user *User) {
	response := Response{}
	response.UserS = user

	ticker := r.URL.Query().Get("ticker")
	typeReq := r.URL.Query().Get("type")
	ordIDstr := r.URL.Query().Get("order_id")
	if ordIDstr == "" {
		ordIDstr = "0"
	}
	ordID, errOi := strconv.Atoi(ordIDstr)
	amount, errAm := strconv.Atoi(r.URL.Query().Get("amount"))
	price, errPr := strconv.ParseFloat(r.URL.Query().Get("price"), 64)
	if errPr == nil && errAm == nil && errOi == nil {
		if ticker != "" && price > 0 && amount > 0 {
			switch typeReq {
			case "buy", "sell":
				if !requestExchange(user, ticker, amount, price, typeReq) {
					response.Error = "Произошла ошибка сервера"
				}
			case "cancel":
				if ordID == 0 {
					response.Error = "Введен некорректный ID заказа"
				} else if !cancelExchange(user, ordID) {
					response.Error = "Произошла ошибка сервера"
				}
			default:
				response.Error = "Указан не поддерживаемый тип операции"
			}
		} else {
			response.Error = "Указаны неправильные данные"
		}
	} else {
		response.Error = "Указаны неправильные данные или некоторые данные отсутствуют"
	}

	checkAndRespond(response, w)
}

func requestExchange(user *User, ticker string, amount int, price float64, typeReq string) bool {
	result, err := db.Exec(
		"INSERT INTO `positions` (`user_id`, `ticker`, `amount`, `price`, `type`, `status`) "+
			"VALUES (?, ?, ?, ?, ?, ?) ",
		user.ID, ticker, amount, price, typeReq, 0,
	)
	checkErr(err)
	lastID, err := result.LastInsertId()
	checkErr(err)

	request := bk.TransactionRequestType{
		BrokerID:        user.Broker.ID,
		Amount:          int32(amount),
		Ticker:          ticker,
		Price:           float32(price),
		ClientID:        int32(user.ID),
		OrderID:         int32(lastID),
		Action:          "add",
		TransactionType: typeReq,
	}
	if requestStream.Send(&request) != nil {
		return false
	}
	return true
}

func cancelExchange(user *User, ordID int) bool {
	request := bk.TransactionRequestType{
		BrokerID: user.Broker.ID,
		ClientID: int32(user.ID),
		OrderID:  int32(ordID),
		Action:   "remove",
	}
	if requestStream.Send(&request) != nil {
		return false
	}
	return true
}
