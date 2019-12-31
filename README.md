simgo是一个统一的服务模拟服务，常用于契约测试。
[![Build Status](https://travis-ci.org/feiyuw/simgo.svg?branch=master)](https://travis-ci.org/feiyuw/simgo)

## Client Simulator

### Scenarios

- [x] 知道grpc服务的IP:PORT，我可以知道这个服务暴露的rpc方法有哪些
- [x] 能向某个grpc服务发送请求，并得到正确的返回
- [x] 当grpc服务出错时，client能得到错误信息

### Client Architecture

[TODO]

## Server Simulator

### Scenarios

- [x] 基于一个服务的proto文件，我可以模拟一个服务接口，让它总是返回确定的内容
- [x] 基于一个服务的proto文件，我可以模拟一个服务接口，让它针对不同的数据返回不同的内容
- [x] 基于一个服务的proto文件，我可以模拟一个服务接口，让它延时特定的时间再返回内容

### Client Architecture

[TODO]

