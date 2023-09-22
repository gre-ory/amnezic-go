#!/bin/bash

TAG=${1}
if [[ "${TAG}" == "" ]]; then
    echo "missing tag!"
    exit 1
fi

echo "~> git fetch --all --tags"
git fetch --all --tags

echo "~> git what \"${TAG}\""
git what "${TAG}" > /dev/null 2> /dev/null  || exit 1

echo "~> git checkout \"${TAG}\""
git checkout "${TAG}"

echo "~> make build"
make build
