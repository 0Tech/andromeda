project(install
  LANGUAGES NONE
)

set(GO_VERSION 1.21.6 CACHE STRING "go version")
set(GO_LINT_VERSION 1.55.2 CACHE STRING "golangci-lint version")

set(BUF_VERSION 1.29.0 CACHE STRING "buf version")

set(PROTOC_GEN_GO_VERSION 1.32 CACHE STRING "protoc-gen-go version")
set(PROTOC_GEN_GO_GRPC_VERSION 1.3.0 CACHE STRING "protoc-gen-go-grpc version")
set(PROTOC_GEN_GRPC_GATEWAY_VERSION 1.16.0 CACHE STRING "protoc-gen-grpc-gateway version")
set(PROTOC_GEN_GOCOSMOS_VERSION 1.4.11 CACHE STRING "protoc-gen-gocosmos version")
set(PROTOC_GEN_GO_PULSAR_VERSION 1.0.0-beta.3 CACHE STRING "protoc-gen-go-pulsar version")

add_custom_target(install_builders
  DEPENDS install_go install_golangci-lint install_buf install_protoc-gen-gocosmos install_protoc-gen-go-pulsar install_protoc-gen-grpc-gateway
)

add_custom_target(install_go
  COMMAND cmake -D VERSION=${GO_VERSION} -P ${CMAKE_CURRENT_SOURCE_DIR}/go.cmake
  VERBATIM
)

add_custom_target(install_buf
  DEPENDS install_go
  COMMAND cmake -D VERSION=${BUF_VERSION} -P ${CMAKE_CURRENT_SOURCE_DIR}/buf.cmake
  VERBATIM
)

add_custom_target(install_golangci-lint
  DEPENDS install_go
  COMMAND cmake -D VERSION=${GO_LINT_VERSION} -P ${CMAKE_CURRENT_SOURCE_DIR}/golangci-lint.cmake
  VERBATIM
)

add_custom_target(install_protoc-gen-gocosmos
  DEPENDS install_go
  COMMAND cmake -D VERSION=${PROTOC_GEN_GOCOSMOS_VERSION} -P ${CMAKE_CURRENT_SOURCE_DIR}/protoc-gen-gocosmos.cmake
  VERBATIM
)

add_custom_target(install_protoc-gen-go-pulsar
  DEPENDS install_go
  COMMAND cmake -D VERSION=${PROTOC_GEN_GO_PULSAR_VERSION} -P ${CMAKE_CURRENT_SOURCE_DIR}/protoc-gen-go-pulsar.cmake
  VERBATIM
)

add_custom_target(install_protoc-gen-grpc-gateway
  DEPENDS install_go install_protoc-gen-go-grpc
  COMMAND cmake -D VERSION=${PROTOC_GEN_GRPC_GATEWAY_VERSION} -P ${CMAKE_CURRENT_SOURCE_DIR}/protoc-gen-grpc-gateway.cmake
  VERBATIM
)

add_custom_target(install_protoc-gen-go-grpc
  DEPENDS install_go install_protoc-gen-go
  COMMAND cmake -D VERSION=${PROTOC_GEN_GO_GRPC_VERSION} -P ${CMAKE_CURRENT_SOURCE_DIR}/protoc-gen-go-grpc.cmake
  VERBATIM
)

add_custom_target(install_protoc-gen-go
  DEPENDS install_go
  COMMAND cmake -D VERSION=${PROTOC_GEN_GO_VERSION} -P ${CMAKE_CURRENT_SOURCE_DIR}/protoc-gen-go.cmake
  VERBATIM
)
