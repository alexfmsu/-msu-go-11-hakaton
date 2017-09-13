package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// посмотреть последнюю истории торгов
// http://127.0.0.1:7777/bargain?token=rfBd56ti2SMtYvSgD5xAV0YU9&ticker=Rim7
func handleBargain(w http.ResponseWriter, r *http.Request, user *User) {
	response := Response{}
	response.UserS = user

	ticker := r.URL.Query().Get("ticker")
	if ticker != "" {
		if hist, ok := history[user.Broker.ID]; ok {
			response.Data = hist.Data[ticker]
		} else {
			response.Error = "Нет данных о торгах " + ticker
		}
	} else {
		response.Error = "Тикер отстуствует"
	}

	checkAndRespond(response, w)
}
