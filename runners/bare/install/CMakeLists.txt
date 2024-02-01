add_custom_target(install_builder_bare
  DEPENDS install_golangci-lint install_buf install_protoc-gen-gocosmos install_protoc-gen-go-pulsar install_protoc-gen-grpc-gateway
)

add_custom_target(install_go
  COMMAND cmake -D VERSION=${GO_VERSION} -P ${CMAKE_CURRENT_SOURCE_DIR}/go.cmake
  VERBATIM
)

add_custom_target(install_buf
  COMMAND cmake -D VERSION=${BUF_VERSION} -P ${CMAKE_CURRENT_SOURCE_DIR}/buf.cmake
  VERBATIM
)

add_custom_target(install_golangci-lint
  COMMAND cmake -D VERSION=${GO_LINT_VERSION} -P ${CMAKE_CURRENT_SOURCE_DIR}/golangci-lint.cmake
  VERBATIM
)

add_custom_target(install_protoc-gen-gocosmos
  COMMAND cmake -D VERSION=${PROTOC_GEN_GOCOSMOS_VERSION} -P ${CMAKE_CURRENT_SOURCE_DIR}/protoc-gen-gocosmos.cmake
  VERBATIM
)

add_custom_target(install_protoc-gen-go-pulsar
  COMMAND cmake -D VERSION=${PROTOC_GEN_GO_PULSAR_VERSION} -P ${CMAKE_CURRENT_SOURCE_DIR}/protoc-gen-go-pulsar.cmake
  VERBATIM
)

add_custom_target(install_protoc-gen-grpc-gateway
  DEPENDS install_protoc-gen-go-grpc
  COMMAND cmake -D VERSION=${PROTOC_GEN_GRPC_GATEWAY_VERSION} -P ${CMAKE_CURRENT_SOURCE_DIR}/protoc-gen-grpc-gateway.cmake
  VERBATIM
)

add_custom_target(install_protoc-gen-go-grpc
  DEPENDS install_protoc-gen-go
  COMMAND cmake -D VERSION=${PROTOC_GEN_GO_GRPC_VERSION} -P ${CMAKE_CURRENT_SOURCE_DIR}/protoc-gen-go-grpc.cmake
  VERBATIM
)

add_custom_target(install_protoc-gen-go
  COMMAND cmake -D VERSION=${PROTOC_GEN_GO_VERSION} -P ${CMAKE_CURRENT_SOURCE_DIR}/protoc-gen-go.cmake
  VERBATIM
)
