package storage

import "github.com/shyam0507/pd-order/types"

type Storage interface {
	CreateOrder(types.Order) error
	UpdateOrder(id string, status string) error
}

type Producer interface {
	ProduceOrderCreated(string, types.OrderCreatedEvent) error
}

type Consumer interface {
}
