project(app
  LANGUAGES NONE
)

add_custom_target(tidy_app
  COMMAND go mod tidy
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
  VERBATIM
)
add_dependencies(tidy
  tidy_app
)

add_custom_target(lint_fix_app
  COMMAND golangci-lint run --fix
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
  VERBATIM
)
add_dependencies(lint_fix
  lint_fix_app
)

add_custom_target(build_app
  COMMAND go build ./...
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
  VERBATIM
)
add_dependencies(build
  build_app
)

add_custom_target(lint_app
  COMMAND golangci-lint run
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
  VERBATIM
)
add_dependencies(lint
  lint_app
)

add_test(
  NAME app
  COMMAND go test ./...
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)
