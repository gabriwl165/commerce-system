package main

import (
	"log"
	"net"

	"github.com/gabriwl165/commerce-system/infra/env"
	"github.com/gabriwl165/commerce-system/internal/infra/kafka"
	"github.com/gabriwl165/commerce-system/internal/interfaces/grpc_handler"
	"github.com/gabriwl165/commerce-system/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	envManager := env.Init()
	kafkaBroker, _ := envManager.Read("KAFKA_BROKER")
	// Create a new gRPC server instance
	grpcServer := grpc.NewServer()
	brokers := []string{
		kafkaBroker.(string),
	}
	topic_producer := "processor-storage-consumption"
	producer := kafka.InitProducer(brokers, topic_producer)
	server := grpc_handler.UsageServiceServer{
		BrokerProducer: producer,
	}

	// Register the Greeter server
	proto.RegisterUsageServiceServer(grpcServer, &server)

	// Log that the server is running
	log.Printf("Server listening at %v", lis.Addr())

	reflection.Register(grpcServer)
	kafka.InitProducer(
		brokers,
		"resource-consumption",
	)

	// Start serving incoming connections
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
