package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

var ErrProducerNotInitialized = errors.New("kafka producer not initialized")

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) Produce(ctx context.Context, key string, value any) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(string(bytes)),
	}

	return p.writer.WriteMessages(ctx, msg)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

var (
	producer *Producer
	once     sync.Once
)

func InitProducer(brokers []string, topic string) {
	once.Do(func() {
		producer = NewProducer(brokers, topic)
	})
}

func SendEvent(ctx context.Context, key string, payload any) error {
	if producer == nil {
		return ErrProducerNotInitialized
	}
	return producer.Produce(ctx, key, payload)
}
