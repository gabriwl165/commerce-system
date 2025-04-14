package process_consumption

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockProducer struct {
	CalledWithKey   any
	CalledWithValue map[string]interface{}
}

func (m *MockProducer) Write(ctx context.Context, key any, value any) error {
	m.CalledWithKey = key
	m.CalledWithValue = value.(map[string]interface{})
	return nil
}

func TestOneProcessConsumption(t *testing.T) {
	var wg sync.WaitGroup

	mockProducer := &MockProducer{}
	consumptionChan := make(chan map[string]interface{})

	for range []int{1} {
		go func() {
			defer wg.Done()
			wg.Add(1)
			ProcessConsumption(mockProducer, consumptionChan)
		}()
	}

	mockMessage := map[string]interface{}{
		"tenant":      "user123",
		"product_sku": "vm",
		"used_amount": 3.5,
		"use_unity":   "kb",
	}
	consumptionChan <- mockMessage
	close(consumptionChan)
	wg.Wait()
	assert.Equal(t, mockMessage["tenant"], mockProducer.CalledWithKey)
}

func TestManyProcessConsumption(t *testing.T) {
	var wg sync.WaitGroup

	mockProducer := &MockProducer{}
	consumptionChan := make(chan map[string]interface{})

	for range []int{1} {
		go func() {
			defer wg.Done()
			wg.Add(1)
			ProcessConsumption(mockProducer, consumptionChan)
		}()
	}

	max := 100_000.0
	usedAmount := 2.0
	tenant := "user123"

	for i := 0; i < int(max); i++ {
		mockMessage := map[string]interface{}{
			"tenant":      tenant,
			"product_sku": "vm",
			"used_amount": usedAmount,
			"use_unity":   "kb",
		}
		consumptionChan <- mockMessage
	}

	close(consumptionChan)
	wg.Wait()
	amount, _ := mockProducer.CalledWithValue["used_amount"].(float64)
	assert.Equal(t, max*usedAmount, amount)
}
