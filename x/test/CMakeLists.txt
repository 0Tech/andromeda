project(x-test
  LANGUAGES NONE
)

add_custom_target(tidy_x-test
  COMMAND go mod tidy
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
  VERBATIM
)
add_dependencies(tidy
  tidy_x-test
)

add_custom_target(lint_x-test
  COMMAND golangci-lint run --fix
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
  VERBATIM
)
add_dependencies(lint
  lint_x-test
)

add_custom_target(build_x-test
  COMMAND go build ./...
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
  VERBATIM
)
add_dependencies(build
  build_x-test
)

add_test(
  NAME x-test
  COMMAND go test ./...
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)

add_custom_target(buf_build_x-test
  COMMAND buf build
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/proto
  VERBATIM
)
add_dependencies(buf_build
  buf_build_x-test
)

add_custom_target(buf_lint_x-test
  COMMAND buf lint
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/proto
  VERBATIM
)
add_dependencies(buf_lint
  buf_lint_x-test
)

add_custom_target(buf_generate_x-test
  COMMAND buf generate --template buf.gen.gogo.yaml
  COMMAND buf generate --template buf.gen.pulsar.yaml
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/proto
  VERBATIM
)
add_dependencies(buf_generate
  buf_generate_x-test
)
