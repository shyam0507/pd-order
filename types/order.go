package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Products   []Product          `json:"products" bson:"products"`
	CustomerId string             `json:"customer_id" bson:"customer_id"`
	AddressId  string             `json:"address_id" bson:"address_id"`
	Status     string             `json:"status" bson:"status"`
	Total      float64            `json:"total" bson:"total"`
}

type Product struct {
	ProductID string `json:"product_id" bson:"product_id"`
	Quantity  int    `json:"quantity" bson:"quantity"`
}
