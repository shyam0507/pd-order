package types

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Products   []Product          `json:"products" bson:"products"`
	CustomerId string             `json:"customer_id" bson:"customer_id"`
	AddressId  string             `json:"address_id" bson:"address_id"`
	Status     string             `json:"status" bson:"status"`
	Total      float64            `json:"total,omitempty" bson:"total"`
}

type Product struct {
	ProductID string `json:"product_id" bson:"product_id"`
	Quantity  int    `json:"quantity" bson:"quantity"`
}

func (o *Order) Validate() error {
	if len(o.Products) < 1 {
		return fmt.Errorf("at least one product is required for an order")
	}

	if o.CustomerId == "" {
		return fmt.Errorf("customer id is required")
	}

	if o.AddressId == "" {
		return fmt.Errorf("address is required")
	}

	return nil
}

func (o *Order) CalculateTotal() (float64, error) {
	return 100, nil
}
