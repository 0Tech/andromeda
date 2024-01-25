project(and
  LANGUAGES NONE
)

add_custom_target(tidy_and
  COMMAND go mod tidy
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
  VERBATIM
)
add_dependencies(tidy
  tidy_and
)

add_custom_target(lint_fix_and
  COMMAND golangci-lint run --fix
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
  VERBATIM
)
add_dependencies(lint_fix
  lint_fix_and
)

add_custom_target(build_and
  COMMAND go build -o ${andromeda_BINARY_DIR} ./...
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
  VERBATIM
)
add_dependencies(build
  build_and
)

add_custom_target(lint_and
  COMMAND golangci-lint run
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
  VERBATIM
)
add_dependencies(lint
  lint_and
)

add_test(
  NAME and
  COMMAND go test ./...
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)
