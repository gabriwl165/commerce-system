package pulses_consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gabriwl165/commerce-system/internal/pkg"
)

func StartConsumer(consumer pkg.BrokerConsumer, consumptionChan chan<- map[string]interface{}, max_read_time time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), max_read_time*time.Second)
	defer cancel()

	for {
		m, err := consumer.Read(ctx)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				break
			}
			fmt.Println("read error:", err)
			return
		}
		if m == nil {
			fmt.Println("received nil message")
			return
		}

		jsonStr, err := m.AsString()
		if err != nil {
			log.Print("Error while converting to string")
		}
		var value map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &value); err != nil {
			log.Printf("failed to unmarshal JSON string: %v", err)
			return
		}

		consumptionChan <- value
	}
}
