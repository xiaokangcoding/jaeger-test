package main

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/metadata"
	"log"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

	hello "github.com/go-micro/examples/greeter/srv/proto/hello"
	"go-micro.dev/v4"
)

func traceClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &traceWrapper{c}
	}
}

type traceWrapper struct {
	client.Client
}
func (t *traceWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}
	spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier(md))
	span := opentracing.StartSpan(req.Service()+"."+req.Endpoint(), ext.RPCServerOption(spanCtx))
	defer span.Finish()

	mdCarrier := opentracing.TextMapCarrier(md)
	if err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, mdCarrier); err != nil {
		return err
	}

	for k, v := range mdCarrier {
		md[k] = v
	}

	ctx = opentracing.ContextWithSpan(ctx, span)
	ctx = metadata.NewContext(ctx, md)

	return t.Client.Call(ctx, req, rsp, opts...)
}
/*
func (t *traceWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}
	spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier(md))
	span := opentracing.StartSpan(req.Service()+"."+req.Endpoint(), ext.RPCServerOption(spanCtx))
	defer span.Finish()

	mdCarrier := opentracing.TextMapCarrier(md)
	if err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, mdCarrier); err != nil {
		return err
	}

	ctx = opentracing.ContextWithSpan(ctx, span)
	ctx = metadata.NewContext(ctx, md)

	return t.Client.Call(ctx, req, rsp, opts...)
}
*/
func main() {

	cfg := jaegercfg.Configuration{
		ServiceName: "greeter-client-service",
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
		log.Fatal("Error: cannot setup Jaeger tracer: %v", err)
	}
	opentracing.SetGlobalTracer(tracer)

	// create a new service
	service := micro.NewService(
		micro.WrapClient(traceClientWrapper()),
	)

	// parse command line flags
	service.Init()

	// Use the generated client stub
	cl := hello.NewSayService("go.micro.srv.greeter", service.Client())

	// Make request
	rsp, err := cl.Hello(context.Background(), &hello.Request{
		Name: "John",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rsp.Msg)
}
