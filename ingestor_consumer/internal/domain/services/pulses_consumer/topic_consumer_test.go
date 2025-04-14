package pulses_consumer

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/gabriwl165/commerce-system/internal/pkg"
	"github.com/stretchr/testify/assert"
)

type MockMessageContent struct{}

func (m *MockMessageContent) AsString() (string, error) {
	return `{
		"tenant": "user123",
		"product_sku": "vm",
		"used_amount": 3,
		"use_unity": "kb"
	}`, nil
}

func (m *MockMessageContent) AsBytes() ([]byte, error) {
	// Retrieve the JSON string using the AsString() method.
	jsonStr, err := m.AsString()
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON string into an interface to ensure it's valid.
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, err
	}

	// Marshal the data back into a compact JSON byte slice and return it.
	return json.Marshal(data)
}

type MockBrokerConsumer struct{}

func (m *MockBrokerConsumer) Read(ctx context.Context) (pkg.MessageContent, error) {
	select {
	case <-ctx.Done():
		return nil, nil
	default:
		return &MockMessageContent{}, nil
	}
}

func TestStartConsumer(t *testing.T) {
	var wg sync.WaitGroup

	consumptionChan := make(chan map[string]interface{})
	mockConsumer := &MockBrokerConsumer{}

	go func() {
		wg.Add(1)
		StartConsumer(mockConsumer, consumptionChan, time.Duration(5))
	}()

	select {
	case msg := <-consumptionChan:
		assert.NotNil(t, msg, "Expected a non-nil message from consumptionChan")
		wg.Done()
		break

	case <-time.After(10 * time.Second):
		assert.Fail(t, "Timeout", "Did not receive message on consumptionChan")
		wg.Done()
		break
	}
	wg.Wait()
	close(consumptionChan)
}
