#!/bin/sh

CONFIG_FILE="/testdata/netflowv5.yaml"
LOG_FILE="/testdata/stanza.log"
STDOUT_FILE="/testdata/stdout.log"
OUTPUT_FILE="/testdata/out.log"

# clear the log if it exists, is is crucial that each test
# starts with empty files
> "${LOG_FILE}"
> "${STDOUT_FILE}"
> "${OUTPUT_FILE}"

chmod 0666 $LOG_FILE
chmod 0666 $STDOUT_FILE
chmod 0666 $OUTPUT_FILE

/stanza_home/stanza \
    --config "${CONFIG_FILE}" \
    --log_file "${LOG_FILE}" >"${STDOUT_FILE}" 2>&1
