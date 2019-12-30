package protocols

import (
	"context"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"net"
	"testing"
	"time"

	"simgo/logger"

	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	ecpb "google.golang.org/grpc/examples/features/proto/echo"
	hwpb "google.golang.org/grpc/examples/helloworld/helloworld"
)

func TestGrpcClient(t *testing.T) {
	s := startDemoServer(3999)
	defer s.Stop()

	Convey("proto files have higher priority than server reflection", t, func() {
		client := NewGrpcClient("127.0.0.1:3999", []string{"helloworld.proto"}, grpc.WithInsecure())
		svcs, err := client.ListServices()
		So(err, ShouldBeNil)
		So(len(svcs), ShouldEqual, 1)
		So(svcs[0], ShouldEqual, "helloworld.Greeter")
		mtds, err := client.ListMethods("helloworld.Greeter")
		So(len(mtds), ShouldEqual, 1)
		So(mtds[0], ShouldEqual, "helloworld.Greeter.SayHello")
	})

	Convey("list grpc services and methos from multi proto files", t, func() {
		client := NewGrpcClient("127.0.0.1:3999", []string{"helloworld.proto", "echo.proto"}, grpc.WithInsecure())
		svcs, err := client.ListServices()
		So(err, ShouldBeNil)
		So(svcs, ShouldResemble, []string{"grpc.examples.echo.Echo", "helloworld.Greeter"})
		mtds, err := client.ListMethods("helloworld.Greeter")
		So(err, ShouldBeNil)
		So(len(mtds), ShouldEqual, 1)
		So(mtds[0], ShouldEqual, "helloworld.Greeter.SayHello")
		mtds, err = client.ListMethods("grpc.examples.echo.Echo")
		So(err, ShouldBeNil)
		So(mtds, ShouldResemble, []string{
			"grpc.examples.echo.Echo.BidirectionalStreamingEcho",
			"grpc.examples.echo.Echo.ClientStreamingEcho",
			"grpc.examples.echo.Echo.ServerStreamingEcho",
			"grpc.examples.echo.Echo.UnaryEcho",
		})
	})

	Convey("list grpc services and methods from server reflection", t, func() {
		client := NewGrpcClient("127.0.0.1:3999", []string{}, grpc.WithInsecure())
		svcs, err := client.ListServices()
		So(err, ShouldBeNil)
		So(len(svcs), ShouldEqual, 3)
		So(svcs[0], ShouldEqual, "grpc.examples.echo.Echo")
		So(svcs[1], ShouldEqual, "grpc.reflection.v1alpha.ServerReflection")
		So(svcs[2], ShouldEqual, "helloworld.Greeter")
		mtds, err := client.ListMethods("grpc.examples.echo.Echo")
		So(mtds, ShouldResemble, []string{
			"grpc.examples.echo.Echo.BidirectionalStreamingEcho",
			"grpc.examples.echo.Echo.ClientStreamingEcho",
			"grpc.examples.echo.Echo.ServerStreamingEcho",
			"grpc.examples.echo.Echo.UnaryEcho",
		})
	})

	Convey("invoke rpc of sync method with proto files", t, func() {
		client := NewGrpcClient("127.0.0.1:3999", []string{"helloworld.proto"}, grpc.WithInsecure())
		out, err := client.InvokeRPC("helloworld.Greeter.SayHello", map[string]interface{}{"name": "you"})
		So(err, ShouldBeNil)
		So(out["message"], ShouldEqual, "Hello you")
	})

	Convey("invoke rpc of sync method without proto files", t, func() {
		client := NewGrpcClient("127.0.0.1:3999", []string{}, grpc.WithInsecure())
		out, err := client.InvokeRPC("grpc.examples.echo.Echo.UnaryEcho", map[string]interface{}{"message": "hello"})
		So(err, ShouldBeNil)
		So(out["message"], ShouldEqual, "hello")
	})

	Convey("invoke rpc of streaming method", t, func() {
		client := NewGrpcClient("127.0.0.1:3999", []string{}, grpc.WithInsecure())
		out, err := client.InvokeRPC("grpc.examples.echo.Echo.ClientStreamingEcho", map[string]interface{}{"message": "hello"})
		So(err, ShouldBeNil)
		So(out["message"], ShouldEqual, "hello")

		out, err = client.InvokeRPC("grpc.examples.echo.Echo.BidirectionalStreamingEcho", map[string]interface{}{"message": "hello"})
		So(err, ShouldBeNil)
		So(out["message"], ShouldEqual, "hello")
	})
}

