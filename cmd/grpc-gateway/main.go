package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/kpacha/krakend-grpc-gateway-post/gateway"
)

var (
	helloworldEndpoint = flag.String("hello_endpoint", "localhost:50051", "endpoint of GreeterServer")
	routeguideEndpoint = flag.String("route_endpoint", "localhost:50052", "endpoint of RouteGuideServer")
	port               = flag.Int("p", 8080, "port of the service")
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux, err := gateway.New(ctx, *helloworldEndpoint, *routeguideEndpoint)
	if err != nil {
		log.Printf("Setting up the gateway: %s", err.Error())
		return
	}

	srvAddr := fmt.Sprintf(":%d", *port)

	s := &http.Server{
		Addr:    srvAddr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		log.Printf("Shutting down the http server")
		if err := s.Shutdown(context.Background()); err != nil {
			log.Printf("Failed to shutdown http server: %v", err)
		}
	}()

	log.Printf("Starting listening at %s", srvAddr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("Failed to listen and serve: %v", err)
	}
}
