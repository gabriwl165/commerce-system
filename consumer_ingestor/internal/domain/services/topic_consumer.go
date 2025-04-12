package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func StartConsumer(kafkaReader *kafka.Reader, consumptionChan chan<- map[string]interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer kafkaReader.Close()

	for {
		m, err := kafkaReader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				break
			}
			fmt.Println("read error:", err)
			return
		}

		jsonStr := string(m.Value) // treat message as JSON string
		var value map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &value); err != nil {
			log.Printf("failed to unmarshal JSON string: %v", err)
			return
		}

		consumptionChan <- value
	}

	log.Print("Ending Consumer for some time")
}
