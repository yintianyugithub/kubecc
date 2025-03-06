.PHONY: build clean run

build:
	goctl rpc protoc add.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=. -style go_zero
