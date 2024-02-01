set(GO_VERSION 1.21 CACHE STRING "go version")
set(GO_LINT_VERSION 1.55.2 CACHE STRING "golangci-lint version")

set(BUF_VERSION 1.29.0 CACHE STRING "buf version")

set(PROTOC_GEN_GO_VERSION 1.32 CACHE STRING "protoc-gen-go version")
set(PROTOC_GEN_GO_GRPC_VERSION 1.3.0 CACHE STRING "protoc-gen-go-grpc version")
set(PROTOC_GEN_GRPC_GATEWAY_VERSION 1.16.0 CACHE STRING "protoc-gen-grpc-gateway version")
set(PROTOC_GEN_GOCOSMOS_VERSION 1.4.11 CACHE STRING "protoc-gen-gocosmos version")
set(PROTOC_GEN_GO_PULSAR_VERSION 1.0.0-beta.3 CACHE STRING "protoc-gen-go-pulsar version")

set(USE_DOCKER TRUE CACHE BOOL "use docker to run binaries")
mark_as_advanced(USE_DOCKER)

if(USE_DOCKER)
  add_subdirectory(docker)
else()
  add_subdirectory(bare)
endif()
