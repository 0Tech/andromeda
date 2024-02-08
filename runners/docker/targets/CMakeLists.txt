function(add_docker_target _type _project)
  _get_builder_working_dir(_builder_working_dir)
  _get_builder_image(_builder_image)

  set(_target ${_type})
  if(_project)
	string(APPEND _target _${_project})
  endif()

  add_custom_target(${_target}
	COMMAND docker run --rm --mount type=bind,src=${CMAKE_SOURCE_DIR},dst=${_builder_working_dir} ${_builder_image} sh -c "cmake -S . -B /tmp/build -D USE_DOCKER=NO && cmake --build /tmp/build --target ${_target}"
	VERBATIM
  )
endfunction()

function(add_docker_test _project)
  _get_builder_working_dir(_builder_working_dir)
  _get_builder_image(_builder_image)

  set(_target ${_project})

  add_test(
	NAME ${_target}
	COMMAND docker run --rm --mount type=bind,src=${CMAKE_SOURCE_DIR},dst=${_builder_working_dir} ${_builder_image} sh -c "cmake -S . -B /tmp/build -D USE_DOCKER=NO && ctest --test-dir /tmp/build -R ${_target}"
  )
endfunction()

function(cache_go_mod)
  _get_builder_binary_dir(_builder_binary_dir)

  file(RELATIVE_PATH _rel_source_dir ${CMAKE_SOURCE_DIR} ${CMAKE_CURRENT_SOURCE_DIR})

  configure_file(go.mod ${_builder_binary_dir}/mods/${_rel_source_dir}/go.mod
	COPYONLY
  )
  configure_file(go.sum ${_builder_binary_dir}/mods/${_rel_source_dir}/go.sum
	COPYONLY
  )
endfunction()

function(add_go_tidy)
  cache_go_mod()
  add_docker_target(go_tidy ${PROJECT_NAME})
endfunction()
add_docker_target(go_tidy "")

function(add_go_lint_fix)
  add_docker_target(go_lint_fix ${PROJECT_NAME})
endfunction()
add_docker_target(go_lint_fix "")

function(add_go_build)
  add_docker_target(go_build ${PROJECT_NAME})
endfunction()
add_docker_target(go_build "")

function(add_go_build_executable)
  add_docker_target(go_build_executable ${PROJECT_NAME})
endfunction()
add_docker_target(go_build_executable "")

function(add_go_lint)
  add_docker_target(go_lint ${PROJECT_NAME})
endfunction()
add_docker_target(go_lint "")

function(add_go_test)
  add_docker_test(${PROJECT_NAME})
endfunction()

function(add_buf_build)
  add_docker_target(buf_build ${PROJECT_NAME})
endfunction()
add_docker_target(buf_build "")

function(add_buf_lint)
  add_docker_target(buf_lint ${PROJECT_NAME})
endfunction()
add_docker_target(buf_lint "")

function(add_buf_generate)
  add_docker_target(buf_generate ${PROJECT_NAME})
endfunction()
add_docker_target(buf_generate "")
