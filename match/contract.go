package match

import (
	"exchange-tutorials/types"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"
)

type FuturesEngine struct {
	Engine      *MatchingEngine
	users       sync.Map // map[string]*User
	liquidation chan string
}

func NewFuturesEngine() *FuturesEngine {
	fe := &FuturesEngine{
		Engine:      NewMatchingEngine(),
		liquidation: make(chan string, 100),
	}
	go fe.fundingRateLoop()
	return fe
}

func (fe *FuturesEngine) OpenPosition(order *types.Order) error {
	user, _ := fe.users.LoadOrStore(order.UserID, &types.User{
		Balance:   make(map[string]*big.Float),
		Margin:    make(map[string]*big.Float),
		Positions: make(map[string]*types.Position),
		Mode:      "isolated",
	})
	u := user.(*types.User)

	marginRequired := new(big.Float).Quo(
		new(big.Float).Mul(order.Price, order.Quantity),
		big.NewFloat(float64(order.Leverage)),
	)
	if u.Balance["USD"].Cmp(marginRequired) < 0 {
		return fmt.Errorf("insufficient balance")
	}

	u.Balance["USD"].Sub(u.Balance["USD"], marginRequired)
	if u.Mode == "isolated" {
		u.Margin[order.Symbol] = marginRequired
	} else {
		u.Margin["USD"].Add(u.Margin["USD"], marginRequired)
	}

	pos := &types.Position{
		Symbol:     order.Symbol,
		Side:       order.Side,
		Quantity:   order.Quantity,
		EntryPrice: order.Price,
		Leverage:   order.Leverage,
		Margin:     marginRequired,
	}
	u.Positions[order.Symbol] = pos

	trades, err := fe.Engine.AddOrder(order)
	if err != nil {
		return err
	}
	go fe.monitorLiquidation(u, pos, trades)
	return nil
}

func (fe *FuturesEngine) monitorLiquidation(user *types.User, pos *types.Position, trades []types.Trade) {
	latestPrice := trades[len(trades)-1].Price
	liqPrice := fe.calculateLiquidationPrice(user, pos)
	if (pos.Side == "long" && latestPrice.Cmp(liqPrice) <= 0) ||
		(pos.Side == "short" && latestPrice.Cmp(liqPrice) >= 0) {
		fe.liquidate(user, pos)
		fe.liquidation <- user.ID
	}
}

func (fe *FuturesEngine) calculateLiquidationPrice(user *types.User, pos *types.Position) *big.Float {
	mmr := big.NewFloat(0.05) // 维持保证金率 5%
	margin := pos.Margin
	if user.Mode == "cross" {
		margin = user.Margin["USD"]
	}
	maintenance := new(big.Float).Mul(margin, mmr)
	if pos.Side == "long" {
		return new(big.Float).Sub(pos.EntryPrice, new(big.Float).Quo(maintenance, pos.Quantity))
	}
	return new(big.Float).Add(pos.EntryPrice, new(big.Float).Quo(maintenance, pos.Quantity))
}

func (fe *FuturesEngine) liquidate(user *types.User, pos *types.Position) {
	delete(user.Positions, pos.Symbol)
	if user.Mode == "isolated" {
		user.Margin[pos.Symbol] = big.NewFloat(0)
	} else {
		user.Margin["USD"].Sub(user.Margin["USD"], pos.Margin)
	}
	log.Printf("User %s liquidated: %s", user.ID, pos.Symbol)
}

func (fe *FuturesEngine) fundingRateLoop() {
	ticker := time.NewTicker(8 * time.Hour) // 每 8 小时计算资金费率
	for range ticker.C {
		fe.applyFundingRate()
	}
}

func (fe *FuturesEngine) applyFundingRate() {
	rate := big.NewFloat(0.0001) // 示例资金费率 0.01%
	fe.users.Range(func(key, value interface{}) bool {
		user := value.(*types.User)
		for _, pos := range user.Positions {
			funding := new(big.Float).Mul(
				new(big.Float).Mul(pos.EntryPrice, pos.Quantity),
				rate,
			)
			if pos.Side == "long" {
				user.Balance["USD"].Sub(user.Balance["USD"], funding)
			} else {
				user.Balance["USD"].Add(user.Balance["USD"], funding)
			}
		}
		return true
	})
}