func TestGrpcServer(t *testing.T) {
	s := NewGrpcServer(":4999", []string{"echo.proto", "helloworld.proto"})
	s.SetMethodHandler("grpc.examples.echo.Echo.UnaryEcho", func(in *dynamic.Message, out *dynamic.Message) error {
		out.SetFieldByName("message", "hello")
		return nil
	})
	s.SetMethodHandler("helloworld.Greeter.SayHello", func(in *dynamic.Message, out *dynamic.Message) error {
		out.SetFieldByName("message", in.GetFieldByName("name"))
		return nil
	})
	s.Start()
	defer s.Stop()
	time.Sleep(time.Microsecond)

	client := NewGrpcClient("127.0.0.1:4999", []string{"echo.proto", "helloworld.proto"}, grpc.WithInsecure())

	Convey("simulated server always return the same data", t, func() {
		out, err := client.InvokeRPC("grpc.examples.echo.Echo.UnaryEcho", map[string]interface{}{"message": "xxxx"})
		So(err, ShouldBeNil)
		So(out["message"], ShouldEqual, "hello")
	})

	Convey("echo server always return the same data", t, func() {
		out, err := client.InvokeRPC("helloworld.Greeter.SayHello", map[string]interface{}{"name": "this is a sentence"})
		So(err, ShouldBeNil)
		So(out["message"], ShouldEqual, "this is a sentence")

		out, err = client.InvokeRPC("helloworld.Greeter.SayHello", map[string]interface{}{"name": "中文：你好世界！hello world"})
		So(err, ShouldBeNil)
		So(out["message"], ShouldEqual, "中文：你好世界！hello world")
	})

	Convey("reply after 1 millisecond delay", t, func() {
		s.SetMethodHandler("helloworld.Greeter.SayHello", func(in *dynamic.Message, out *dynamic.Message) error {
			time.Sleep(time.Millisecond)
			out.SetFieldByName("message", "after sleep")
			return nil
		})
		start := time.Now().UnixNano()
		out, err := client.InvokeRPC("helloworld.Greeter.SayHello", map[string]interface{}{"name": "what to do"})
		duration := time.Now().UnixNano() - start
		So(err, ShouldBeNil)
		So(out["message"], ShouldEqual, "after sleep")
		So(duration, ShouldBeGreaterThanOrEqualTo, 1_000_000)
	})

	Convey("streaming API", t, func() {
		s.SetMethodHandler("grpc.examples.echo.Echo.BidirectionalStreamingEcho", func(in *dynamic.Message, out *dynamic.Message) error {
			out.SetFieldByName("message", "dodododo")
			return nil
		})
		out, err := client.InvokeRPC("grpc.examples.echo.Echo.BidirectionalStreamingEcho", map[string]interface{}{"message": "xxxx"})
		So(err, ShouldBeNil)
		So(out["message"], ShouldEqual, "dodododo")
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

func (s *ecServer) ClientStreamingEcho(stream ecpb.Echo_ClientStreamingEchoServer) error {
	// Read requests and send responses.
	var message string
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&ecpb.EchoResponse{Message: message})
		}
		message = in.Message
		if err != nil {
			return err
		}
	}
}

func (s *ecServer) BidirectionalStreamingEcho(stream ecpb.Echo_BidirectionalStreamingEchoServer) error {
	// Read requests and send responses.
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if err := stream.Send(&ecpb.EchoResponse{Message: in.Message}); err != nil {
			return err
		}
	}
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
