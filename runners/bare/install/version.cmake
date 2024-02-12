function(get_version _version _binary _module)
  find_file(_binary_path ${_binary})
  execute_process(
	COMMAND go version -m ${_binary_path}
	OUTPUT_VARIABLE _output
  )

  string(REGEX MATCH "[ \t]*mod[ \t]*${_module}[ \t]*v[^ \t]*" _grep "${_output}")
  string(REGEX MATCH "v[^ \t]*$" _tag ${_grep})
  string(SUBSTRING ${_tag} 1 -1 ${_version})

  return(PROPAGATE ${_version})
endfunction()
