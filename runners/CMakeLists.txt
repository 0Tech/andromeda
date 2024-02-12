set(USE_DOCKER TRUE CACHE BOOL "use docker to run binaries")

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
