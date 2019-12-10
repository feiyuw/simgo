package protocols

import (
	"context"

	"simgo/logger"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// ================================== client ==================================

// client
type GrpcClient struct {
	addr string // connected service addr, eg. 127.0.0.1:1988
	conn *grpc.ClientConn
	desc grpcurl.DescriptorSource
}

// Create a new grpc client
func NewGrpcClient(addr string, protos []string, opts ...grpc.DialOption) *GrpcClient {
	var descSource grpcurl.DescriptorSource

	if len(protos) > 0 {
		descSource, err := grpcurl.DescriptorSourceFromProtoFiles([]string{}, protos...)
		if err != nil {
			logger.Fatalf("protocols/grpc", "cannot parse proto file: %v", err)
		}
		return &GrpcClient{addr: addr, desc: descSource}
	}

	if addr == "" {
		logger.Fatal("protocols/grpc", "addr or protos should be set")
	}

	// fetch from server reflection RPC
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		logger.Fatalf("protocols/grpc", "did not connect: %v", err)
	}
	c := rpb.NewServerReflectionClient(conn)
	ctx := context.Background() // TODO: add timeout, dialtime options
	refClient := grpcreflect.NewClient(ctx, c)
	descSource = grpcurl.DescriptorSourceFromServer(ctx, refClient)

	return &GrpcClient{addr: addr, desc: descSource}
}

func (gc *GrpcClient) ListServices() ([]string, error) {
	svcs, err := grpcurl.ListServices(gc.desc)
	if err != nil {
		return nil, err
	}
	return svcs, nil
}

func (gc *GrpcClient) ListMethods(svcName string) ([]string, error) {
	mtds, err := grpcurl.ListMethods(gc.desc, svcName)
	if err != nil {
		return nil, err
	}
	return mtds, nil
}

// ================================== server ==================================

type GrpcServer struct {
}

// create a new grpc server
func NewGrpcServer(addr string) *GrpcServer {
	return &GrpcServer{}
}
