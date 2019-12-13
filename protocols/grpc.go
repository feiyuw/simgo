package protocols

import (
	"bytes"
	"context"
	"encoding/json"
	"net"

	"simgo/logger"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	hwpb "google.golang.org/grpc/examples/helloworld/helloworld"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// ================================== client ==================================
var (
	ctx = context.Background() // TODO: add timeout, dialtime options
)

// client
type GrpcClient struct {
	addr string           // connected service addr, eg. 127.0.0.1:1988
	conn *grpc.ClientConn // connection with grpc server
	desc grpcurl.DescriptorSource
}

// Create a new grpc client
// if protos set, will get services and methods from proto files
// if addr set but protos empty, will get services and methods from server reflection
func NewGrpcClient(addr string, protos []string, opts ...grpc.DialOption) *GrpcClient {
	var descSource grpcurl.DescriptorSource

	if addr == "" {
		logger.Fatal("protocols/grpc", "addr should not be empty")
	}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		logger.Fatalf("protocols/grpc", "did not connect: %v", err)
	}

	if len(protos) > 0 {
		descSource, err := grpcurl.DescriptorSourceFromProtoFiles([]string{}, protos...)
		if err != nil {
			logger.Fatalf("protocols/grpc", "cannot parse proto file: %v", err)
		}
		return &GrpcClient{addr: addr, conn: conn, desc: descSource}
	}

	// fetch from server reflection RPC
	c := rpb.NewServerReflectionClient(conn)
	refClient := grpcreflect.NewClient(ctx, c)
	descSource = grpcurl.DescriptorSourceFromServer(ctx, refClient)

	return &GrpcClient{addr: addr, conn: conn, desc: descSource}
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

func (gc *GrpcClient) InvokeRPC(mtdName string, reqData map[string]interface{}) (map[string]interface{}, error) {
	// TODO: in change to io.Reader type
	// TODO: out change to io.Writer type
	var in, out bytes.Buffer

	reqBytes, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}
	in.Write(reqBytes)

	rf, formatter, err := grpcurl.RequestParserAndFormatterFor(grpcurl.FormatJSON, gc.desc, true, false, &in)
	if err != nil {
		return nil, err
	}
	h := grpcurl.NewDefaultEventHandler(&out, gc.desc, formatter, false)
	if err = grpcurl.InvokeRPC(ctx, gc.desc, gc.conn, mtdName, []string{}, h, rf.Next); err != nil {
		return nil, err
	}

	resp := make(map[string]interface{})
	if err = json.Unmarshal(out.Bytes(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// ================================== server ==================================

type GrpcServer struct {
	addr   string
	desc   grpcurl.DescriptorSource
	server *grpc.Server
}

// create a new grpc server
func NewGrpcServer(addr string, protos []string, opts ...grpc.ServerOption) *GrpcServer {
	desc, err := grpcurl.DescriptorSourceFromProtoFiles([]string{}, protos...)
	if err != nil {
		logger.Fatalf("protocols/grpc", "cannot parse proto file: %v", err)
	}
	return &GrpcServer{addr: addr, desc: desc, server: grpc.NewServer(opts...)}
}

func (gs *GrpcServer) Start() error {
	lis, err := net.Listen("tcp", gs.addr)
	if err != nil {
		return err
	}
	logger.Infof("protocols/grpc", "server listening at %v", lis.Addr())
	gs.server = grpc.NewServer()
	services, err := grpcurl.ListServices(gs.desc)
	if err != nil {
		logger.Error("protocols/grpc", "failed to start server")
		return err
	}
	for _, svcName := range services {
		dsc, err := gs.desc.FindSymbol(svcName)
		if err != nil {
			return err
		}
		sd := dsc.(*desc.ServiceDescriptor)

		unaryMethods := []grpc.MethodDesc{}
		streamMethods := []grpc.StreamDesc{}
		for _, mtd := range sd.GetMethods() {
			logger.Debugf("protocols/grpc", "try to add method: %v of service: %s", mtd, svcName)
			if mtd.IsClientStreaming() && mtd.IsServerStreaming() {
				streamMethods = append(streamMethods, grpc.StreamDesc{
					StreamName:    mtd.GetName(),
					Handler:       bidiStreamHandler,
					ServerStreams: true,
					ClientStreams: true,
				})
			} else if mtd.IsClientStreaming() {
				streamMethods = append(streamMethods, grpc.StreamDesc{
					StreamName:    mtd.GetName(),
					Handler:       clientStreamHandler,
					ClientStreams: true,
				})
			} else if mtd.IsServerStreaming() {
				streamMethods = append(streamMethods, grpc.StreamDesc{
					StreamName:    mtd.GetName(),
					Handler:       serverStreamHandler,
					ServerStreams: true,
				})
			} else {
				unaryMethods = append(unaryMethods, grpc.MethodDesc{
					MethodName: mtd.GetName(),
					Handler:    getUnaryHandler(mtd),
				})
			}
		}

		svcDesc := grpc.ServiceDesc{
			ServiceName: svcName,
			HandlerType: (*mockServerIface)(nil),
			Methods:     unaryMethods,
			Streams:     streamMethods,
			Metadata:    sd.GetFile().GetName(),
		}
		gs.server.RegisterService(&svcDesc, &mockServer{})
	}

	go func() {
		if err := gs.server.Serve(lis); err != nil {
			logger.Errorf("protocols/grpc", "failed to serve: %v", err)
		}
	}()

	return nil
}

func listMethods(source grpcurl.DescriptorSource, serviceName string) ([]*desc.MethodDescriptor, error) {
	dsc, err := source.FindSymbol(serviceName)
	if err != nil {
		return nil, err
	}
	sd := dsc.(*desc.ServiceDescriptor)

	return sd.GetMethods(), nil
}

func getUnaryHandler(mtd *desc.MethodDescriptor) func(interface{}, context.Context, func(interface{}) error, grpc.UnaryServerInterceptor) (interface{}, error) {
	return func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		in := new(hwpb.HelloRequest)
		if err := dec(in); err != nil {
			return nil, err
		}
		if interceptor == nil {
			return SayHello(ctx, in)
		}
		info := &grpc.UnaryServerInfo{
			Server:     srv,
			FullMethod: mtd.GetFullyQualifiedName(),
		}
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return SayHello(ctx, req.(*hwpb.HelloRequest))
		}
		return interceptor(ctx, in, info, handler)
	}
}

func clientStreamHandler(srv interface{}, stream grpc.ServerStream) error {
	// TODO
	return nil
}

func serverStreamHandler(srv interface{}, stream grpc.ServerStream) error {
	// TODO
	return nil
}

func bidiStreamHandler(srv interface{}, stream grpc.ServerStream) error {
	// TODO
	return nil
}

// mock server interface for service descriptor
type mockServerIface interface {
}

// mock server struct for service descriptor
type mockServer struct {
}

func SayHello(ctx context.Context, in *hwpb.HelloRequest) (*hwpb.HelloReply, error) {
	return &hwpb.HelloReply{Message: "Hello " + in.Name}, nil
}
