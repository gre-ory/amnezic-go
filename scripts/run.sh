#!/bin/bash

# args 

PHASE=${1:-loc}
PORT=${2:-9999}

CURRENT_DIR=$PWD
SCRIPTS_DIR=$( dirname $0 )

cd "${SCRIPTS_DIR}/.."

export PORT=${PORT}

BIN="${PWD}/bin/amnezic-go"

# build

make build

# source

source go-source "${BIN}" "${PHASE}"

print-info "LOG_FILE=${LOG_FILE}"
print-info "SERVER_ADDRESS=$( color yellow ${SERVER_ADDRESS} )"
print-info "SQLITE_DATA_SOURCE=${SQLITE_DATA_SOURCE}"

# db up 

go-goose-sqlite "${BIN}" "${PHASE}" up

# run

exec go-run "${BIN}" "${PHASE}"