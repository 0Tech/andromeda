configure_file(Dockerfile Dockerfile
  COPYONLY
)

set(UID 1000 CACHE STRING "uid of the user")
set(GID 1000 CACHE STRING "gid of the user")

function(_get_builder_binary_dir _output)
  file(RELATIVE_PATH _rel_builder_source_dir ${CMAKE_SOURCE_DIR} ${CMAKE_CURRENT_FUNCTION_LIST_DIR})
  file(REAL_PATH ${_rel_builder_source_dir} ${_output} BASE_DIRECTORY ${CMAKE_BINARY_DIR})
  return(PROPAGATE ${_output})
endfunction()

add_custom_target(install_builders
  COMMAND cmake -D IDFILE=id.txt -D UID=${UID} -D GID=${GID} -P ${CMAKE_CURRENT_SOURCE_DIR}/docker-build.cmake
)

file(COPY ${CMAKE_SOURCE_DIR}/go.work ${CMAKE_SOURCE_DIR}/go.work.sum
  DESTINATION mods/
)
function(cache_go_mod)
  _get_builder_binary_dir(_builder_binary_dir)

  file(RELATIVE_PATH _rel_source_dir ${CMAKE_SOURCE_DIR} ${CMAKE_CURRENT_SOURCE_DIR})

  file(COPY go.mod go.sum
	DESTINATION ${_builder_binary_dir}/mods/${_rel_source_dir}/
  )
endfunction()
