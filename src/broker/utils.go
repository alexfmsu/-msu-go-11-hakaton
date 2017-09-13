package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type handlerFuncUser func(w http.ResponseWriter, r *http.Request, user *User)
type handlerFunc func(w http.ResponseWriter, r *http.Request)

func checkToken(function handlerFuncUser) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if user, ok := sessions[token]; ok {
			function(w, r, user)
		} else {
			response := Response{}
			response.Error = "Токен отстуствует"
			data, err := json.Marshal(response)
			checkErr(err)
			w.Write(data)
		}
	}
}

func getStatusName(st int) string {
	switch st {
	case 0:
		return "статус не определён"
	case 1:
		return "запрос на транзакцию добавлен"
	case 2:
		return "транзакция совершена полностью"
	case 3:
		return "транзакция совершена частично"
	case 4:
		return "запрос на транзакцию удалён"
	case -1:
		return "ошибка, несуществующий Ticker"
	case -2:
		return "ошибка, неправильный тип запроса"
	case -3:
		return "ошибка, неправильное действие с запросом"
	default:
		return ""
	}
}

func checkAndRespond(response Response, w http.ResponseWriter) {
	data, err := json.Marshal(response)
	if err != nil {
		respondError(w, 0)
	} else {
		w.Write(data)
	}
}

func respondError(w http.ResponseWriter, repeat int) {
	resp := Response{Error: "Произошла ошибка сервера. Попробуйте позже"}
	data, err := json.Marshal(resp)
	if err != nil {
		if repeat < 2 {
			time.Sleep(time.Second)
			respondError(w, repeat+1)
		} else {
			w.Write([]byte("Да чтож такое! Не твой день сегодня!"))
		}
	} else {
		w.Write(data)
	}
}

func checkErr(err error) {
	if err != nil {
		// log.Fatalln(err)
		panic(err)
	}
}

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
