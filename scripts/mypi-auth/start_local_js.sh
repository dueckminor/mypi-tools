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

DIR_AUTH_SERVER="${DIR_MYPI_ROOT}/etc/auth/server"
DIR_AUTH_CLIENTS="${DIR_MYPI_ROOT}/etc/auth/clients"

mkdir -p "${DIR_AUTH_CLIENTS}"
mkdir -p "${DIR_AUTH_SERVER}"

if [[ ! -f "${DIR_MYPI_ROOT}/etc/users.yml" ]]; then
    cat > "${DIR_MYPI_ROOT}/etc/users.yml" <<__EOF__
- name: admin
  password: \$2a\$10\$MsH/XNruBjWZ0irP06JiWuWvHPyiHxfFaymQZtaSGfFtMHibE/iKi
__EOF__
fi

if [[ ! -f "${DIR_AUTH_SERVER}/server_priv.pem" ]]; then
    openssl genpkey -algorithm RSA -out "${DIR_AUTH_SERVER}/server_priv.pem" -pkeyopt rsa_keygen_bits:2048
    openssl rsa -pubout -in "${DIR_AUTH_SERVER}/server_priv.pem" -out "${DIR_AUTH_SERVER}/server_pub.pem"
fi

if [[ ! -f "${DIR_AUTH_CLIENTS}/sample.yml" ]]; then
    SECRET="$(openssl rand -hex 16)"
    PUB="$(cat "${DIR_AUTH_SERVER}/server_pub.pem")"

    jq -n --arg client_secret "${SECRET}" --arg server_key "${PUB}" \
    '{
        "client_id":"sample", 
        "client_secret":$client_secret,
        "server_key":$server_key
    }' > "${DIR_AUTH_CLIENTS}/sample.yml"
fi

PID="$(pgrep -f " --port 8082" || true)"
if [[ -n "${PID}" ]]; then
    kill -9 "${PID}"
fi

cd "${DIR_MYPI_TOOLS}/web/mypi-auth"

npm run dev
