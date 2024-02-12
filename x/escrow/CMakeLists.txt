project(x-escrow
  LANGUAGES NONE
)

add_go_tidy()
add_go_lint()
add_go_build()
add_go_test()

add_buf_build()
add_buf_format()
add_buf_lint()
add_buf_generate()
