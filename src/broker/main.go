package main

import (
	cfg "broker/config"
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

var sessions map[string]*User
var configData *cfg.Config
var history map[int32]History

func init() {
	sessions = map[string]*User{}
	configData = cfg.ParseConfigFromFile()
	history = map[int32]History{}

	var err error
	db, err = sql.Open("mysql", "root@tcp(localhost:3306)/broker?charset=utf8&interpolateParams=true")
	checkErr(err)
	db.SetMaxOpenConns(10)
	checkErr(db.Ping())

	initGRPC()
}

func main() {
	defer grcpConn.Close()
	// curl -X POST -d 'login=Vasily&password=123456' http://localhost:7777/
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			checkErr(r.ParseForm())

			var username = r.FormValue("login")
			var password = r.FormValue("password")

			var auth = Response{}
			if user, ok := getUser(username, password); ok {
				user.Token = randStringRunes(25)
				user.Broker = configData.Brokers[0]
				sessions[user.Token] = user
				auth.UserS = user
			} else {
				auth.Error = "Неверный логин/пароль"
			}

			checkAndRespond(auth, w)
		}
	})
	http.HandleFunc("/tickers", func(w http.ResponseWriter, r *http.Request) { // получить список тикеров
		checkAndRespond(Response{
			Data: configData.Tickers,
		}, w)
	})
	http.HandleFunc("/me", checkToken(handleMe))           // посмотреть свои позиции и баланс
	http.HandleFunc("/deals", checkToken(handleDeals))     // посмотрить историю своих сделок
	http.HandleFunc("/request", checkToken(handleRequest)) // отправить на биржу заявку на покупку или продажу тикера; отменить ранее отправленную заявку
	http.HandleFunc("/bargain", checkToken(handleBargain)) // посмотреть последнюю истории торгов
	http.ListenAndServe(configData.Host+":"+configData.Port, nil)
}

func getUser(name string, password string) (*User, bool) {
	var user = User{Name: name}
	var realPassword string
	err := db.QueryRow(
		"SELECT `id`, `password`, `balance` "+
			"FROM `clients` WHERE `login`=?", name,
	).Scan(&user.ID, &realPassword, &user.Balance)
	if err == sql.ErrNoRows {
		return nil, false
	}
	checkErr(err)

	if password == realPassword {
		return &user, true
	}
	return nil, false
}
