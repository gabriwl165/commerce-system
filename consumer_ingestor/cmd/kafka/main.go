package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gabriwl165/commerce-system/internal/domain/services"
	"github.com/gabriwl165/commerce-system/internal/infra/kafka"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	brokers := []string{"localhost:9092"}
	topic := "resource-consumption"

	consumptionChan := make(chan map[string]interface{})
	go services.ProcessConsumption(consumptionChan)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutdown signal received, exiting...")
			return
		default:
			kafkaConsumer := kafka.NewConsumer(brokers, topic, "ingestor-consumer")
			log.Print("Running Consumer.")
			services.StartConsumer(kafkaConsumer, consumptionChan)
			time.Sleep(10 * time.Second)
		}
	}
}
