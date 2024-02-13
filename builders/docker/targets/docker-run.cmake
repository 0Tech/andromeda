file(READ ${IDFILE} id)
message(STATUS "Running target on Docker ${id}")
execute_process(
  COMMAND docker run --rm --mount type=bind,src=${SOURCE_DIR_SRC},dst=${SOURCE_DIR_DST} --mount type=bind,src=${BINARY_DIR_SRC},dst=${BINARY_DIR_DST} ${id} sh -c ${COMMAND}
  COMMAND_ERROR_IS_FATAL ANY
)
