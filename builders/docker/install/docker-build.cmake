execute_process(
  COMMAND docker build --build-arg UID=${UID} --build-arg GID=${GID} -q .
  OUTPUT_VARIABLE id
)
string(STRIP ${id} id)
file(WRITE ${IDFILE} ${id})
