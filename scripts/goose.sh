#!/bin/bash

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
>&2 echo "~> SQLITE_DATA_SOURCE = ${SQLITE_DATA_SOURCE}"

>&2 echo "~> goose -dir ./db sqlite3 ${SQLITE_DATA_SOURCE} ${2:-up}"
goose -dir ./db sqlite3 ${SQLITE_DATA_SOURCE} ${2:-up}
