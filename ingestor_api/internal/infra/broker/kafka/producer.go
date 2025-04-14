package kafka

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

var ErrProducerNotInitialized = errors.New("kafka producer not initialized")

func NewProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *KafkaProducer) Write(ctx context.Context, key any, value any) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	stringKey, ok := key.(string)
	if !ok {
		return errors.New("String must be a string")
	}

	msg := kafka.Message{
		Key:   []byte(stringKey),
		Value: []byte(string(bytes)),
	}

	err = p.writer.WriteMessages(ctx, msg)
	return err
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}

func InitProducer(brokers []string, topic string) *KafkaProducer {
	producer := NewProducer(brokers, topic)
	return producer
}
