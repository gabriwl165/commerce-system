package kafka

import (
	"context"
	"errors"

	"github.com/gabriwl165/commerce-system/internal/pkg"
	"github.com/segmentio/kafka-go"
)

type KafkaBrokerConsumer struct {
	reader *kafka.Reader
}

func (k *KafkaBrokerConsumer) Init(brokers []string, topic string, groupID string) {
	k.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   topic,
	})
}

func (k *KafkaBrokerConsumer) Read(ctx context.Context) (pkg.MessageContent, error) {
	msg, err := k.reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	return &KafkaMessageContent{Message: msg}, nil
}

type KafkaMessageContent struct {
	Message kafka.Message
}

func (m *KafkaMessageContent) AsString() (string, error) {
	if m.Message.Value == nil {
		return "", errors.New("message value is nil")
	}
	return string(m.Message.Value), nil
}

func (m *KafkaMessageContent) AsBytes() ([]byte, error) {
	if m.Message.Value == nil {
		return nil, errors.New("message value is nil")
	}
	return m.Message.Value, nil
}
