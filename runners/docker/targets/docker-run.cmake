file(READ ${IDFILE} id)
message(STATUS "Running target on Docker ${id}")
execute_process(
  COMMAND docker run --rm --mount type=bind,src=${BIND_DIR},dst=/workspace --workdir /workspace ${id} sh -c ${COMMAND}
  COMMAND_ERROR_IS_FATAL ANY
)
