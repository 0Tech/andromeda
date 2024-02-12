function(add_docker_target _type)
  _get_builder_binary_dir(_builder_binary_dir)

  set(_target ${_type}_${PROJECT_NAME})

  add_custom_target(${_target}
	COMMAND cmake -D IDFILE=${_builder_binary_dir}/id.txt -D SOURCE_DIR_SRC=${CMAKE_SOURCE_DIR} -D SOURCE_DIR_DST=/workspace/source -D BINARY_DIR_SRC=${CMAKE_CURRENT_BINARY_DIR} -D BINARY_DIR_DST=/workspace/binary -D "COMMAND=rm -rf /workspace/binary/docker && cmake -S /workspace/source -B /workspace/binary/docker -D FIX_IS_ERROR=${FIX_IS_ERROR} && cmake --build /workspace/binary/docker --target ${_target}" -P ${CMAKE_CURRENT_FUNCTION_LIST_DIR}/docker-run.cmake
	VERBATIM
  )
  add_dependencies(${_target}
	install_builders
  )
  add_dependencies(${_type}
	${_target}
  )
endfunction()

function(add_docker_test)
  _get_builder_binary_dir(_builder_binary_dir)

  set(_target ${PROJECT_NAME})

  add_test(
	NAME ${_target}
	COMMAND cmake -D IDFILE=${_builder_binary_dir}/id.txt -D SOURCE_DIR_SRC=${CMAKE_SOURCE_DIR} -D SOURCE_DIR_DST=/workspace/source -D BINARY_DIR_SRC=${CMAKE_CURRENT_BINARY_DIR} -D BINARY_DIR_DST=/workspace/binary -D "COMMAND=rm -rf /workspace/binary/docker && cmake -S /workspace/source -B /workspace/binary/docker -D FIX_IS_ERROR=${FIX_IS_ERROR} && ctest --test-dir /workspace/binary/docker --output-on-failure -R ${_target}" -P ${CMAKE_CURRENT_FUNCTION_LIST_DIR}/docker-run.cmake
  )
endfunction()

function(add_go_tidy)
  cache_go_mod()
  add_docker_target(go_tidy)
endfunction()

function(add_go_lint)
  add_docker_target(go_lint)
endfunction()

function(add_go_build)
  add_docker_target(go_build)
endfunction()

function(add_go_build_executable)
  add_docker_target(go_build)
endfunction()

function(add_go_test)
  add_docker_test()
endfunction()

function(add_buf_build)
  add_docker_target(buf_build)
endfunction()

function(add_buf_format)
  add_docker_target(buf_format)
endfunction()

function(add_buf_lint)
  add_docker_target(buf_lint)
endfunction()

function(add_buf_generate)
  add_docker_target(buf_generate)
endfunction()
