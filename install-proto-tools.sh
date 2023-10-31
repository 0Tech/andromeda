#!/bin/sh

go install github.com/bufbuild/buf/cmd/buf@latest
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@latest
go install github.com/cosmos/gogoproto/protoc-gen-gocosmos@latest
go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
