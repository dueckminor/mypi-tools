#!/usr/bin/env bash

set -ex

DIR_MYPI_TOOLS="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.."; pwd)"

#DIR_MYPI_ROOT="${DIR_MYPI_TOOLS}/.mypi"

PID="$(pgrep -f " --port 9501" || true)"
if [[ -n "${PID}" ]]; then
    kill -9 "${PID}"
fi
PID="$(pgrep -f " --web-port 9501" || true)"
if [[ -n "${PID}" ]]; then
    kill -9 "${PID}"
fi

cd "${DIR_MYPI_TOOLS}/flutter/mypi-setup"

( sleep 5; open "http://localhost:9500" )  &

flutter run -d web --web-port 9501
