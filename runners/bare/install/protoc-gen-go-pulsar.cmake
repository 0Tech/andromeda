#!/usr/bin/env cmake -P

cmake_policy(SET CMP0140 NEW)

include(version.cmake)

get_version(version protoc-gen-go-pulsar github.com/cosmos/cosmos-proto)
if(VERSION VERSION_EQUAL version)
  return()
endif()

execute_process(
  COMMAND go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@v${VERSION}
  COMMAND_ERROR_IS_FATAL ANY
)

message(STATUS "Installed ${VERSION}")
