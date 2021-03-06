package protocols

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/feiyuw/simgo/logger"

	"github.com/jhump/protoreflect/dynamic"
	"github.com/robertkrimen/otto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	ecpb "google.golang.org/grpc/examples/features/proto/echo"
	hwpb "google.golang.org/grpc/examples/helloworld/helloworld"
)

func TestGrpcClient(t *testing.T) {
	s := startDemoServer(3999)
	defer s.Stop()

	Convey("proto files have higher priority than server reflection", t, func() {
		client, _ := NewGrpcClient("127.0.0.1:3999", []string{"helloworld.proto"}, grpc.WithInsecure())
		svcs, err := client.ListServices()
		So(err, ShouldBeNil)
		So(len(svcs), ShouldEqual, 1)
		So(svcs[0], ShouldEqual, "helloworld.Greeter")
		mtds, err := client.ListMethods("helloworld.Greeter")
		So(len(mtds), ShouldEqual, 1)
		So(mtds[0], ShouldEqual, "helloworld.Greeter.SayHello")
	})

	Convey("list grpc services and methos from multi proto files", t, func() {
		client, _ := NewGrpcClient("127.0.0.1:3999", []string{"helloworld.proto", "echo.proto"}, grpc.WithInsecure())
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
		client, _ := NewGrpcClient("127.0.0.1:3999", []string{}, grpc.WithInsecure())
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
		client, _ := NewGrpcClient("127.0.0.1:3999", []string{"helloworld.proto"}, grpc.WithInsecure())
		out, err := client.InvokeRPC("helloworld.Greeter.SayHello", map[string]interface{}{"name": "you"})
		So(err, ShouldBeNil)
		So(out.(map[string]interface{})["message"], ShouldEqual, "Hello you")
	})

	Convey("invoke rpc of sync method without proto files", t, func() {
		client, _ := NewGrpcClient("127.0.0.1:3999", []string{}, grpc.WithInsecure())
		out, err := client.InvokeRPC("grpc.examples.echo.Echo.UnaryEcho", map[string]interface{}{"message": "hello"})
		So(err, ShouldBeNil)
		So(out.(map[string]interface{})["message"], ShouldEqual, "hello")
	})

	Convey("invoke rpc of streaming method", t, func() {
		client, _ := NewGrpcClient("127.0.0.1:3999", []string{}, grpc.WithInsecure())
		out, err := client.InvokeRPC("grpc.examples.echo.Echo.ClientStreamingEcho", map[string]interface{}{"message": "hello"})
		So(err, ShouldBeNil)
		So(out.(map[string]interface{})["message"], ShouldEqual, "hello")

		out, err = client.InvokeRPC("grpc.examples.echo.Echo.BidirectionalStreamingEcho", map[string]interface{}{"message": "hello"})
		So(err, ShouldBeNil)
		So(out.(map[string]interface{})["message"], ShouldEqual, "hello")
	})
}

