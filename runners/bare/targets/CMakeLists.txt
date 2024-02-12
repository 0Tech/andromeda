function(add_go_tidy)
  set(_target go_tidy_${PROJECT_NAME})
  add_custom_target(${_target}
	COMMAND go mod tidy
	WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
	VERBATIM
  )
  add_dependencies(${_target}
	install_go
  )
  add_dependencies(go_tidy
	${_target}
  )
endfunction()

function(add_go_lint_fix)
  set(_target go_lint_fix_${PROJECT_NAME})
  add_custom_target(${_target}
	COMMAND golangci-lint run --fix
	WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
	VERBATIM
  )
  add_dependencies(${_target}
	install_golangci-lint
  )
  add_dependencies(go_lint_fix
	${_target}
  )
endfunction()

function(add_go_build)
  set(_target go_build_${PROJECT_NAME})
  add_custom_target(${_target}
	COMMAND go build ./...
	WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
	VERBATIM
  )
  add_dependencies(${_target}
	install_go
  )
  add_dependencies(go_build
	${_target}
  )
endfunction()

function(add_go_build_executable)
  set(_target go_build_${PROJECT_NAME})
  add_custom_target(${_target}
	COMMAND go build -o ${CMAKE_BINARY_DIR} ./...
	WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
	VERBATIM
  )
  add_dependencies(${_target}
	install_go
  )
  add_dependencies(go_build
	${_target}
  )
endfunction()

function(add_go_lint)
  set(_target go_lint_${PROJECT_NAME})
  add_custom_target(${_target}
	COMMAND golangci-lint run
	WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
	VERBATIM
  )
  add_dependencies(${_target}
	install_golangci-lint
  )
  add_dependencies(go_lint
	${_target}
  )
endfunction()

function(add_go_test)
  add_test(
	NAME ${PROJECT_NAME}
	COMMAND go test -coverprofile ${CMAKE_CURRENT_BINARY_DIR}/cover.out ./...
	WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
  )
endfunction()

function(add_buf_build)
  set(_target buf_build_${PROJECT_NAME})
  add_custom_target(${_target}
	COMMAND buf build
	WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/proto
	VERBATIM
  )
  add_dependencies(${_target}
	install_buf
  )
  add_dependencies(buf_build
	${_target}
  )
endfunction()

function(add_buf_lint)
  set(_target buf_lint_${PROJECT_NAME})
  add_custom_target(${_target}
	COMMAND buf lint
	WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/proto
	VERBATIM
  )
  add_dependencies(${_target}
	install_buf
  )
  add_dependencies(buf_lint
	${_target}
  )
endfunction()

function(add_buf_generate)
  set(_target buf_generate_${PROJECT_NAME})
  add_custom_target(${_target}
	COMMAND buf generate --template buf.gen.gogo.yaml
	COMMAND buf generate --template buf.gen.pulsar.yaml
	WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/proto
	VERBATIM
  )
  add_dependencies(${_target}
	install_buf
	install_protoc-gen-go
	install_protoc-gen-go-grpc
  )
  add_dependencies(buf_generate
	${_target}
  )
endfunction()
