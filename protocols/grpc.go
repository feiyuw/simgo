package protocols

import (
	"context"
	"fmt"
	"time"

	"simgo/logger"

	"google.golang.org/grpc"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// ================================== client ==================================

// client
type GrpcClient struct {
	addr string // connected service addr, eg. 127.0.0.1:1988
	conn *grpc.ClientConn
}

// one endpoint contain multi services
type GrpcService struct {
	name string // service name
}

// one service contain multi methods
type GrpcMethod struct {
	name string
}

// Create a new grpc client
func NewGrpcClient(addr string, opts ...grpc.DialOption) *GrpcClient {
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		logger.Fatalf("protocols/grpc", "did not connect: %v", err)
	}
	return &GrpcClient{addr: addr, conn: conn}
}

func (gc *GrpcClient) ListServices() ([]*GrpcService, error) {
	c := rpb.NewServerReflectionClient(gc.conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := c.ServerReflectionInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get ServerReflectionInfo: %v", err)
	}
	if err := stream.Send(&rpb.ServerReflectionRequest{
		MessageRequest: &rpb.ServerReflectionRequest_ListServices{},
	}); err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	r, err := stream.Recv()
	if err != nil {
		// io.EOF is not ok.
		return nil, fmt.Errorf("failed to recv response: %v", err)
	}

	switch r.MessageResponse.(type) {
	case *rpb.ServerReflectionResponse_ListServicesResponse:
		services := r.GetListServicesResponse().Service
		grpcServices := make([]*GrpcService, len(services))
		for idx, svc := range services {
			grpcServices[idx] = &GrpcService{name: svc.Name}
		}
		return grpcServices, nil
	default:
		return nil, fmt.Errorf("ListServices = %v, want type <ServerReflectionResponse_ListServicesResponse>", r.MessageResponse)
	}
}

func (gc *GrpcClient) ListMethods() ([]*GrpcMethod, error) {
	return []*GrpcMethod{}, nil
}

// ================================== server ==================================

type GrpcServer struct {
}

// create a new grpc server
func NewGrpcServer(addr string) *GrpcServer {
	return &GrpcServer{}
}
