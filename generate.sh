#!/bin/bash

for i in "$@"
do
   : 
    echo "generating ${i} service"
    
    echo "  - grpc bindings"
	protoc -I/usr/local/include -I. \
	  -I$GOPATH/src \
	  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --go_out=plugins=grpc:. \
	  protos/${i}.proto
    
    echo "  - grpc-gateway"
	protoc -I/usr/local/include -I. \
	  -I$GOPATH/src \
	  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --grpc-gateway_out=logtostderr=true,grpc_api_configuration=protos/${i}.yml:. \
	  protos/${i}.proto

	# move generated sources
	mkdir -p generated/${i}
	mv protos/${i}.*.go generated/${i}/.
    
    echo "  - grpc-gateway swagger"
	protoc -I/usr/local/include -I. \
	  -I$GOPATH/src \
	  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --swagger_out=logtostderr=true,grpc_api_configuration=protos/${i}.yml:. \
	  protos/${i}.proto

done

echo "generating static files"
# move and pack the swagger definitions
mkdir -p protos/swagger
mv protos/*.swagger.json protos/swagger/.

[ ! -d statik ] || rm -rf statik
[ ! -d gateway/statik ] || rm -rf gateway/statik
go get github.com/rakyll/statik
go install github.com/rakyll/statik
statik -src=protos/swagger
mv statik gateway/.
rm -rf protos/swagger

echo "services generated"