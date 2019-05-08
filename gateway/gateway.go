package gateway

import (
	"context"
	"errors"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"

	_ "github.com/kpacha/krakend-grpc-post/gateway/statik"
	"github.com/kpacha/krakend-grpc-post/generated/helloworld"
	"github.com/kpacha/krakend-grpc-post/generated/routeguide"
)

func New(ctx context.Context, extra map[string]interface{}) (http.Handler, error) {
	cfg := parse(extra)
	if cfg == nil {
		return nil, errors.New("wrong config")
	}

	gw := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := helloworld.RegisterGreeterHandlerFromEndpoint(ctx, gw, cfg.helloEndpoint, opts); err != nil {
		return nil, err
	}

	if err := routeguide.RegisterRouteGuideHandlerFromEndpoint(ctx, gw, cfg.routeEndpoint, opts); err != nil {
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

func parse(extra map[string]interface{}) *opts {
	name, ok := extra["name"].(string)
	if !ok {
		return nil
	}

	rawEs, ok := extra["endpoints"]
	if !ok {
		return nil
	}
	es, ok := rawEs.([]interface{})
	if !ok || len(es) < 2 {
		return nil
	}
	endpoints := make([]string, len(es))
	for i, e := range es {
		endpoints[i] = e.(string)
	}

	return &opts{
		name:          name,
		helloEndpoint: endpoints[0],
		routeEndpoint: endpoints[1],
	}
}

type opts struct {
	name          string
	helloEndpoint string
	routeEndpoint string
}
