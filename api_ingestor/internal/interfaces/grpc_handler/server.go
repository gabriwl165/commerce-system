package grpc_handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gabriwl165/commerce-system/internal/infra/kafka"
	"github.com/gabriwl165/commerce-system/proto"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
)

// server is used to implement helloworld.GreeterServer.
type UsageServiceServer struct {
	proto.UsageServiceServer
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
	kafka.SendEvent(context.Background(), "resource-consumption", result)

	return nil, nil
}
