package main

import (
	"context"
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
	f(string(r), gateway.New)
}
