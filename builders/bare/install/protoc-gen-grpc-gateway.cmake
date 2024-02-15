#!/usr/bin/env cmake -P

cmake_policy(SET CMP0140 NEW)

include(common.cmake)

require_variables(VERSION)

get_version(version protoc-gen-grpc-gateway github.com/grpc-ecosystem/grpc-gateway)
if(VERSION VERSION_EQUAL version)
  return()
endif()

execute_process(
  COMMAND go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v${VERSION}
  COMMAND_ERROR_IS_FATAL ANY
)

message(STATUS "Installed ${VERSION}")
