package protocols

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"reflect"

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
func NewGrpcClient(addr string, protos []string, opts ...grpc.DialOption) (*GrpcClient, error) {
	var descSource grpcurl.DescriptorSource

	if addr == "" {
		return nil, fmt.Errorf("addr should not be empty")
	}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("did not connect: %v", err)
	}

	if len(protos) > 0 {
		descSource, err := grpcurl.DescriptorSourceFromProtoFiles([]string{}, protos...)
		if err != nil {
			return nil, fmt.Errorf("cannot parse proto file: %v", err)
		}
		return &GrpcClient{addr: addr, conn: conn, desc: descSource}, nil
	}

	// fetch from server reflection RPC
	c := rpb.NewServerReflectionClient(conn)
	refClient := grpcreflect.NewClient(ctx, c)
	descSource = grpcurl.DescriptorSourceFromServer(ctx, refClient)

	return &GrpcClient{addr: addr, conn: conn, desc: descSource}, nil
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

func (gc *GrpcClient) Close() error {
	if gc.conn == nil {
		return nil
	}
	return gc.conn.Close()
}

type rpcResponse struct {
	messages []bytes.Buffer
}

func (rr *rpcResponse) Write(p []byte) (int, error) {
	var msg bytes.Buffer
	n, err := msg.Write(p)
	rr.messages = append(rr.messages, msg)
	return n, err
}

func (rr *rpcResponse) ToJSON() (interface{}, error) {
	switch len(rr.messages) {
	case 0:
		return map[string]interface{}{}, nil
	case 1:
		resp := make(map[string]interface{})
		if err := json.Unmarshal(rr.messages[0].Bytes(), &resp); err != nil {
			return nil, err
		}

		return resp, nil
	default:
		resp := make([]map[string]interface{}, len(rr.messages))
		for idx, msg := range rr.messages {
			oneResp := make(map[string]interface{})
			if err := json.Unmarshal(msg.Bytes(), &oneResp); err != nil {
				return nil, err
			}
			resp[idx] = oneResp
		}
		return resp, nil
	}
}

func (gc *GrpcClient) InvokeRPC(mtdName string, reqData interface{}) (interface{}, error) {
	var in bytes.Buffer
	var out = rpcResponse{messages: []bytes.Buffer{}}

	switch reflect.TypeOf(reqData).Kind() {
	case reflect.Slice:
		for _, data := range reqData.([]map[string]interface{}) {
			reqBytes, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}
			in.Write(reqBytes)
		}
	case reflect.Map:
		reqBytes, err := json.Marshal(reqData)
		if err != nil {
			return nil, err
		}
		in.Write(reqBytes)
	default:
		in.WriteString(reqData.(string))
	}

	rf, formatter, err := grpcurl.RequestParserAndFormatterFor(grpcurl.FormatJSON, gc.desc, true, false, &in)
	if err != nil {
		return nil, err
	}
	h := grpcurl.NewDefaultEventHandler(&out, gc.desc, formatter, false)
	if err = grpcurl.InvokeRPC(ctx, gc.desc, gc.conn, mtdName, []string{}, h, rf.Next); err != nil {
		return nil, err
	}

	return out.ToJSON()
}

// ================================== server ==================================

type GrpcServer struct {
	addr     string
	desc     grpcurl.DescriptorSource
	server   *grpc.Server
	handlerM map[string]func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error
}

// create a new grpc server
func NewGrpcServer(addr string, protos []string, opts ...grpc.ServerOption) (*GrpcServer, error) {
	descFromProto, err := grpcurl.DescriptorSourceFromProtoFiles([]string{}, protos...)
	if err != nil {
		return nil, fmt.Errorf("cannot parse proto file: %v", err)
	}
	gs := &GrpcServer{
		addr:     addr,
		desc:     descFromProto,
		server:   grpc.NewServer(opts...),
		handlerM: map[string]func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error{},
	}

	gs.server = grpc.NewServer()
	services, err := grpcurl.ListServices(gs.desc)
	if err != nil {
		return nil, fmt.Errorf("failed to list services")
	}
	for _, svcName := range services {
		dsc, err := gs.desc.FindSymbol(svcName)
		if err != nil {
			return nil, fmt.Errorf("unable to find service: %s, error: %v", svcName, err)
		}
		sd := dsc.(*desc.ServiceDescriptor)

		unaryMethods := []grpc.MethodDesc{}
		streamMethods := []grpc.StreamDesc{}
		for _, mtd := range sd.GetMethods() {
			logger.Debugf("protocols/grpc", "try to add method: %v of service: %s", mtd, svcName)

			if mtd.IsClientStreaming() || mtd.IsServerStreaming() {
				streamMethods = append(streamMethods, grpc.StreamDesc{
					StreamName:    mtd.GetName(),
					Handler:       gs.getStreamHandler(mtd),
					ServerStreams: mtd.IsServerStreaming(),
					ClientStreams: mtd.IsClientStreaming(),
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
			HandlerType: (*interface{})(nil),
			Methods:     unaryMethods,
			Streams:     streamMethods,
			Metadata:    sd.GetFile().GetName(),
		}
		gs.server.RegisterService(&svcDesc, &mockServer{})
	}

	return gs, nil
}

func (gs *GrpcServer) Start() error {
	lis, err := net.Listen("tcp", gs.addr)
	if err != nil {
		return err
	}
	logger.Infof("protocols/grpc", "server listening at %v", lis.Addr())

	go func() {
		if err := gs.server.Serve(lis); err != nil {
			logger.Errorf("protocols/grpc", "failed to serve: %v", err)
		}
	}()

	return nil
}

func (gs *GrpcServer) Close() error {
	if gs.server != nil {
		gs.server.Stop()
		gs.server = nil
		gs.handlerM = map[string]func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error{}
		logger.Infof("protocols/grpc", "grpc server %s stopped", gs.addr)
	}

	return nil
}

// set specified method handler, one method only have one handler, it's the highest priority
// if you want to return error, see https://github.com/avinassh/grpc-errors/blob/master/go/server.go
func (gs *GrpcServer) SetMethodHandler(mtd string, handler func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error) error {
	if _, exists := gs.handlerM[mtd]; exists {
		logger.Warnf("protocols/grpc", "handler for method %s exists, will be overrided", mtd)
	}
	gs.handlerM[mtd] = handler
	return nil
}

func (gs *GrpcServer) getMethodHandler(mtd string) (func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error, error) {
	handler, ok := gs.handlerM[mtd]
	if !ok {
		return nil, fmt.Errorf("handler for method %s not found", mtd)
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
		if err := handler(in, out, nil); err != nil {
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

func (gs *GrpcServer) getStreamHandler(mtd *desc.MethodDescriptor) func(interface{}, grpc.ServerStream) error {
	return func(srv interface{}, stream grpc.ServerStream) error {
		handler, err := gs.getMethodHandler(mtd.GetFullyQualifiedName())
		if err != nil {
			return err
		}
		in := dynamic.NewMessage(mtd.GetInputType())
		out := dynamic.NewMessage(mtd.GetOutputType())
		if err := handler(in, out, stream); err != nil {
			return err
		}
		return nil
	}
}

// mock server struct for service descriptor
type mockServer struct {
}
