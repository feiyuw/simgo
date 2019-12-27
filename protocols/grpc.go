package protocols

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"

	"simgo/logger"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
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
	addr            string
	desc            grpcurl.DescriptorSource
	server          *grpc.Server
	handlerM        map[string]func(in *dynamic.Message, out *dynamic.Message) error
	defaultHandlerM map[string]func(in *dynamic.Message, out *dynamic.Message) error
}

// create a new grpc server
func NewGrpcServer(addr string, protos []string, opts ...grpc.ServerOption) *GrpcServer {
	desc, err := grpcurl.DescriptorSourceFromProtoFiles([]string{}, protos...)
	if err != nil {
		logger.Fatalf("protocols/grpc", "cannot parse proto file: %v", err)
	}
	return &GrpcServer{
		addr:            addr,
		desc:            desc,
		server:          grpc.NewServer(opts...),
		handlerM:        map[string]func(in *dynamic.Message, out *dynamic.Message) error{},
		defaultHandlerM: map[string]func(in *dynamic.Message, out *dynamic.Message) error{},
	}
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
					Handler:       getStreamHandler(mtd),
					ServerStreams: true,
					ClientStreams: true,
				})
			} else if mtd.IsClientStreaming() {
				streamMethods = append(streamMethods, grpc.StreamDesc{
					StreamName:    mtd.GetName(),
					Handler:       getStreamHandler(mtd),
					ClientStreams: true,
				})
			} else if mtd.IsServerStreaming() {
				streamMethods = append(streamMethods, grpc.StreamDesc{
					StreamName:    mtd.GetName(),
					Handler:       getStreamHandler(mtd),
					ServerStreams: true,
				})
			} else {
				unaryMethods = append(unaryMethods, grpc.MethodDesc{
					MethodName: mtd.GetName(),
					Handler:    gs.getUnaryHandler(mtd),
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

func (gs *GrpcServer) Stop() error {
	if gs.server != nil {
		gs.server.Stop()
		gs.server = nil
		gs.handlerM = map[string]func(in *dynamic.Message, out *dynamic.Message) error{}
		gs.defaultHandlerM = map[string]func(in *dynamic.Message, out *dynamic.Message) error{}
		logger.Infof("protocols/grpc", "grpc server %s stopped", gs.addr)
	}

	return nil
}

// TODO: stream handler support
// set specified method handler, one method only have one handler, it's the highest priority
func (gs *GrpcServer) SetMethodHandler(mtd string, handler func(in *dynamic.Message, out *dynamic.Message) error) error {
	if _, exists := gs.handlerM[mtd]; exists {
		logger.Warnf("protocols/grpc", "handler for method %s exists, will be overrided", mtd)
	}
	gs.handlerM[mtd] = handler
	return nil
}

// TODO: stream handler support
// set default method handler, one method only have one default handler, it's the lowest priority
func (gs *GrpcServer) SetDefaultMethodHandler(mtd string, handler func(in *dynamic.Message, out *dynamic.Message) error) error {
	if _, exists := gs.defaultHandlerM[mtd]; exists {
		logger.Warnf("protocols/grpc", "default handler for method %s exists, will be overrided", mtd)
	}
	gs.defaultHandlerM[mtd] = handler
	return nil
}

func (gs *GrpcServer) getMethodHandler(mtd string) (func(in *dynamic.Message, out *dynamic.Message) error, error) {
	handler, ok := gs.handlerM[mtd]
	if !ok {
		handler, ok = gs.defaultHandlerM[mtd]
		if !ok {
			return nil, fmt.Errorf("handler for method %s not found", mtd)
		}
	}
	return handler, nil
}

func (gs *GrpcServer) getUnaryHandler(mtd *desc.MethodDescriptor) func(interface{}, context.Context, func(interface{}) error, grpc.UnaryServerInterceptor) (interface{}, error) {
	return func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		handler, err := gs.getMethodHandler(mtd.GetFullyQualifiedName())
		if err != nil {
			return nil, err
		}
		in := dynamic.NewMessage(mtd.GetInputType())
		if err := dec(in); err != nil {
			return nil, err
		}
		out := dynamic.NewMessage(mtd.GetOutputType())
		if err := handler(in, out); err != nil {
			return nil, err
		}

		if interceptor == nil {
			return out, nil
		}

		info := &grpc.UnaryServerInfo{
			Server:     srv,
			FullMethod: mtd.GetFullyQualifiedName(),
		}
		wrapper := func(ctx context.Context, req interface{}) (interface{}, error) {
			return out, nil
		}
		return interceptor(ctx, in, info, wrapper)
	}
}

func getStreamHandler(mtd *desc.MethodDescriptor) func(interface{}, grpc.ServerStream) error {
	// TODO
	return func(srv interface{}, stream grpc.ServerStream) error {
		return nil
	}
}

// mock server interface for service descriptor
type mockServerIface interface {
}

// mock server struct for service descriptor
type mockServer struct {
}
