option(USE_DOCKER "use docker to run binaries")

add_custom_target(go_tidy)
add_custom_target(go_lint)
add_custom_target(go_build)

add_custom_target(buf_build)
add_custom_target(buf_format)
add_custom_target(buf_lint)
add_custom_target(buf_generate)

if(USE_DOCKER)
  add_subdirectory(docker)

  block()
	_get_builder_binary_dir(_builder_binary_dir)
	file(COPY bare/install
	  DESTINATION ${_builder_binary_dir}
	)
  endblock()
else()
  add_subdirectory(bare)
endif()
