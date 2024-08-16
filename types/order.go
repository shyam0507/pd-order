package types

type Order struct {
	OrderID    string    `json:"id" bson:"_id",omitempty`
	Products   []Product `json:"products"`
	CustomerId string
	AddressId  string
	Status     string
	Total      float64
}

type Product struct {
	ProductID string
	Quantity  int
}
