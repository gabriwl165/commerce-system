package pkg

import "context"

type BrokerProducer interface {
	Write(ctx context.Context, key any, value any) error
}
