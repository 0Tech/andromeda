#!/usr/bin/env cmake -P

cmake_policy(SET CMP0140 NEW)

include(common.cmake)

function(get_version _version)
  set(${_version} 0)

  execute_process(
	COMMAND go version
	OUTPUT_VARIABLE _output
  )
  string(REGEX MATCH "[0-9](\\.[0-9]+)*" ${_version} "${_output}")
  return(PROPAGATE ${_version})
endfunction()

function(coerce _input _output)
  string(REGEX MATCH "[0-9]\\.[0-9]+" ${_output} "${_input}")
  return(PROPAGATE ${_output})
endfunction()

require_variables(VERSION)

get_version(version)
block()
  coerce(${VERSION} expecting)
  coerce("${version}" actual)
  if(actual VERSION_EQUAL expecting)
	return()
  endif()
endblock()

message(WARNING "No go ${VERSION} found. Automatic installation is destructive and requires write permission on your filesystem. Manual installation is highly recommended.")

# TODO: fix hard coding
file(DOWNLOAD https://go.dev/dl/go${VERSION}.linux-amd64.tar.gz go.tar.gz)

# libarchive bug
if(NOT DEFINED ENV{LANG})
  set(ENV{LANG} C.UTF-8)
endif()
file(ARCHIVE_EXTRACT
  INPUT ${CMAKE_CURRENT_BINARY_DIR}/go.tar.gz
  TOUCH
)

file(INSTALL ${CMAKE_CURRENT_BINARY_DIR}/go/
  DESTINATION /usr/local
  USE_SOURCE_PERMISSIONS
)

message(STATUS "Installed ${VERSION}")
