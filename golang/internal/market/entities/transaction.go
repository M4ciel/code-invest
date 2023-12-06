package entities

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID           string
	SellingOrder *Order
	BuyingOrder  *Order
	Shares       int
	Price        float64
	Total        float64
	DateTime     time.Time
}

func NewTransaction(sellingOrder *Order, buyingOrder *Order, shares int, price float64) *Transaction {
	total := float64(shares) * price

	return &Transaction{
		ID:           uuid.New().String(),
		SellingOrder: sellingOrder,
		BuyingOrder:  buyingOrder,
		Shares:       shares,
		Price:        price,
		Total:        total,
		DateTime:     time.Now(),
	}
}

func CloseTransaction(optionOrder *Order) {
	if optionOrder.PendingShares == 0 {
		optionOrder.Status = "CLOSED"
	}
}

func GetMinShares(buyingShares int, sellingShares int) int {
	minShares := sellingShares
	if buyingShares < sellingShares {
		minShares = buyingShares
	}

	return minShares
}

func (t *Transaction) CalculateTotal() {
	t.Total = float64(t.Shares) * t.Price
}

func (t *Transaction) AddBuyOrderPendingShares(shares int) {
	t.BuyingOrder.PendingShares += shares
}

func (t *Transaction) AddSellOrderPendingShares(shares int) {
	t.SellingOrder.PendingShares += shares
}
