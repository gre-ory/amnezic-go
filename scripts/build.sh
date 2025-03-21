#!/bin/bash

CURRENT_DIR=$PWD
SCRIPTS_DIR=$( dirname $0 )

cd "${SCRIPTS_DIR}/.."

make build

cd $CURRENT_DIR