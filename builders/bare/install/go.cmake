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

require_variables(VERSION CMAKE_SYSTEM_NAME CMAKE_SYSTEM_PROCESSOR)

get_version(version)
block()
  coerce(${VERSION} expecting)
  coerce("${version}" actual)
  if(actual VERSION_EQUAL expecting)
	return()
  endif()
endblock()

message(WARNING "No go ${VERSION} found. Automatic installation is destructive and requires write permission on your filesystem. Manual installation is highly recommended.")

# Determine Operating System
if(${CMAKE_SYSTEM_NAME} STREQUAL "Linux")
  set(OS "linux")
elseif(${CMAKE_SYSTEM_NAME} STREQUAL "Windows")
  set(OS "windows")
elseif(${CMAKE_SYSTEM_NAME} STREQUAL "Darwin")
  set(OS "darwin")
else()
  message(FATAL_ERROR "Unsupported operating system.")
endif()

# Determine Architecture
if(${CMAKE_SYSTEM_PROCESSOR} STREQUAL "x86_64")
  set(ARCH "amd64")
elseif(${CMAKE_SYSTEM_PROCESSOR} MATCHES "arm")
  set(ARCH "arm64")
else()
  message(FATAL_ERROR "Unsupported architecture.")
endif()

set(GO_BINARY_URL "https://dl.google.com/go/go${VERSION}.${OS}-${ARCH}.tar.gz")

file(DOWNLOAD ${GO_BINARY_URL} go.tar.gz)

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
