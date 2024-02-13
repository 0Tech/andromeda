#!/usr/bin/env cmake -P

cmake_policy(SET CMP0140 NEW)

include(common.cmake)

function(get_version _version)
  set(${_version} 0)

  execute_process(
	COMMAND protoc-gen-go-grpc -version
	OUTPUT_VARIABLE _output
  )
  string(REGEX MATCH "[0-9](\\.[0-9]+)*" ${_version} "${_output}")
  return(PROPAGATE ${_version})
endfunction()

require_variables(VERSION)

get_version(version)
if(VERSION VERSION_EQUAL version)
  return()
endif()

execute_process(
  COMMAND go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${VERSION}
  COMMAND_ERROR_IS_FATAL ANY
)

message(STATUS "Installed ${VERSION}")
