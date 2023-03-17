#!/bin/bash
set -e
set -x

REACT_DIRECTORY=${1:-$( dirname $PWD )/amnezic-react}
if [[ ! -e "${REACT_DIRECTORY}" ]]; then
    echo "missing react directory"
    exit 1
fi

CURRENT_DIRECTORY=${PWD}
FINAL_DIRECTORY="${PWD}/www"
BUILD_DIRECTORY="${REACT_DIRECTORY}/build"

rm -rf "${BUILD_DIRECTORY}" "${FINAL_DIRECTORY}"

cd ${REACT_DIRECTORY}
yarn build
cp -r "${BUILD_DIRECTORY}" "${FINAL_DIRECTORY}"
cd -
