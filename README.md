# gRPC-Resolver demo
学习 gRPC Resolver 写的demo, 主要功能:利用 etcd 实现服务发现.

## 如何使用
首先请确保你已经安装 [proto](https://developers.google.com/protocol-buffers)

### 1. 启动 etcd
etcd 可以使用 docker 启动,并且监听在 `2379` 端口上.

### 2. 编译 pb 文件
```shell
cd proto
protoc --go_out=plugins=grpc:. hello.proto
```

### 3. 启动服务端
```shell
cd server
go run .
```
启动后,服务器会自己注册到 etcd 中, 我们可以在终端中查看: 
```shell
etcdctl get "hello" --prefix
```

### 4. 启动客户端
```shell
cd client
go run .
```
然后就能看见控制台中打印出了: `get reply: hello wero`