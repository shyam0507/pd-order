package storage

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"github.com/shyam0507/pd-order/types"
)

type KafkaProducer struct {
	writer *kafka.Writer
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

func NewKafkaProducer(topic string, brokers []string) Producer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return KafkaProducer{writer: w}
}
