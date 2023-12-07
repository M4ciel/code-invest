package entities

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuyAsset(t *testing.T) {
	asset1 := NewAsset("asset1", "Asset 1", 100)

	investor := NewInvestor("1")
	investor2 := NewInvestor("2")

	investorAssetPosition := NewAssetPosition("asset1", 10)
	investor.AddAssetPosition(investorAssetPosition)

	wg := sync.WaitGroup{}
	orderChan := make(chan *Order)
	orderChanOut := make(chan *Order)

	book := NewBook(orderChan, orderChanOut, &wg)
	go book.Trade()

	// add buy order
	wg.Add(1)
	order := NewOrder("1", investor, asset1, 5, 5, "SELL")
	orderChan <- order

	// add sell order

	order2 := NewOrder("2", investor2, asset1, 5, 5, "BUY")
	orderChan <- order2
	wg.Wait()

	asset := assert.New(t)
	asset.Equal("CLOSED", order.Status, "Order 1 should be closed")
	asset.Equal(0, order.PendingShares, "Order 1 should have 0 PendingShares")
	asset.Equal("CLOSED", order2.Status, "Order 2 should be closed")
	asset.Equal(0, order2.PendingShares, "Order 2 should have 0 PendingShares")

	asset.Equal(5, investorAssetPosition.Shares, "Investor 1 should have 5 shares of asset 1")
	asset.Equal(5, investor2.GetAssetPosition("asset1").Shares, "Investor 2 should have 5 shares of asset 1")
}

func TestBuyAssetWithDifferentAssents(t *testing.T) {
	asset1 := NewAsset("asset1", "Asset 1", 100)
	asset2 := NewAsset("asset2", "Asset 2", 100)

	investor := NewInvestor("1")
	investor2 := NewInvestor("2")

	investorAssetPosition := NewAssetPosition("asset1", 10)
	investor.AddAssetPosition(investorAssetPosition)

	investorAssetPosition2 := NewAssetPosition("asset2", 10)
	investor2.AddAssetPosition(investorAssetPosition2)

	wg := sync.WaitGroup{}
	orderChan := make(chan *Order)
	orderChanOut := make(chan *Order)

	book := NewBook(orderChan, orderChanOut, &wg)
	go book.Trade()

	order := NewOrder("1", investor, asset1, 5, 5, "SELL")
	orderChan <- order

	order2 := NewOrder("2", investor2, asset2, 5, 5, "BUY")
	orderChan <- order2

	asset := assert.New(t)
	asset.Equal("OPEN", order.Status, "Order 1 should be closed")
	asset.Equal(5, order.PendingShares, "Order 1 should have 5 PendingShares")
	asset.Equal("OPEN", order2.Status, "Order 2 should be closed")
	asset.Equal(5, order2.PendingShares, "Order 2 should have 5 PendingShares")
}

func TestBuyPartialAsset(t *testing.T) {
	asset1 := NewAsset("asset1", "Asset 1", 100)

	investor := NewInvestor("1")
	investor2 := NewInvestor("2")
	investor3 := NewInvestor("3")

	investorAssetPosition := NewAssetPosition("asset1", 3)
	investor.AddAssetPosition(investorAssetPosition)

	investorAssetPosition2 := NewAssetPosition("asset1", 5)
	investor3.AddAssetPosition(investorAssetPosition2)

	wg := sync.WaitGroup{}
	orderChan := make(chan *Order)
	orderChanOut := make(chan *Order)

	book := NewBook(orderChan, orderChanOut, &wg)
	go book.Trade()

	wg.Add(1)
	// investidor 2 quer comprar 5 shares
	order2 := NewOrder("1", investor2, asset1, 5, 5.0, "BUY")
	orderChan <- order2

	// investidor 1 quer vender 3 shares
	order := NewOrder("2", investor, asset1, 3, 5.0, "SELL")
	orderChan <- order

	asset := assert.New(t)
	go func() {
		for range orderChanOut {
		}
	}()

	wg.Wait()

	// assert := assert.New(t)
	asset.Equal("CLOSED", order.Status, "Order 1 should be closed")
	asset.Equal(0, order.PendingShares, "Order 1 should have 0 PendingShares")

	asset.Equal("OPEN", order2.Status, "Order 2 should be OPEN")
	asset.Equal(2, order2.PendingShares, "Order 2 should have 2 PendingShares")

	asset.Equal(0, investorAssetPosition.Shares, "Investor 1 should have 0 shares of asset 1")
	asset.Equal(3, investor2.GetAssetPosition("asset1").Shares, "Investor 2 should have 3 shares of asset 1")

	wg.Add(1)
	order3 := NewOrder("3", investor3, asset1, 2, 5.0, "SELL")
	orderChan <- order3
	wg.Wait()

	asset.Equal("CLOSED", order3.Status, "Order 3 should be closed")
	assert.Equal(&testing.T{}, order3.PendingShares, "Order 3 should have 0 PendingShares")

	asset.Equal("CLOSED", order2.Status, "Order 2 should be CLOSED")
	asset.Equal(0, order2.PendingShares, "Order 2 should have 0 PendingShares")

	asset.Equal(2, len(book.Transactions), "Should have 2 transactions")
	asset.Equal(15.0, book.Transactions[0].Total, "Transaction should have price 15")
	asset.Equal(10.0, book.Transactions[1].Total, "Transaction should have price 10")
}

