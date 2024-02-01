configure_file(Dockerfile.in Dockerfile
  @ONLY
)

function(_get_builder_working_dir _output)
  set(${_output} /workspace)
  return(PROPAGATE ${_output})
endfunction()

function(_get_builder_image _output)
  _get_builder_working_dir(_builder_working_dir)

  file(RELATIVE_PATH _rel_builder_source_dir ${CMAKE_SOURCE_DIR} ${CMAKE_CURRENT_FUNCTION_LIST_DIR})
  file(REAL_PATH ${_rel_builder_source_dir} _builder_binary_dir BASE_DIRECTORY ${CMAKE_BINARY_DIR})

  execute_process(
	COMMAND docker build -q --build-arg workdir=${_builder_working_dir} ${_builder_binary_dir}
	OUTPUT_VARIABLE ${_output}
  )
  string(STRIP ${${_output}} ${_output})

  return(PROPAGATE ${_output})
endfunction()
