package main

import (
	"encoding/json"
	// "log"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const server_url string = "http://127.0.0.1:7777"

var user User

var tickers []string
var ticker string

var msg = map[string]string{
	"msg":  "",
	"back": "",
}

func root(w http.ResponseWriter, r *http.Request) {} // ?

func Login(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("src/client/login.gtpl")
	t.Execute(w, nil)
}

func SellFormHandler(w http.ResponseWriter, r *http.Request) {
	if user.Token == "" {
		r.ParseForm()

		// log.Println("Sell:\n")

		login := r.Form["login"][0]
		password := r.Form["password"][0]

		err := get_token(login, password)

		if err != nil || user.Token == "" {
			PrintResultPage(w, r, "Неверно введены логин/пароль", "/")

			return
		}

		get_tickers(w, r)
	}

	t, _ := template.ParseFiles("src/client/sell_form.gtpl")
	t.Execute(w, tickers)
}

func BuyFormHandler(w http.ResponseWriter, r *http.Request) {
	if user.Token == "" {
		r.ParseForm()

		// log.Println("Buy:\n")

		login := r.Form["login"][0]
		password := r.Form["password"][0]

		get_token(login, password)
	}

	t, _ := template.ParseFiles("src/client/buy_form.gtpl")
	t.Execute(w, tickers)
}

func request(w http.ResponseWriter, r *http.Request, option string, gtpl string) {
	r.ParseForm()

	amount := r.Form["amount"][0]
	price := r.Form["price"][0]
	ticker := r.Form["ticker"][0]

	form := url.Values{}

	var err error

	resp, err := http.Post(server_url+"/request?token="+user.Token+"&amount="+amount+"&price="+price+"&ticker="+ticker+"&type="+option, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))

	if err != nil {
		log.Println("Error:", err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error:", err)

		PrintResultPage(w, r, "Bad request", "src/client/sell_form")

		return
	}

	var res Response

	err = json.Unmarshal(body, &res)

	if err != nil {
		// log.Println("Error:", err)
		PrintResultPage(w, r, "Bad request", "src/client/sell_form")

		return
	}

	var deals []Deal

	v, err := json.Marshal(&res.Data)

	if err != nil {
		log.Println("Error:", err)
	}

	err = json.Unmarshal(v, &deals)

	if err != nil {
		log.Println("Error:", err)
	}

	t, _ := template.ParseFiles(gtpl)

	t.Execute(w, nil)

	// log.Printf(option+":\n%#v\n", res)

	// log.Println("_amount:", amount)
	// log.Println("_price:", price)
	// log.Println("_ticker:", ticker)
}

func sell(w http.ResponseWriter, r *http.Request) {
	request(w, r, "sell", "src/client/sold.gtpl")
}

func buy(w http.ResponseWriter, r *http.Request) {
	request(w, r, "buy", "src/client/bought.gtpl")
}

func me(w http.ResponseWriter, r *http.Request) {
	form := url.Values{}

	var err error

	resp, err := http.Post(server_url+"/me?token="+user.Token, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))

	if err != nil {
		log.Println("Error:", err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error:", err)
	}

	var res Response

	err = json.Unmarshal(body, &res)

	if err != nil {
		log.Println("Error:", err)
	}

	if res.Error != "" {
		log.Println("Error:", res.Error)

		return
	}

	user = *res.UserS

	var deals []Deal

	v, err := json.Marshal(&res.Data)
	if err != nil {
		log.Println("err:2", err)
	}

	err = json.Unmarshal(v, &deals)

	if err != nil {
		log.Println("Error:", err)
	}

	t, _ := template.ParseFiles("src/client/me.gtpl")
	t.Execute(w, deals)
}

func get_tickers(w http.ResponseWriter, r *http.Request) {
	form := url.Values{}

	var err error

	resp, err := http.Post(server_url+"/tickers?token="+user.Token, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))

	if err != nil {
		log.Println("err:", err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("err:", err)
	}

	var res Response

	err = json.Unmarshal(body, &res)

	if err != nil {
		log.Println("Error:", err)
	}

	if res.Error != "" {
		log.Println("Error:", res.Error)

		return
	}

	v, err := json.Marshal(&res.Data)

	if err != nil {
		log.Println("Error:", err)
	}

	var _t []Ticker

	err = json.Unmarshal(v, &_t)

	if err != nil {
		log.Println("Error2:", err)
	}

	// log.Printf("Tickers:\n%#v\n", res)
	// log.Printf("_tickers:\n%#v\n", _t)

	// log.Printf("%#v\n", res)

	tickers = make([]string, 0)

	for _, tick := range _t {
		tickers = append(tickers, tick.Name)
	}

	ticker = tickers[0]
}

func bargain(w http.ResponseWriter, r *http.Request) {
	// log.Println("ticker:", ticker)
	// log.Println("token:", user.Token)

	form := url.Values{}

	var err error

	resp, err := http.Post(server_url+"/bargain?token="+user.Token+"&ticker="+ticker, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))

	if err != nil {
		log.Println("err:", err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("err:", err)
	}

	var res Response

	err = json.Unmarshal(body, &res)

	if err != nil {
		log.Println("Error:", err)
	}

	if res.Error != "" {
		log.Println("Error:", res.Error)

		return
	}

	user = *res.UserS

	v, err := json.Marshal(&res.Data)

	if err != nil {
		log.Println("Error:", err)
	}

	var hist []OHLCV

	err = json.Unmarshal(v, &hist)

	if err != nil {
		log.Println("Error2:", err)
	}

	log.Printf("Bargain:\n%#v\n", res)
	log.Printf("_hist:\n%#v\n", hist)

	t, _ := template.ParseFiles("src/client/hist.gtpl")
	t.Execute(w, hist)
}

func deals(w http.ResponseWriter, r *http.Request) {
	// log.Println("ticker:", ticker)
	log.Println("token:", user.Token)

	form := url.Values{}

	var err error

	resp, err := http.Post(server_url+"/deals?token="+user.Token, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))

	if err != nil {
		log.Println("err:", err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("err:", err)
	}

	var res Response

	err = json.Unmarshal(body, &res)

	if err != nil {
		log.Println("Error:", err)
	}

	if res.Error != "" {
		log.Println("Error:", res.Error)

		return
	}

	user = *res.UserS

	v, err := json.Marshal(&res.Data)

	if err != nil {
		log.Println("Error:", err)
	}

	var deals []Deal

	err = json.Unmarshal(v, &deals)

	if err != nil {
		log.Println("Error2:", err)
	}

	log.Printf("Deals:\n%#v\n", res)
	log.Printf("_deals:\n%#v\n", deals)

	t, _ := template.ParseFiles("src/client/deals.gtpl")
	t.Execute(w, deals)
}

func cancel(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	log.Println(r.Form["cancel"][0])
	form := url.Values{}

	var err error

	resp, err := http.Post(server_url+"/request?token="+user.Token+r.Form["cancel"][0]+"&type=cancel", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))

	if err != nil {
		log.Println("Error:", err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error:", err)
	}

	var res Response

	err = json.Unmarshal(body, &res)

	if err != nil {
		log.Println("Error:", err)
	}

	if res.Error != "" {
		PrintResultPage(w, r, "Error: "+res.Error, "/me")

		return
	}

	PrintResultPage(w, r, "Ваша позиция удалена", "/me")

	// log.Println("Заказ удален")

	// log.Printf("Cancel:\n%#v\n", res)

	t, _ := template.ParseFiles("result_page.gtpl")
	t.Execute(w, msg)
}

func main() {
	// http.HandleFunc("/", root)
	// http.HandleFunc("/login", Login)
	http.HandleFunc("/", Login)

	http.HandleFunc("/sell_form", SellFormHandler)
	http.HandleFunc("/buy_form", BuyFormHandler)

	http.HandleFunc("/sell", sell)
	http.HandleFunc("/buy", buy)
	http.HandleFunc("/me", me)

	http.HandleFunc("/bargain", bargain)
	http.HandleFunc("/deals", deals)

	http.HandleFunc("/cancel", cancel)

	err := http.ListenAndServe(":9090", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func get_token(login string, password string) error {
	form := url.Values{}
	form.Add("login", login)
	form.Add("password", password)

	var err error

	resp, err := http.Post(server_url, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))

	if err != nil {
		log.Println("err:", err)

		return err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("err:", err)

		return err
	}

	var res Response

	err = json.Unmarshal(body, &res)

	if err != nil {
		log.Println("err:", err)
	}

	if res.Error != "" {
		log.Println("Error:", res.Error)

		return nil
	}

	user = *res.UserS

	// log.Printf("Get Token:\n%#v\n", res)

	return nil
}

// --------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
type Response struct {
	Data  interface{} `json:"data"`
	UserS *User       `json:"user"`
	Error string      `json:"error"`
}

type Broker struct {
	ID   int32  `yaml:"id" json:"id"`
	Name string `yaml:"name" json:"name"`
}

type User struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
	Token   string  `json:"token"`
	broker  Broker  `json:"broker"`
}

type Deal struct {
	ID     int     `json:"id"`
	Time   int     `json:"time"`
	UserID int     `json:"-"`
	Ticker string  `json:"ticker"`
	Amount int     `json:"amount"`
	Price  float64 `json:"price"`
	Status string  `json:"status"`
}

type History struct {
	BrokerID int32
	Data     map[string][]OHLCV
}

type OHLCV struct {
	ID       int64   `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
	Time     int32   `protobuf:"varint,2,opt,name=Time" json:"Time,omitempty"`
	Interval int32   `protobuf:"varint,3,opt,name=Interval" json:"Interval,omitempty"`
	Open     float32 `protobuf:"fixed32,4,opt,name=Open" json:"Open,omitempty"`
	High     float32 `protobuf:"fixed32,5,opt,name=High" json:"High,omitempty"`
	Low      float32 `protobuf:"fixed32,6,opt,name=Low" json:"Low,omitempty"`
	Close    float32 `protobuf:"fixed32,7,opt,name=Close" json:"Close,omitempty"`
	Ticker   string  `protobuf:"bytes,8,opt,name=Ticker" json:"Ticker,omitempty"`
}

type Ticker struct {
	ID   int32  `yaml:"id" json:"id"`
	Name string `yaml:"name" json:"name"`
}

// --------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
func PrintResultPage(w http.ResponseWriter, r *http.Request, _msg string, _back string) {
	msg["msg"] = _msg
	msg["back"] = _back

	t, _ := template.ParseFiles("src/client/result_page.gtpl")
	t.Execute(w, msg)
}
