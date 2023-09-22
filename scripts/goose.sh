#!/bin/bash

# go install github.com/pressly/goose/v3/cmd/goose@latest

PHASE=$( echo ${1:-loc} | awk '{ print tolower($0); }' )

SCRIPT_DIR=$( dirname $0 )
DB_DIR=$( cd "${SCRIPT_DIR}/../db"; pwd )
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
>&2 echo "~> SQLITE_DATA_SOURCE = ${SQLITE_DATA_SOURCE}"

>&2 echo "~> goose -dir ${DB_DIR} sqlite3 ${SQLITE_DATA_SOURCE} ${2:-up}"
goose -dir ${DB_DIR} sqlite3 ${SQLITE_DATA_SOURCE} ${2:-up}