func TestGrpcServer(t *testing.T) {
	s, _ := NewGrpcServer(":4999", []string{"echo.proto", "helloworld.proto"})
	s.SetMethodHandler("grpc.examples.echo.Echo.UnaryEcho", func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error {
		out.SetFieldByName("message", "hello")
		return nil
	})
	s.SetMethodHandler("helloworld.Greeter.SayHello", func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error {
		out.SetFieldByName("message", in.GetFieldByName("name"))
		return nil
	})
	s.Start()
	defer s.Close()
	time.Sleep(time.Millisecond) // make sure server started

	client, _ := NewGrpcClient("127.0.0.1:4999", []string{"echo.proto", "helloworld.proto"}, grpc.WithInsecure())

	Convey("simulated server always return the same data", t, func() {
		out, err := client.InvokeRPC("grpc.examples.echo.Echo.UnaryEcho", map[string]interface{}{"message": "xxxx"})
		So(err, ShouldBeNil)
		So(out.(map[string]interface{})["message"], ShouldEqual, "hello")
	})

	Convey("echo server always return the same data", t, func() {
		out, err := client.InvokeRPC("helloworld.Greeter.SayHello", map[string]interface{}{"name": "this is a sentence"})
		So(err, ShouldBeNil)
		So(out.(map[string]interface{})["message"], ShouldEqual, "this is a sentence")

		out, err = client.InvokeRPC("helloworld.Greeter.SayHello", map[string]interface{}{"name": "中文：你好世界！hello world"})
		So(err, ShouldBeNil)
		So(out.(map[string]interface{})["message"], ShouldEqual, "中文：你好世界！hello world")
	})

	Convey("change method handler", t, func() {
		s.SetMethodHandler("grpc.examples.echo.Echo.UnaryEcho", func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error {
			out.SetFieldByName("message", "world")
			return nil
		})
		defer s.SetMethodHandler("grpc.examples.echo.Echo.UnaryEcho", func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error {
			out.SetFieldByName("message", "hello")
			return nil
		})
		out, err := client.InvokeRPC("grpc.examples.echo.Echo.UnaryEcho", map[string]interface{}{"message": "xxxx"})
		So(err, ShouldBeNil)
		So(out.(map[string]interface{})["message"], ShouldEqual, "world")
	})

	Convey("use javascript handler", t, func() {
		s.SetMethodHandler("grpc.examples.echo.Echo.UnaryEcho", func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error {
			vm := otto.New()
			vm.Set("ctx", map[string]interface{}{
				"in":     in,
				"out":    out,
				"stream": stream,
			})
			vm.Run(`ctx.out.SetFieldByName("message", "javascript")`)
			return nil
		})
		defer s.SetMethodHandler("grpc.examples.echo.Echo.UnaryEcho", func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error {
			out.SetFieldByName("message", "hello")
			return nil
		})
		out, err := client.InvokeRPC("grpc.examples.echo.Echo.UnaryEcho", map[string]interface{}{"message": "zxyabc"})
		So(err, ShouldBeNil)
		So(out.(map[string]interface{})["message"], ShouldEqual, "javascript")
	})

	Convey("reply after 1 millisecond delay", t, func() {
		s.SetMethodHandler("helloworld.Greeter.SayHello", func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error {
			time.Sleep(time.Millisecond)
			out.SetFieldByName("message", "after sleep")
			return nil
		})
		start := time.Now().UnixNano()
		out, err := client.InvokeRPC("helloworld.Greeter.SayHello", map[string]interface{}{"name": "what to do"})
		duration := time.Now().UnixNano() - start
		So(err, ShouldBeNil)
		So(out.(map[string]interface{})["message"], ShouldEqual, "after sleep")
		So(duration, ShouldBeGreaterThanOrEqualTo, 1000000)
	})

	Convey("client streaming API", t, func() {
		recvedMsgs := make([]*dynamic.Message, 0)
		s.SetMethodHandler("grpc.examples.echo.Echo.ClientStreamingEcho", func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error {
			for {
				msg := dynamic.NewMessage(in.GetMessageDescriptor())
				err := stream.RecvMsg(msg)
				if err == io.EOF {
					break
				}
				recvedMsgs = append(recvedMsgs, msg)
			}
			out.SetFieldByName("message", "xxxyyy")
			stream.SendMsg(out)
			return nil
		})
		out, err := client.InvokeRPC("grpc.examples.echo.Echo.ClientStreamingEcho", []map[string]interface{}{
			map[string]interface{}{"message": "xxxx"},
			map[string]interface{}{"message": "yyyy"},
		})
		So(err, ShouldBeNil)
		So(len(recvedMsgs), ShouldEqual, 2)
		So(recvedMsgs[0].GetFieldByName("message"), ShouldEqual, "xxxx")
		So(recvedMsgs[1].GetFieldByName("message"), ShouldEqual, "yyyy")
		So(out.(map[string]interface{})["message"], ShouldEqual, "xxxyyy")
	})

	Convey("server streaming API", t, func() {
		s.SetMethodHandler("grpc.examples.echo.Echo.ServerStreamingEcho", func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error {
			stream.RecvMsg(in)
			out.SetFieldByName("message", in.GetFieldByName("message"))
			stream.SendMsg(out)
			out.SetFieldByName("message", "end")
			stream.SendMsg(out)
			return nil
		})
		out, err := client.InvokeRPC("grpc.examples.echo.Echo.ServerStreamingEcho", map[string]interface{}{"message": "xxxx"})
		So(err, ShouldBeNil)
		So(out.([]map[string]interface{})[0]["message"], ShouldEqual, "xxxx")
		So(out.([]map[string]interface{})[1]["message"], ShouldEqual, "end")
	})

	Convey("bidi streaming API", t, func() {
		s.SetMethodHandler("grpc.examples.echo.Echo.BidirectionalStreamingEcho", func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error {
			for {
				if err := stream.RecvMsg(in); err == io.EOF {
					break
				}
				out.SetFieldByName("message", in.GetFieldByName("message"))
				stream.SendMsg(out)
			}
			return nil
		})
		out, err := client.InvokeRPC("grpc.examples.echo.Echo.BidirectionalStreamingEcho", map[string]interface{}{"message": "xxxx"})
		So(err, ShouldBeNil)
		So(out.(map[string]interface{})["message"], ShouldEqual, "xxxx")

		out, err = client.InvokeRPC("grpc.examples.echo.Echo.BidirectionalStreamingEcho", []map[string]interface{}{
			map[string]interface{}{"message": "a"},
			map[string]interface{}{"message": "b"},
			map[string]interface{}{"message": "c"},
		})
		So(err, ShouldBeNil)
		outSlice := out.([]map[string]interface{})
		So(len(outSlice), ShouldEqual, 3)
		So(outSlice[0]["message"], ShouldEqual, "a")
		So(outSlice[1]["message"], ShouldEqual, "b")
		So(outSlice[2]["message"], ShouldEqual, "c")
	})

	Convey("handle listeners", t, func() {
		msgCnt := 0
		inMsgs := [][]string{}
		outMsgs := [][]string{}
		s.AddListener(func(mtd, direction, from, to, body string) error {
			msgCnt++
			return nil
		})
		s.AddListener(func(mtd, direction, from, to, body string) error {
			switch direction {
			case "in":
				inMsgs = append(inMsgs, []string{mtd, from, to, body})
			case "out":
				outMsgs = append(outMsgs, []string{mtd, from, to, body})
			}

			return nil
		})

		client.InvokeRPC("helloworld.Greeter.SayHello", map[string]interface{}{"name": "what to do"})
		So(msgCnt, ShouldEqual, 2)
		So(len(inMsgs), ShouldEqual, 1)
		So(inMsgs[0][0], ShouldEqual, "helloworld.Greeter.SayHello")
		So(inMsgs[0][2], ShouldEqual, ":4999")
		So(len(outMsgs), ShouldEqual, 1)
		So(outMsgs[0][0], ShouldEqual, "helloworld.Greeter.SayHello")
		So(outMsgs[0][1], ShouldEqual, ":4999")
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
