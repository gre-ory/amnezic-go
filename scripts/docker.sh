#!/bin/bash

CMD=${1:-run}
TAG="amnezic"

if [[ "${CMD}" == "run" ]]; then

    #
    # docker build
    #

    clear

    >&2 print-info "~> ${TAG}: make build"
    make build

    >&2 print-info "~> ${TAG}: docker build -t ${TAG} ."
    docker build -t ${TAG} .

    >&2 print-info "~> ${TAG}: docker run --env-file scripts/loc.env -p 9090:9090 -p 9091:9091 ${TAG}"
    docker run --env-file scripts/loc.env -p 9090:9090 -p 9091:9091 ${TAG}

else 

    DOCKER_ID=$( docker ps | grep "${TAG}" | head -n 1 | print-1 )
    if [[ "${DOCKER_ID}" == "" ]]; then
        print-warning "${TAG} not running!"
        exit 1
    fi

    if [[ "${CMD}" == "log" ]]; then

        >&2 print-info "~> ${TAG}: docker logs \"${DOCKER_ID}\""
        docker logs "${DOCKER_ID}" 2>&1

    elif [[ "${CMD}" == "rm" ]]; then

        >&2 print-info "~> ${TAG}: docker stop \"${DOCKER_ID}\""
        docker stop "${DOCKER_ID}"

        >&2 print-info "~> ${TAG}: docker rm \"${DOCKER_ID}\""
        docker rm "${DOCKER_ID}"

    else

        print-warning "unknown command! ( $1 )"
        exit 1

    fi

fi
