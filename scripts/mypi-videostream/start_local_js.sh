#!/usr/bin/env bash

set -ex

DIR_MYPI_TOOLS="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.."; pwd)"

if [[ ! -d "${DIR_MYPI_TOOLS}/../mypi-workspace/.git" ]]; then
    echo
    echo "Please do the following:"
    echo 
    echo "cd \"$(cd ${DIR_MYPI_TOOLS}/..)\""
    echo "git clone git@github.com:dueckminor/mypi-workspace.git"
    echo
    exit 1
fi

DIR_MYPI_WORKSPACE="$(cd "${DIR_MYPI_TOOLS}/../mypi-workspace"; pwd)"
DIR_MYPI_ROOT="${DIR_MYPI_WORKSPACE}/.mypi"

PID="$(pgrep -f " --port 8084" || true)"
if [[ -n "${PID}" ]]; then
    kill -9 "${PID}"
fi

cd "${DIR_MYPI_TOOLS}/web/mypi-videostream"

npm run dev
