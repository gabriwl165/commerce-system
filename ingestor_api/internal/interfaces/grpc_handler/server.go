package grpc_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gabriwl165/commerce-system/internal/pkg"
	"github.com/gabriwl165/commerce-system/proto"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
)

// server is used to implement helloworld.GreeterServer.
type UsageServiceServer struct {
	proto.UsageServiceServer
	BrokerProducer pkg.BrokerProducer
}

func (s *UsageServiceServer) Consume(ctx context.Context, usageInfo *proto.UsageInfoRequest) (*emptypb.Empty, error) {
	marshaller := protojson.MarshalOptions{
		UseProtoNames:   true, // field names same as proto
		EmitUnpopulated: true, // include zero-value fields
	}

	bytes, _ := marshaller.Marshal(usageInfo)
	var result map[string]interface{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	log.Print("Producing message into Kafka")
	tenant := result["tenant"].(string)
	s.BrokerProducer.Write(context.Background(), tenant, result)
	return nil, nil
}
