package main

import (
	"context"
	"fmt"
	hello "github.com/wanzeping72/resolver_demo/proto"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"sync"
	"time"
)

var scheme = "test"
var serviceName = "hello"

func main() {
	target := fmt.Sprintf("%s://%s/%s", scheme, "127.0.0.1:2379", serviceName)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	clientConn, err := grpc.DialContext(ctx, target, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	reply, err := hello.NewGreeterClient(clientConn).SayHello(context.Background(), &hello.HelloRequest{
		Name: "wero",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("get reply:", reply.Message)
}

func init() {
	resolver.Register(newBuilder())
}

var once sync.Once

type exampleBuilder struct {
}

func newBuilder() *exampleBuilder {
	return &exampleBuilder{}
}

func (b *exampleBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &exampleResolver{cc: cc, target: target}
	r.resolve()
	return r, nil
}

func (b *exampleBuilder) Scheme() string {
	return scheme
}

type exampleResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
	cli    *clientv3.Client
}

func (e *exampleResolver) ResolveNow(_ resolver.ResolveNowOptions) {
	fmt.Println("ResolverNow")
}

func (e *exampleResolver) Close() {
	fmt.Println("Close")
}

func (e *exampleResolver) resolve() {
	once.Do(func() {
		clientConfig := clientv3.Config{
			Endpoints:   []string{e.target.Authority},
			DialTimeout: 120 * time.Second,
		}
		cli, err := clientv3.New(clientConfig)
		if err != nil {
			panic(err)
		}
		e.cli = cli
	})
	addList := make([]resolver.Address, 0)
	if result, err := e.cli.Get(context.Background(), e.target.Endpoint, clientv3.WithPrefix()); err == nil {
		for _, kvs := range result.Kvs {
			addList = append(addList, resolver.Address{
				Addr: string(kvs.Value),
			})
			fmt.Printf("resolved addr: %s\n", string(kvs.Value))
		}
	}
	e.cc.UpdateState(resolver.State{Addresses: addList})
}
