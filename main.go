package main

import (
	"simgo/protocols"
)

// client: grpcurl -plaintext -proto ../../private/simgo/protocols/helloworld.proto -d '{"name": "world"}' 127.0.0.1:1777 helloworld.Greeter/SayHello
func main() {
	var ch chan bool

	server := protocols.NewGrpcServer(":1777", []string{"protocols/helloworld.proto"})
	server.Start()

	<-ch
}
