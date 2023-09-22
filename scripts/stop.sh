#!/bin/bash

PHASE=$( echo ${1} | awk '{ print tolower($0); }' )
if [[ "${PHASE}" == "" ]]; then
    echo "missing phase!"
    exit 1
fi
if [[ "${PHASE}" != "stg" && "${PHASE}" != "prd"  ]]; then
    echo "invalid phase!"
    exit 1
fi

echo "~> sudo systemctl stop amnezic.${PHASE}.service"
sudo systemctl stop amnezic.${PHASE}.service