package protocols

import (
	"context"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net"
	"testing"
	"time"

	"simgo/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	ecpb "google.golang.org/grpc/examples/features/proto/echo"
	hwpb "google.golang.org/grpc/examples/helloworld/helloworld"
)

func TestGrpcServiceInspect(t *testing.T) {
	Convey("list grpc services from proto files", t, func() {
		s := startDemoServer(2999)
		defer s.Stop()
		client := NewGrpcClient("127.0.0.1:2999", []string{"helloworld.proto"}, grpc.WithInsecure())
		svcs, err := client.ListServices()
		So(err, ShouldBeNil)
		So(len(svcs), ShouldEqual, 1)
		So(svcs[0].name, ShouldEqual, "helloworld.Greeter")
	})

	Convey("list grpc services from server reflection", t, func() {
		s := startDemoServer(3999)
		defer s.Stop()
		client := NewGrpcClient("127.0.0.1:3999", []string{}, grpc.WithInsecure())
		svcs, err := client.ListServices()
		So(err, ShouldBeNil)
		So(len(svcs), ShouldEqual, 3)
		So(svcs[0].name, ShouldEqual, "grpc.examples.echo.Echo")
		So(svcs[1].name, ShouldEqual, "grpc.reflection.v1alpha.ServerReflection")
		So(svcs[2].name, ShouldEqual, "helloworld.Greeter")
	})
}

// hwServer is used to implement helloworld.GreeterServer.
type hwServer struct {
	hwpb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *hwServer) SayHello(ctx context.Context, in *hwpb.HelloRequest) (*hwpb.HelloReply, error) {
	return &hwpb.HelloReply{Message: "Hello " + in.Name}, nil
}

type ecServer struct {
	ecpb.UnimplementedEchoServer
}

func (s *ecServer) UnaryEcho(ctx context.Context, req *ecpb.EchoRequest) (*ecpb.EchoResponse, error) {
	return &ecpb.EchoResponse{Message: req.Message}, nil
}

func startDemoServer(port int) *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatalf("protocols/grpc", "failed to listen: %v", err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())
	s := grpc.NewServer()

	// Register Greeter on the server.
	hwpb.RegisterGreeterServer(s, &hwServer{})

	// Register RouteGuide on the same server.
	ecpb.RegisterEchoServer(s, &ecServer{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	go func() {
		if err := s.Serve(lis); err != nil {
			logger.Fatalf("protocols/grpc", "failed to serve: %v", err)
		}
	}()
	time.Sleep(time.Microsecond)

	return s
}
