simgo is a simple service simulator for both client and server simulation. It can be used as a go unit test module, a web based test tool, or a test automation library(planned).

[![Build Status](https://travis-ci.org/feiyuw/simgo.svg?branch=master)](https://travis-ci.org/feiyuw/simgo)

Supported protocols:

| protocol | client simulator | server simulator |
| -------- | ---------------- | ---------------- |
| gRPC     |    √             |     √            |
| HTTP     |    ×             |     ×            |
| Dubbo    |    ×             |     ×            |

## Used in go unit test

status: **done**

```go
import (
	"time"
	"testing"

	"github.com/feiyuw/simgo/protocols"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAsMockClient(t *testing.T) {
	Convey("test helloworld grpc service", t, func() {
		client, err := protocols.NewGrpcClient("127.0.0.1:3999", []string{"helloworld.proto"}, grpc.WithInsecure())
		So(err, ShouldBeNil)
		svcs, err := client.ListServices()
		So(err, ShouldBeNil)
		So(len(svcs), ShouldEqual, 1)
		So(svcs[0], ShouldEqual, "helloworld.Greeter")
		mtds, err := client.ListMethods("helloworld.Greeter")
		So(len(mtds), ShouldEqual, 1)
		So(mtds[0], ShouldEqual, "helloworld.Greeter.SayHello")
		out, err := client.InvokeRPC("helloworld.Greeter.SayHello", map[string]interface{}{"name": "you"})
		So(err, ShouldBeNil)
		So(out.(map[string]interface{})["message"], ShouldEqual, "Hello you")
	})
}

func TestAsMockServer(t *testing.T) {
	s, _ := protocols.NewGrpcServer(":4999", []string{"echo.proto"})
	s.SetMethodHandler("grpc.examples.echo.Echo.UnaryEcho", func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error {
		out.SetFieldByName("message", "hello")
		return nil
	})
	s.Start()
	defer s.Close()
	time.Sleep(time.Millisecond) // make sure server started

	Convey("use dynamic field in handler", t, func() {
		s.SetMethodHandler("grpc.examples.echo.Echo.UnaryEcho", func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error {
			out.SetFieldByName("message", in.GetFieldByName("message"))
			return nil
		})
		// your test code ...
	})

	Convey("use javascript handler", t, func() {
		s.SetMethodHandler("grpc.examples.echo.Echo.UnaryEcho", func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error {
			vm := otto.New()
			vm.Set("ctx", map[string]interface{}{
				"in":     in,
				"out":    out,
				"stream": stream,
			})
			vm.Run(`ctx.out.SetFieldByName("message", ctx.in.GetFieldByName("message"))`)
			return nil
		})
		// your test code ...
	})

	Convey("delay reply", t, func() {
		s.SetMethodHandler("grpc.examples.echo.Echo.UnaryEcho", func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error {
			time.Sleep(time.Second)
			out.SetFieldByName("message", "delay 1 sec")
			return nil
		})
		// your test code ...
	})

	Convey("streaming API", t, func() {
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
		// your test code ...
	})
}
```

## As a web based test tool

status: **draft**

`simgo` can be used as a web based RPC test tool. Download the latest [release](https://github.com/feiyuw/simgo/releases).

As a client simulator.

![client](https://github.com/feiyuw/simgo/raw/master/snapshot_client.png)

As a server simulator.

![server](https://github.com/feiyuw/simgo/raw/master/snapshot_server.png)

### Handler examples

1. static response

	type: raw

	content: {"message": "hello world"}

1. dynamic response

	type: javascript

	content: 

	```javascript
		ctx.out.SetFieldByName("message", "hello " + ctx.in.GetFieldByName("message"))
	```

1. delay response

	type: javascript

	content: 

	```javascript
		ctx.Sleep(1) // 1 second
		ctx.out.SetFieldByName("message", ctx.in.GetFieldByName("message"))
	```

1. streaming response

	type: javascript

	content: 

	```javascript
		ctx.stream.RecvMsg(ctx.in)  // read input stream
		ctx.out.SetFieldByName("message", ctx.in.GetFieldByName("message"))
		ctx.stream.SendMsg(ctx.out)  // send output stream
	```

1. error response

## As a test automation library

status: **planned**
