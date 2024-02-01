cmake_minimum_required(VERSION 3.25.1)

cmake_policy(SET CMP0140 NEW)

option(CMAKE_COLOR_DIAGNOSTICS
  "Enable/Disable color output during build.")

project(andromeda
  LANGUAGES NONE
)
enable_testing()

add_subdirectory(runners)

add_subdirectory(app)
add_subdirectory(cmd/and)

add_subdirectory(x/escrow)
add_subdirectory(x/test)
