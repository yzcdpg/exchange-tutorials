package db

import (
	"database/sql"
	"exchange-tutorials/types"
	_ "github.com/lib/pq"
	"log"
)

func initDB() *sql.DB {
	db, err := sql.Open("postgres", "user=postgres password=secret dbname=exchange sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db.Exec(`
        CREATE TABLE trades (
            id TEXT PRIMARY KEY,
            symbol TEXT,
            price TEXT,
            quantity TEXT,
            buyer_id TEXT,
            seller_id TEXT,
            timestamp TIMESTAMP
        )
    `)
	return db
}

func SaveTrade(db *sql.DB, trade types.Trade) {
	_, err := db.Exec(
		"INSERT INTO trades (id, symbol, price, quantity, buyer_id, seller_id, timestamp) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		trade.ID, trade.Symbol, trade.Price.String(), trade.Quantity.String(), trade.BuyerID, trade.SellerID, trade.Timestamp,
	)
	if err != nil {
		log.Println("Failed to save trade:", err)
	}
}

func GetUserFromDB(email string) types.User {
	return types.User{}
}
