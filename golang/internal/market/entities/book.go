package entities

import (
	"container/heap"
	"sync"
)

type Book struct {
	Orders        []*Order
	Transactions  []*Transaction
	OrdersChan    chan *Order
	OrdersChanOut chan *Order
	Wg            *sync.WaitGroup
}

func NewBook(orderChan chan *Order, orderChanOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Orders:        []*Order{},
		Transactions:  []*Transaction{},
		OrdersChan:    orderChan,
		OrdersChanOut: orderChanOut,
		Wg:            wg,
	}
}

func (b *Book) Trade() {
	buyOrders := make(map[string]*OrderQueue)
	sellOrders := make(map[string]*OrderQueue)

	for order := range b.OrdersChan {
		asset := order.Asset.ID
		buyOrder := buyOrders[asset]
		sellOrder := sellOrders[asset]

		buyOrder = InitNewOrder(buyOrder)
		sellOrder = InitNewOrder(sellOrder)

		if order.OrderType == "BUY" {
			buyOrder.Push(order)
			b.TradeOption(sellOrder, order)
		} else if order.OrderType == "SELL" {
			sellOrder.Push(order)
			b.TradeOption(buyOrder, order)
		}
	}
}

func InitNewOrder(optionOrder *OrderQueue) *OrderQueue {
	if optionOrder == nil {
		optionOrder = NewOrderQueue()
		heap.Init(optionOrder)
	}
	return optionOrder
}

func (b *Book) TradeOption(optionOrders *OrderQueue, order *Order) {
	if optionOrders.Len() > 0 && optionOrders.Orders[0].Price >= order.Price {
		optionOrder := optionOrders.Pop().(*Order)
		if optionOrder.PendingShares > 0 {
			transaction := NewTransaction(order, optionOrder, order.Shares, optionOrder.Price)
			b.AddTransaction(transaction, b.Wg)

			optionOrder.Transactions = append(optionOrder.Transactions, transaction)
			order.Transactions = append(order.Transactions, transaction)

			b.OrdersChanOut <- optionOrder
			b.OrdersChanOut <- order

			if optionOrder.PendingShares > 0 {
				optionOrders.Push(optionOrder)
			}
		}
	}
}

func (b *Book) AddTransaction(transaction *Transaction, wg *sync.WaitGroup) {
	defer wg.Done()

	sellingShares := transaction.SellingOrder.PendingShares
	buyingShares := transaction.BuyingOrder.PendingShares

	minShares := GetMinShares(buyingShares, sellingShares)

	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.ID, -minShares)
	transaction.AddSellOrderPendingShares(-minShares)
	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.BuyingOrder.Asset.ID, minShares)
	transaction.AddBuyOrderPendingShares(-minShares)

	transaction.CalculateTotal()

	CloseTransaction(transaction.BuyingOrder)
	CloseTransaction(transaction.SellingOrder)

	b.Transactions = append(b.Transactions, transaction)
}
