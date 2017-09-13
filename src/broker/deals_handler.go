package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// посмотрить историю своих сделок
// http://127.0.0.1:7777/deals?token=rfBd56ti2SMtYvSgD5xAV0YU9
func handleDeals(w http.ResponseWriter, r *http.Request, user *User) {
	response := Response{}
	response.UserS = user

	var Data []Deal

	rows, err := db.Query("SELECT id, time, ticker, amount, price FROM orders_history WHERE user_id = ?", user.ID)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		deal := Deal{}
		err = rows.Scan(&deal.ID, &deal.Time, &deal.Ticker, &deal.Amount, &deal.Price)
		checkErr(err)
		Data = append(Data, deal)
	}

	response.Data = Data
	checkAndRespond(response, w)
}
