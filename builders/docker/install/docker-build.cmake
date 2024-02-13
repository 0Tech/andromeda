execute_process(
  COMMAND docker build -q .
  OUTPUT_VARIABLE id
)
string(STRIP ${id} id)
file(WRITE ${IDFILE} ${id})
