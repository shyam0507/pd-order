package storage

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/shyam0507/pd-order/src/src/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

type KafkaConsumer struct {
	reader   *kafka.Reader
	storage  Storage
	producer Producer
}

// ProduceOrderCreated implements Producer.
func (k KafkaProducer) ProduceOrderCreated(key string, value types.OrderCreatedEvent) error {

	m, _ := json.Marshal(value)
	err := k.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: m,
		},
	)

	if err != nil {
		slog.Error("failed to write messages:", "Err", err)
		return err
	}

	slog.Info("kafka event published for Key", "Key", key)

	return nil
}

func (k KafkaProducer) ProduceOrderConfirmed(key string, value types.OrderConfirmedEvent) error {

	m, _ := json.Marshal(value)
	err := k.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: m,
		},
	)

	if err != nil {
		slog.Error("failed to write messages:", "Err", err)
		return err
	}

	slog.Info("kafka event published for Key", "Key", key)

	return nil
}

func NewKafkaProducer(topic string, brokers []string) Producer {
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")
	mechanism := plain.Mechanism{
		Username: username,
		Password: password,
	}

	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Transport: &kafka.Transport{
			SASL: mechanism,
		},
		AllowAutoTopicCreation: true,
	}

	return &KafkaProducer{writer: w}
}

func NewKafkaConsumer(topic string, brokers []string, storage Storage, producer Producer) Consumer {
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")
	mechanism := plain.Mechanism{
		Username: username,
		Password: password,
	}

	dialer := &kafka.Dialer{
		SASLMechanism: mechanism,
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     "pd-order-group",
		StartOffset: kafka.LastOffset,
		MaxBytes:    10e6, // 10MB
		Dialer:      dialer,
	})

	return KafkaConsumer{reader: r, storage: storage, producer: producer}
}

// ConsumePaymentReceived implements Consumer.
//
// It will block until a message is consumed, then it will return the value of the message.
// If the consumer is closed, it will return an error.
func (k KafkaConsumer) ConsumePaymentReceived() error {
	slog.Info("Consuming payment received event")
	for {
		m, err := k.reader.ReadMessage(context.Background())

		if err != nil {
			slog.Error("failed to read message:", "Err", err)
			continue
		}

		slog.Info("kafka event consumed for Key", "Key", string(m.Key))
		var event types.PaymentReceivedEvent
		err = json.Unmarshal(m.Value, &event)

		if err != nil {
			slog.Error("failed to unmarshal message:", "Err", err)
			continue
		}

		slog.Info("payment received event consumed", "Event", event)

		if err := k.storage.UpdateOrder(event.Data.OrderId, "CONFIRMED"); err != nil {
			slog.Error("failed to update order status:", "Err", err)
			continue
		}
		slog.Info("order status updated to ORDER_CONFIRMED", "OrderId", event.Data.OrderId)

		id, _ := primitive.ObjectIDFromHex(event.Data.OrderId)

		//TODO product event for order confirmed
		confEvent := types.OrderConfirmedEvent{
			Id:              event.Data.OrderId,
			Type:            "OrderConfirmed",
			Source:          "OrderService",
			SpecVersion:     "1.0",
			Data:            types.Order{Id: id, Status: "ORDER_CONFIRMED"},
			DataContentType: "application/json",
		}
		k.producer.ProduceOrderConfirmed("order.confirmed", confEvent)

	}

}
