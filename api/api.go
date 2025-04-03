package api

import (
	"exchange-tutorials/auth"
	"exchange-tutorials/match"
	"exchange-tutorials/types"
	"exchange-tutorials/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
	"log"
	"math/big"
	"net/http"
)

var upgrader = websocket.Upgrader{}
var limiter = rate.NewLimiter(10, 100) // 每秒 10 次，桶容量 100

func main() {
	auth := auth.NewAuthService("secret_key")
	fe := match.NewFuturesEngine()

	r := mux.NewRouter()
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		email, pass := r.FormValue("email"), r.FormValue("password")
		token, err := auth.Login(email, pass)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		w.Write([]byte(token))
	}).Methods("POST")

	r.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		token := r.Header.Get("Authorization")
		userID, err := auth.ValidateToken(token)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		order := &types.Order{
			ID:       utils.GenerateID(),
			UserID:   userID,
			Symbol:   "BTC/USD",
			Side:     "buy",
			Type:     "limit",
			Price:    big.NewFloat(50000),
			Quantity: big.NewFloat(1),
			Leverage: 10,
		}
		if err := fe.OpenPosition(order); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write([]byte("Order placed"))
	}).Methods("POST")

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()
		for trade := range fe.Engine.Trades() {
			conn.WriteJSON(trade)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
