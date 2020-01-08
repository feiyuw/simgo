simgo是一个统一的服务模拟服务，常用于契约测试。

[![Build Status](https://travis-ci.org/feiyuw/simgo.svg?branch=master)](https://travis-ci.org/feiyuw/simgo)

simgo同时支持客户端和服务端的协议模拟，当前支持的协议包括：

- [x] gRPC
- [ ] HTTP
- [ ] Dubbo

你可以将simgo作为library用于工具和单元测试，也可以直接用作测试工具。

## Example

```go
# mock server
import (
    "github.com/feiyuw/simgo"
)

var (
    ch = make(chan, bool)
)

func main() {
    s, _ := simgo.protocols.NewGrpcServer(":4999", []string{"echo.proto"})
	s.SetMethodHandler("grpc.examples.echo.Echo.UnaryEcho", func(in *dynamic.Message, out *dynamic.Message, stream grpc.ServerStream) error {
		out.SetFieldByName("message", in.GetFieldByName("message"))
		return nil
	})
	s.Start()
	<- ch
}
```


