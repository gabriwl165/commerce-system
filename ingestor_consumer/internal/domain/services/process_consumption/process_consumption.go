package process_consumption

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/gabriwl165/commerce-system/internal/pkg"
)

type Aggregator struct {
	mu   sync.Mutex
	data map[string]float64
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		data: make(map[string]float64),
	}
}

func (a *Aggregator) Add(tenant, productSKU, useUnity string, usedAmount float64) {
	key := fmt.Sprintf("%s:%s:%s", tenant, productSKU, useUnity)
	a.mu.Lock()
	defer a.mu.Unlock()
	a.data[key] += usedAmount
}

func (a *Aggregator) GetData() []map[string]interface{} {
	a.mu.Lock()
	defer a.mu.Unlock()
	result := make([]map[string]interface{}, 0)
	for key, value := range a.data {
		parts := strings.Split(key, ":")
		tenant, product, use_unity := parts[0], parts[1], parts[2]
		result = append(result, map[string]interface{}{
			"tenant":      tenant,
			"product":     product,
			"use_unity":   use_unity,
			"used_amount": value,
		})

	}
	return result
}

func ProcessConsumption(producer pkg.BrokerProducer, consumptionChan <-chan map[string]interface{}) {
	agg := NewAggregator()
	for consumption := range consumptionChan {
		tenant := consumption["tenant"].(string)
		productSKU := consumption["product_sku"].(string)
		useUnity := consumption["use_unity"].(string)
		usedAmount := consumption["used_amount"].(float64)

		agg.Add(tenant, productSKU, useUnity, usedAmount)
	}
	pulses := agg.GetData()
	for _, pulse := range pulses {
		tenant := pulse["tenant"].(string)
		producer.Write(context.Background(), tenant, pulse)
	}

}
