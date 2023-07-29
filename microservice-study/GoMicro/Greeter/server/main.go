// Package main
package main

import (
	"context"
	"log"
	"time"

	hello "github.com/go-micro/examples/greeter/srv/proto/hello"
	"go-micro.dev/v4"
	"google.golang.org/grpc"
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