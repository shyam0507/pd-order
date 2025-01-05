package main

import (
	"log/slog"
	"os"

	"github.com/shyam0507/pd-order/src/src/api"
	"github.com/shyam0507/pd-order/src/src/storage"
)

func main() {

	kafkaBrokers := os.Getenv("KAFKA_BROKERS") //localhost:9092
	mongoURI := os.Getenv("MONGO_URI")         //mongodb://localhost:27017

	slog.Info("Kafka Brokers: ", "", kafkaBrokers)
	slog.Info("Mongo URI: ", "", mongoURI)

	mongoStorage := storage.NewMongoStorage(mongoURI, "pd_orders")
	orderCreatedproducer := storage.NewKafkaProducer("order.created", []string{kafkaBrokers})
	orderConfirmedProducer := storage.NewKafkaProducer("order.confirmed", []string{kafkaBrokers})
	consumer := storage.NewKafkaConsumer("payment.received", []string{kafkaBrokers}, mongoStorage, orderConfirmedProducer)

	go consumer.ConsumePaymentReceived()

	server := api.NewServer("3005", mongoStorage, orderCreatedproducer)
	server.Start()
}
