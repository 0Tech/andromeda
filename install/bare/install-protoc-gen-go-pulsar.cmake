#!/usr/bin/env cmake -P

cmake_policy(SET CMP0140 NEW)

function(get_version _version)
  set(${_version} 0)
  return(PROPAGATE ${_version})
endfunction()

get_version(version)
if(VERSION VERSION_EQUAL version)
  return()
endif()

execute_process(
  COMMAND go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@v${VERSION}
  COMMAND_ERROR_IS_FATAL ANY
)

message(STATUS "Installed ${VERSION}")
