package types

import (
	"math/big"
	"sync"
	"time"
)

// 用户
type User struct {
	ID        string
	Email     string
	Password  string // 加密存储
	Balance   map[string]*big.Float
	Margin    map[string]*big.Float
	Positions map[string]*Position
	Mode      string // "isolated"（逐仓）或 "cross"（全仓）
}

// 订单
type Order struct {
	ID        string
	UserID    string
	Symbol    string
	Side      string
	Type      string
	Price     *big.Float
	Quantity  *big.Float
	Filled    *big.Float
	Status    string
	Leverage  int
	Timestamp time.Time
}

// 持仓
type Position struct {
	Symbol       string
	Side         string
	Quantity     *big.Float
	EntryPrice   *big.Float
	Leverage     int
	Margin       *big.Float
	UnrealizedPL *big.Float // 未实现盈亏
}

// 订单簿
type OrderBook struct {
	Symbol string
	Bids   *sync.Map // 价格 -> 数量
	Asks   *sync.Map
}

// 成交
type Trade struct {
	ID        string
	Symbol    string
	Price     *big.Float
	Quantity  *big.Float
	BuyerID   string
	SellerID  string
	Timestamp time.Time
}
