package match

import (
	"exchange-tutorials/types"
	"exchange-tutorials/utils"
	_ "log"
	"math/big"
	"sync"
	"time"
)

type MatchingEngine struct {
	orderBooks sync.Map // map[string]*OrderBook
	trades     chan types.Trade
}

func NewMatchingEngine() *MatchingEngine {
	return &MatchingEngine{
		trades: make(chan types.Trade, 1000),
	}
}

func (me MatchingEngine) Trades() chan types.Trade {
	return me.trades
}

func (me *MatchingEngine) AddOrder(order *types.Order) ([]types.Trade, error) {
	book, _ := me.orderBooks.LoadOrStore(order.Symbol, &types.OrderBook{
		Symbol: order.Symbol,
		Bids:   &sync.Map{},
		Asks:   &sync.Map{},
	})
	return me.matchOrder(order, book.(*types.OrderBook))
}

func (me *MatchingEngine) matchOrder(order *types.Order, book *types.OrderBook) ([]types.Trade, error) {
	var trades []types.Trade
	remaining := new(big.Float).Set(order.Quantity)

	if order.Side == "buy" {
		asks := me.getSortedPrices(book.Asks, true)
		for _, priceStr := range asks {
			price, _ := new(big.Float).SetString(priceStr)
			if order.Type == "limit" && order.Price.Cmp(price) < 0 {
				break
			}
			qty, _ := book.Asks.Load(priceStr)
			tradeQty := new(big.Float).Set(qty.(*big.Float))
			if tradeQty.Cmp(remaining) > 0 {
				tradeQty.Set(remaining)
			}
			if tradeQty.Cmp(big.NewFloat(0)) > 0 {
				trade := types.Trade{
					ID:        utils.GenerateID(),
					Symbol:    order.Symbol,
					Price:     price,
					Quantity:  tradeQty,
					BuyerID:   order.UserID,
					SellerID:  "seller", // 假设
					Timestamp: time.Now(),
				}
				trades = append(trades, trade)
				me.trades <- trade
				remaining.Sub(remaining, tradeQty)
				book.Asks.Store(priceStr, new(big.Float).Sub(qty.(*big.Float), tradeQty))
			}
			if remaining.Cmp(big.NewFloat(0)) <= 0 {
				order.Status = "filled"
				return trades, nil
			}
		}
		if remaining.Cmp(big.NewFloat(0)) > 0 {
			book.Bids.Store(order.Price.String(), remaining)
			order.Filled = new(big.Float).Sub(order.Quantity, remaining)
			order.Status = "open"
		}
	} // 卖单逻辑类似
	return trades, nil
}

func (me *MatchingEngine) getSortedPrices(m *sync.Map, asc bool) []string {
	// 实现价格排序 (省略，类似之前)
	return []string{}
}
