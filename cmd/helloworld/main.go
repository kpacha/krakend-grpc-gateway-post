package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/kpacha/krakend-grpc-post/generated/helloworld"
	grpc "google.golang.org/grpc"
)

func main() {
	port := flag.Int("p", 50051, "port of the service")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	helloworld.RegisterGreeterServer(grpcServer, &server{})

	grpcServer.Serve(lis)
}

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}
