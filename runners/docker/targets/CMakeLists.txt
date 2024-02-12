function(add_docker_target _type)
  _get_builder_binary_dir(_builder_binary_dir)

  set(_target ${_type}_${PROJECT_NAME})

  add_custom_target(${_target}
	COMMAND cmake -D IDFILE=${_builder_binary_dir}/id.txt -D BIND_DIR=${CMAKE_SOURCE_DIR} -D "COMMAND=cmake -S . -B /tmp/build -D FIX_IS_ERROR=${FIX_IS_ERROR} && cmake --build /tmp/build --target ${_target}" -P ${CMAKE_CURRENT_FUNCTION_LIST_DIR}/docker-run.cmake
	VERBATIM
  )
  add_dependencies(${_target}
	install_builders
  )

  set(_root ${_type})
  if(_type STREQUAL go_build_executable)
	set(_root go_build)
  endif()
  add_dependencies(${_root}
	${_target}
  )
endfunction()

function(add_docker_test)
  _get_builder_binary_dir(_builder_binary_dir)

  set(_target ${PROJECT_NAME})

  add_test(
	NAME ${_target}
	COMMAND cmake -D IDFILE=${_builder_binary_dir}/id.txt -D BIND_DIR=${CMAKE_SOURCE_DIR} -D "COMMAND=cmake -S . -B /tmp/build -D FIX_IS_ERROR=${FIX_IS_ERROR} && ctest --test-dir /tmp/build --output-on-failure -R ${_target}" -P ${CMAKE_CURRENT_FUNCTION_LIST_DIR}/docker-run.cmake
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
  add_docker_target(go_build_executable)
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
