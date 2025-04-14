package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gabriwl165/commerce-system/infra/env"
	"github.com/gabriwl165/commerce-system/internal/domain/services/process_consumption"
	"github.com/gabriwl165/commerce-system/internal/domain/services/pulses_consumer"
	"github.com/gabriwl165/commerce-system/internal/infra/broker/kafka"
	"github.com/gabriwl165/commerce-system/internal/pkg"
)

func main() {
	envManager := env.Init()
	kafkaBroker, _ := envManager.Read("KAFKA_BROKER")
	brokers := []string{
		kafkaBroker.(string),
	}
	topic_producer := "processor-storage-consumption"
	producer := kafka.InitProducer(brokers, topic_producer)
	StartScheduleConsumer(producer, envManager)
}

func StartScheduleConsumer(producer pkg.BrokerProducer, envManager *env.EnvManager) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	for {

		consumptionChan := make(chan map[string]interface{})
		go process_consumption.ProcessConsumption(producer, consumptionChan)

		select {
		case <-ctx.Done():
			fmt.Println("Shutdown signal received, exiting...")
			return

		default:

			kafkaBroker, _ := envManager.Read("KAFKA_BROKER")
			brokers := []string{kafkaBroker.(string)}
			topic_consumer := "resource-consumption"

			consumer := &kafka.KafkaBrokerConsumer{}
			consumer.Init(brokers, topic_consumer, "ingestor-consumer")

			log.Print("Running Consumer.")
			pulses_consumer.StartConsumer(consumer, consumptionChan, time.Duration(60))
			close(consumptionChan)
			time.Sleep(10 * time.Second)

		}
	}
}
