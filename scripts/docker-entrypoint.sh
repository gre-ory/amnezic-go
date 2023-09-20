#!/bin/sh
set -e

: ${AMNEZIC_HTTP_TIME_OUT_MS:=0}
: ${AMNEZIC_ARG:=serve}

if [ "$1" = "amnezic" ]; then

	# trigger export job
	echo "[docker] amnezic ${AMNEZIC_ARG}"
	exec server \
		--http-time-out-ms="${AMNEZIC_HTTP_TIME_OUT_MS}" \
		-- \
		"${AMNEZIC_ARG}"

fi

exec "$@"