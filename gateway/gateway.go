package gateway

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"

	_ "github.com/kpacha/krakend-grpc-post/gateway/statik"
	"github.com/kpacha/krakend-grpc-post/generated/helloworld"
	"github.com/kpacha/krakend-grpc-post/generated/routeguide"
)

func New(ctx context.Context, helloEndpoint, routeEndpoint string) (http.Handler, error) {
	gw := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := helloworld.RegisterGreeterHandlerFromEndpoint(ctx, gw, helloEndpoint, opts); err != nil {
		return nil, err
	}

	if err := routeguide.RegisterRouteGuideHandlerFromEndpoint(ctx, gw, routeEndpoint, opts); err != nil {
		return nil, err
	}

	statikFS, err := fs.New()
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(statikFS)))
	mux.Handle("/", gw)

	return mux, nil
}
