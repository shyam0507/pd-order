package main

import (
	"github.com/shyam0507/pd-order/api"
	"github.com/shyam0507/pd-order/storage"
)

func main() {

	mongoStorage := storage.NewMongoStorage("mongodb://localhost:27017", "pd_orders")
	orderCreatedproducer := storage.NewKafkaProducer("order.created", []string{"localhost:9092"})
	orderConfirmedProducer := storage.NewKafkaProducer("order.confirmed", []string{"localhost:9092"})
	consumer := storage.NewKafkaConsumer("payment.received", []string{"localhost:9092"}, mongoStorage, orderConfirmedProducer)

	go consumer.ConsumePaymentReceived()

	server := api.NewServer("3005", mongoStorage, orderCreatedproducer)
	server.Start()
}
