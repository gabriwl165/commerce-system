package main

import (
	"log"
	"net"

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

	// Create a new gRPC server instance
	grpcServer := grpc.NewServer()

	server := grpc_handler.UsageServiceServer{}

	// Register the Greeter server
	proto.RegisterUsageServiceServer(grpcServer, &server)

	// Log that the server is running
	log.Printf("Server listening at %v", lis.Addr())

	reflection.Register(grpcServer)

	brokers := []string{"localhost:9092"}
	kafka.InitProducer(
		brokers,
		"resource-consumption",
	)

	// Start serving incoming connections
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
