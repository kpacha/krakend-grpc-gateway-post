package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/kpacha/krakend-grpc-post/gateway"
)

func init() {
	fmt.Println("krakend-grpc-post plugin loaded!!!")
}

var GRPCRegisterer = registerer("grpc-post")

type registerer string

func (r registerer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
)) {
	f(string(r), func(ctx context.Context, extra map[string]interface{}) (http.Handler, error) {
		cfg := parse(extra)
		if cfg == nil {
			return nil, errors.New("wrong config")
		}
		if cfg.name != string(r) {
			return nil, fmt.Errorf("unknown register %s", cfg.name)
		}
		return gateway.New(ctx, cfg.helloEndpoint, cfg.routeEndpoint)
	})
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
