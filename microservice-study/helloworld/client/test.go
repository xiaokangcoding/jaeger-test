package main

import (
	"context"
	"fmt"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"time"

	pb "github.com/Henry-jk/jaeger-test/microservice-study/helloworld/proto"
)

func loggingUnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		fmt.Printf("Metadata received from server: %v\n", md)
	}
	return err
}

func main() {
	// Jaeger tracer 初始化
	cfg, _ := config.FromEnv()
	tracer, closer, _ := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	conn, err := grpc.Dial(":8888",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_opentracing.UnaryClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.SayHello(ctx, &pb.HelloRequest{})
	if err != nil {
		log.Fatalf("could not call: %v", err)
	}
	// 处理返回的结果
	// ...
}
/*
func main2() {
	// 初始化 Jaeger tracer
	cfg, _ := config.FromEnv()
	tracer, closer, _ := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// 连接到 gRPC 服务器
	conn, err := grpc.Dial("localhost:8888", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// 调用方法
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 埋点
	span, ctx := opentracing.StartSpanFromContext(ctx, "rpc")
	defer span.Finish()

	_, err = c.SayHello(ctx, &pb.HelloRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
}
*/