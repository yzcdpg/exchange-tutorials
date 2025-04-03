package main

import (
	"exchange-tutorials/match"
	"exchange-tutorials/types"
	"exchange-tutorials/utils"
	"github.com/gorilla/mux"
	"log"
	"math/big"
	"net/http"
)

func main() {
	fe := match.NewFuturesEngine()

	// REST API
	r := mux.NewRouter()
	r.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		// 处理订单请求 (省略解析逻辑)
		order := &types.Order{
			ID:       utils.GenerateID(),
			UserID:   "user1",
			Symbol:   "BTC/USD",
			Side:     "buy",
			Type:     "limit",
			Price:    big.NewFloat(50000),
			Quantity: big.NewFloat(1),
			Leverage: 10, // 合约交易
		}
		if err := fe.OpenPosition(order); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write([]byte("Order placed"))
	}).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}
