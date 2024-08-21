package main

import (
	"github.com/shyam0507/pd-order/api"
	"github.com/shyam0507/pd-order/storage"
)

func main() {

	mongoStorage := storage.NewMongoStorage("mongodb://localhost:27017", "pd_orders")
	producer := storage.NewKafkaProducer("order.created", []string{"localhost:9092"})

	server := api.NewServer("3005", mongoStorage, producer)
	server.Start()
}
