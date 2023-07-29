// Package main
package main

import (
	"context"
	hello "github.com/go-micro/examples/greeter/srv/proto/hello"
	_ "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	_ "github.com/opentracing/opentracing-go"
	"go-micro.dev/v4"
	"log"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	log.Println("Received Say.Hello request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

func main() {

	service := micro.NewService(
		micro.Name("go.micro.srv.greeter"),
		//micro.WrapHandler(opentracingplugins.NewHandlerWrapper(opentracing.GlobalTracer())),
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