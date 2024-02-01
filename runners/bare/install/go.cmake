#!/usr/bin/env cmake -P

cmake_policy(SET CMP0140 NEW)

function(get_version _version)
  set(${_version} 0)

  execute_process(
	COMMAND go version
	OUTPUT_VARIABLE _output
  )
  string(REGEX MATCH "[0-9]\\.[0-9]+" ${_version} "${_output}")
  return(PROPAGATE ${_version})
endfunction()

get_version(version)
if(VERSION VERSION_EQUAL version)
  return()
endif()

message(FATAL_ERROR "go ${VERSION} not found")
