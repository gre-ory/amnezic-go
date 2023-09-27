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

SCRIPT_DIR=$( dirname $0 )
BIN_DIR=$( cd "${SCRIPT_DIR}/../bin"; pwd )
ENV_FILE="${SCRIPT_DIR}/${PHASE}.env"
if [[ ! -e "${ENV_FILE}" ]]; then
    >&2 echo -e "\033[0;31m missing env file for ${PHASE}! ( ${ENV_FILE} ) \033[0m"
    exit 1
fi

set -o allexport
source "${ENV_FILE}"
set +o allexport

>&2 echo "~> PHASE = ${PHASE}"
>&2 echo "~> ENV_FILE = ${ENV_FILE}"
>&2 echo "~> FRONTEND_ADDRESS = ${FRONTEND_ADDRESS}"
>&2 echo "~> BACKEND_ADDRESS = ${BACKEND_ADDRESS}"
>&2 echo "~> LOG_FILE = ${LOG_FILE}"

export APPLICATION_NAME=amnezic
export APPLICATION_VERSION=$( cd ${SCRIPT_DIR} > /dev/null 2> /dev/null ; git describe --exact-match --tags $( git log -n1 --pretty='%h' ) 2> /dev/null || echo "v9.9.9"; cd - > /dev/null 2> /dev/null )

>&2 echo "~> APPLICATION_NAME = ${APPLICATION_NAME}"
>&2 echo "~> APPLICATION_VERSION = ${APPLICATION_VERSION}"

>&2 echo "~> ${BIN_DIR}/${APP_NAME} &"
${BIN_DIR}/${APP_NAME} &

child=$!
>&2 echo "~> wait \"$child\""
wait "$child"
