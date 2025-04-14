package pkg

import "context"

type MessageContent interface {
	AsString() (string, error)
	AsBytes() ([]byte, error)
}

type BrokerConsumer interface {
	Read(ctx context.Context) (MessageContent, error)
}
