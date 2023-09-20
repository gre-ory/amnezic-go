#!/bin/bash

APP_NAME="amnezic-go"

_term() {
  # child=$( ps ux | grep "./bin/${APP_NAME}" | grep -v "grep" | print-2 )
  >&2 echo "~> kill -TERM \"$child\" 2>/dev/null"
  kill -TERM "$child" 2>/dev/null
}
trap _term SIGTERM
trap _term SIGINT

PHASE=$( echo ${1:-loc} | awk '{ print tolower($0); }' )

ENV_FILE="$( dirname $0 )/${PHASE}.env"
if [[ ! -e "${ENV_FILE}" ]]; then
    >&2 echo -e "\033[0;31m missing env file for ${PHASE}! ( ${ENV_FILE} ) \033[0m"
    exit 1
fi

set -o allexport
source "${ENV_FILE}"
set +o allexport

>&2 echo "~> PHASE = ${PHASE}"
>&2 echo "~> ENV_FILE = ${ENV_FILE}"

>&2 echo "~> ./bin/${APP_NAME} &"
./bin/${APP_NAME} &

child=$!
>&2 echo "~> wait \"$child\""
wait "$child"
