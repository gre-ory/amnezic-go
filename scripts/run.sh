#!/bin/bash

# args 

PHASE=${1:-loc}
PORT=${2:-9999}

CURRENT_DIR=$PWD
SCRIPTS_DIR=$( dirname $0 )

cd "${SCRIPTS_DIR}/.."

export PORT=${PORT}

BIN="${PWD}/bin/amnezic-go"

# source

source go-source "${BIN}" "${PHASE}"

print-info "SERVER_ADDRESS=$( color yellow ${SERVER_ADDRESS} )"
print-info "SQLITE_DATA_SOURCE=${SQLITE_DATA_SOURCE}"

# db up 

go-goose-sqlite "${BIN}" "${PHASE}" up

# run

go-run "${BIN}" "${PHASE}"

cd $CURRENT_DIR