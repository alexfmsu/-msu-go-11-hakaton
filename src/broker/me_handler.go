package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// посмотреть свои позиции и баланс
// http://127.0.0.1:7777/me?token=rfBd56ti2SMtYvSgD5xAV0YU9
func handleMe(w http.ResponseWriter, r *http.Request, user *User) {
	response := Response{}
	response.UserS = user

	var Data []Deal

	rows, err := db.Query("SELECT id, ticker, amount, price, status FROM positions WHERE user_id = ?", user.ID)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		deal := Deal{}
		var status int
		err = rows.Scan(&deal.ID, &deal.Ticker, &deal.Amount, &deal.Price, &status)
		deal.Status = getStatusName(status)
		checkErr(err)
		Data = append(Data, deal)
	}

	response.Data = Data
	checkAndRespond(response, w)
}
