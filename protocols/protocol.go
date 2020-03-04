package protocols

import (
	"errors"
	"google.golang.org/grpc"
	"strconv"
	"time"
)

type RpcClient interface {
	InvokeRPC(string, interface{}) (interface{}, error)
	Close() error
}

func NewRpcClient(protocol string, server string, options map[string]interface{}) (RpcClient, error) {
	switch protocol {
	case "grpc":
		protos, exists := options["protos"]
		if !exists {
			return nil, errors.New("no protos specified")
		}
		protosStr := make([]string, len(protos.([]interface{})))
		for idx, proto := range protos.([]interface{}) {
			protosStr[idx] = proto.(string)
		}
		return NewGrpcClient(server, protosStr, grpc.WithBlock(), grpc.WithInsecure(), grpc.WithTimeout(2*time.Second)) // TODO: add SSL support
	default:
		return nil, errors.New("unsupported protocol: " + protocol)
	}
}

type RpcServer interface {
	Start() error
	Close() error
	AddListener(func(mtd, direction, from, to, body string) error)
}

func NewRpcServer(protocol string, name string, port int, options map[string]interface{}) (RpcServer, error) {
	switch protocol {
	case "grpc":
		protos, exists := options["protos"]
		if !exists {
			return nil, errors.New("no protos specified")
		}
		protosStr := make([]string, len(protos.([]interface{})))
		for idx, proto := range protos.([]interface{}) {
			protosStr[idx] = proto.(string)
		}
		return NewGrpcServer(":"+strconv.Itoa(port), protosStr)
	default:
		return nil, errors.New("unsupported protocol: " + protocol)
	}
}
