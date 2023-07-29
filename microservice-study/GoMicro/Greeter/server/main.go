// Package main
package main

import (
	"context"
	hello "github.com/go-micro/examples/greeter/srv/proto/hello"
	opentracingplugins "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	"go-micro.dev/v4"
	"log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	log.Println("Received Say.Hello request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

func main() {
	// 配置 Jaeger tracer
	cfg := jaegercfg.Configuration{
		ServiceName: "greeter-service",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, _, err := cfg.NewTracer()
	if err != nil {
		log.Fatalf("Error: cannot setup Jaeger tracer: %v", err)
	}
	opentracing.SetGlobalTracer(tracer)

	service := micro.NewService(
		micro.Name("go.micro.srv.greeter"),
		micro.WrapHandler(opentracingplugins.NewHandlerWrapper(opentracing.GlobalTracer())),
	)

	// optionally setup command line usage
	service.Init()

	// Register Handlers
	hello.RegisterSayHandler(service.Server(), new(Say))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}