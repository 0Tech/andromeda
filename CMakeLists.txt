cmake_minimum_required(VERSION 3.25.1)

option(CMAKE_COLOR_DIAGNOSTICS
  "Enable/Disable color output during build.")

project(andromeda
  LANGUAGES NONE
)
enable_testing()

# go
add_custom_target(tidy)
add_custom_target(lint_fix)
add_custom_target(build ALL)
add_custom_target(lint)

# buf
add_custom_target(buf_build)
add_custom_target(buf_lint)
add_custom_target(buf_generate)

add_subdirectory(install)

add_subdirectory(app)
add_subdirectory(cmd/and)

add_subdirectory(x/escrow)
add_subdirectory(x/test)
