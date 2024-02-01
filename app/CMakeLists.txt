project(app
  LANGUAGES NONE
)

add_go_tidy()
add_go_lint_fix()
add_go_build()
add_go_lint()
add_go_test()
