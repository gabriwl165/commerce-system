package grpc_handler

import (
	"context"

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
	kafka.SendEvent(context.Background(), "resource-consumption", bytes)

	return nil, nil
}