func TestBuyWithDifferentPrice(t *testing.T) {
	asset1 := NewAsset("asset1", "Asset 1", 100)

	investor := NewInvestor("1")
	investor2 := NewInvestor("2")
	investor3 := NewInvestor("3")

	investorAssetPosition := NewAssetPosition("asset1", 3)
	investor.AddAssetPosition(investorAssetPosition)

	investorAssetPosition2 := NewAssetPosition("asset1", 5)
	investor3.AddAssetPosition(investorAssetPosition2)

	wg := sync.WaitGroup{}
	orderChan := make(chan *Order)

	orderChanOut := make(chan *Order)

	book := NewBook(orderChan, orderChanOut, &wg)
	go book.Trade()

	wg.Add(1)
	// investidor 2 quer comprar 5 shares
	order2 := NewOrder("2", investor2, asset1, 5, 5.0, "BUY")
	orderChan <- order2

	// investidor 1 quer vender 3 shares
	order := NewOrder("1", investor, asset1, 3, 4.0, "SELL")
	orderChan <- order

	go func() {
		for range orderChanOut {
		}
	}()
	wg.Wait()

	asset := assert.New(t)
	asset.Equal("CLOSED", order.Status, "Order 1 should be closed")
	asset.Equal(0, order.PendingShares, "Order 1 should have 0 PendingShares")

	asset.Equal("OPEN", order2.Status, "Order 2 should be OPEN")
	asset.Equal(2, order2.PendingShares, "Order 2 should have 2 PendingShares")

	asset.Equal(0, investorAssetPosition.Shares, "Investor 1 should have 0 shares of asset 1")
	asset.Equal(3, investor2.GetAssetPosition("asset1").Shares, "Investor 2 should have 3 shares of asset 1")

	wg.Add(1)
	order3 := NewOrder("3", investor3, asset1, 3, 4.5, "SELL")
	orderChan <- order3

	wg.Wait()

	asset.Equal("OPEN", order3.Status, "Order 3 should be open")
	asset.Equal(1, order3.PendingShares, "Order 3 should have 1 PendingShares")

	asset.Equal("CLOSED", order2.Status, "Order 2 should be CLOSED")
	asset.Equal(0, order2.PendingShares, "Order 2 should have 0 PendingShares")

	// assert.Equal(2, len(book.Transactions), "Should have 2 transactions")
	// assert.Equal(15.0, float64(book.Transactions[0].Total), "Transaction should have price 15")
	// assert.Equal(10.0, float64(book.Transactions[1].Total), "Transaction should have price 10")
}

func TestNoMatch(t *testing.T) {
	asset1 := NewAsset("asset1", "Asset 1", 100)

	investor := NewInvestor("1")
	investor2 := NewInvestor("2")

	investorAssetPosition := NewAssetPosition("asset1", 3)
	investor.AddAssetPosition(investorAssetPosition)

	wg := sync.WaitGroup{}
	orderChan := make(chan *Order)

	orderChanOut := make(chan *Order)

	book := NewBook(orderChan, orderChanOut, &wg)
	go book.Trade()

	wg.Add(0)
	// investidor 1 quer vender 3 shares
	order := NewOrder("1", investor, asset1, 3, 6.0, "SELL")
	orderChan <- order

	// investidor 2 quer comprar 5 shares
	order2 := NewOrder("2", investor2, asset1, 5, 5.0, "BUY")
	orderChan <- order2

	go func() {
		for range orderChanOut {
		}
	}()
	wg.Wait()

	asset := assert.New(t)
	asset.Equal("OPEN", order.Status, "Order 1 should be closed")
	asset.Equal("OPEN", order2.Status, "Order 2 should be OPEN")
	asset.Equal(3, order.PendingShares, "Order 1 should have 3 PendingShares")
	asset.Equal(5, order2.PendingShares, "Order 2 should have 5 PendingShares")
}
