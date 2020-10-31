package main

import (
	"context"
	"fmt"
	hello "github.com/wanzeping72/resolver_demo/proto"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

var serviceName = "hello"

type greeterService struct {
	cli *clientv3.Client
}

func main() {
	gs := NewGreeterService()
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	hello.RegisterGreeterServer(s, gs)
	fmt.Printf("start...")
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve, ", err)
	}
}

func NewGreeterService() *greeterService {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 120 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	g := &greeterService{cli: cli}
	g.Register()
	return g
}

func (g *greeterService) Register() {
	target := fmt.Sprintf("%s", serviceName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	// 没有实现租约功能
	_, err := g.cli.Put(ctx, target, "127.0.0.1:8000")
	if err != nil {
		panic(err)
	}
}

func (g *greeterService) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloReply, error) {
	return &hello.HelloReply{Message: "hello " + req.Name}, nil
}
