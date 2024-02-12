function(get_version _version _binary _module)
  set(${_version} 0)

  find_file(_binary_path ${_binary})
  if(_binary_path STREQUAL _binary_path-NOTFOUND)
	return(PROPAGATE ${_version})
  endif()

  execute_process(
	COMMAND go version -m ${_binary_path}
	OUTPUT_VARIABLE _output
  )

  string(REGEX MATCH "[ \t]*mod[ \t]*${_module}[ \t]*v[^ \t]*" _grep "${_output}")
  string(REGEX MATCH "v[^ \t]*$" _tag ${_grep})
  string(SUBSTRING ${_tag} 1 -1 ${_version})

  return(PROPAGATE ${_version})
endfunction()
