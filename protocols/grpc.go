package protocols

import (
	"bytes"
	"context"
	"encoding/json"

	"simgo/logger"

	"github.com/fullstorydev/grpcurl"
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
}

// create a new grpc server
func NewGrpcServer(addr string) *GrpcServer {
	return &GrpcServer{}
}
