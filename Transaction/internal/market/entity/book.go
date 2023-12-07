package entity

import (
	"container/heap"
	"sync"
)

type Book struct {
	Orders			[]*Order
	Transactions	[]*Transaction
	OrderChan		chan *Order
	OrderChanOut	chan *Order
	Wg				*sync.WaitGroup  //Espera todas as trades terminarem para emitir um sinal
}

func NewBook(orderChan chan *Order, orderChanOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Orders:			[]*Order{},
		Transactions:	[]*Transaction{},
		OrderChan:		orderChan,
		OrderChanOut:	orderChanOut,
		Wg: 			wg,
	}
}

func (b *Book) Trade() {
	buyOrders := make(map[string]*OrderQueue)
	sellOrders := make(map[string]*OrderQueue)

	//buyOrders := NewOrderQueue()
	//sellOrders := NewOrderQueue()

	//heap.Init(buyOrders)
	//heap.Init(sellOrders)

	for order := range b.OrderChan {
		asset := order.Asset.ID

		if buyOrders[asset] == nil {
			buyOrders[asset] = NewOrderQueue()
			heap.Init(buyOrders[asset])
		}

		if sellOrders[asset] == nil {
			sellOrders[asset] = NewOrderQueue()
			heap.Init(sellOrders[asset])
		}

		if order.OrderType == "BUY" {
			buyOrders[asset].Push(order)

			if sellOrders[asset].Len() > 0 && sellOrders[asset].Orders[0].Price <= order.Price {
				sellOrder := sellOrders[asset].Pop().(*Order)

				if sellOrder.PendingShares > 0 {
					sealedShares := min(sellOrder.PendingShares, order.PendingShares)

					transaction := NewTransaction(sellOrder, order, sealedShares, sellOrder.Price)
					b.AddTransaction(transaction, b.Wg, sealedShares)
					sellOrder.Transactions = append(sellOrder.Transactions, transaction)
					order.Transactions = append(order.Transactions, transaction)

					b.OrderChanOut <- sellOrder
					b.OrderChanOut <- order

					if sellOrder.PendingShares > 0 {
						sellOrders[asset].Push(sellOrder)
					}
				}
			}
		} else if order.OrderType == "SELL" {
			sellOrders[asset].Push(order)

			if buyOrders[asset].Len() > 0 && buyOrders[asset].Orders[0].Price >= order.Price {
				buyOrder := buyOrders[asset].Pop().(*Order)

				if buyOrder.PendingShares > 0 {
					sealedShares := min(buyOrder.PendingShares, order.PendingShares)
					
					transaction := NewTransaction(order, buyOrder, sealedShares, order.Price)
					
					b.AddTransaction(transaction, b.Wg, sealedShares)

					buyOrder.AppendTransaction(transaction)
					order.AppendTransaction(transaction)

					b.OrderChanOut <- buyOrder
					b.OrderChanOut <- order

					if buyOrder.PendingShares > 0 {
						buyOrders[asset].Push(buyOrder)
					}
				}
			}
		}
	}
}

func (b *Book) AddTransaction(transaction *Transaction, wg *sync.WaitGroup, soldShares int) {
	defer wg.Done() // defer estabelece que esse será o último comando desse bloco de código

	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.ID, -soldShares)
	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.ID, soldShares)

	transaction.SellingOrder.SellShares(soldShares)
	transaction.BuyingOrder.SellShares(soldShares)

	transaction.CalculateTotal()

	transaction.SellingOrder.Close()
	transaction.BuyingOrder.Close()

	b.Transactions = append(b.Transactions, transaction)

}
