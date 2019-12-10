package protocols

import (
	"context"
	"fmt"
	"time"

	"simgo/logger"

	"github.com/fullstorydev/grpcurl"
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

// one endpoint contain multi services
type GrpcService struct {
	name string // service name
}

// one service contain multi methods
type GrpcMethod struct {
	name string
}

// Create a new grpc client
func NewGrpcClient(addr string, protos []string, opts ...grpc.DialOption) *GrpcClient {
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		logger.Fatalf("protocols/grpc", "did not connect: %v", err)
	}

	var descSource grpcurl.DescriptorSource = nil

	if len(protos) > 0 {
		descSource, err = grpcurl.DescriptorSourceFromProtoFiles([]string{}, protos...)
		if err != nil {
			logger.Fatalf("protocols/grpc", "cannot parse proto file: %v", err)
		}
	}
	return &GrpcClient{addr: addr, conn: conn, desc: descSource}
}

func (gc *GrpcClient) ListServices() ([]*GrpcService, error) {
	if gc.desc != nil { // from protos
		return gc.listServicesFromProtos()
	}
	return gc.listServicesFromReflection()
}

func (gc *GrpcClient) listServicesFromProtos() ([]*GrpcService, error) {
	svcs, err := grpcurl.ListServices(gc.desc)
	if err != nil {
		return nil, err
	}
	grpcServices := make([]*GrpcService, len(svcs))
	for idx, svc := range svcs {
		grpcServices[idx] = &GrpcService{name: svc}
	}
	return grpcServices, nil
}

func (gc *GrpcClient) listServicesFromReflection() ([]*GrpcService, error) {
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
