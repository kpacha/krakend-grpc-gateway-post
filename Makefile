.PHONY: generate build

all: generate build

generate:
	@sh generate.sh helloworld routeguide

build:
	go build ./cmd/grpc-gateway
	go build ./cmd/helloworld
	go build ./cmd/routeguide

build-plugin:
	go build -buildmode=plugin -o grpc-gateway-post.so ./plugin

run:
	./routeguide &
	./helloworld &
	./grpc-gateway &

stop:
	killall routeguide
	killall helloworld
	killall grpc-gateway