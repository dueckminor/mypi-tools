#!/usr/bin/env bash

set -ex

DIR_MYPI_AUTH="$(cd "$(dirname "${BASH_SOURCE[0]}")/.."; pwd)"

DIR_MYPI_ROOT="${DIR_MYPI_AUTH}/.mypi"

DIR_AUTH_SERVER="${DIR_MYPI_ROOT}/etc/mypi-auth/server"
DIR_AUTH_CLIENTS="${DIR_MYPI_ROOT}/etc/mypi-auth/clients"

mkdir -p "${DIR_AUTH_CLIENTS}"
mkdir -p "${DIR_AUTH_SERVER}"

if [[ ! -f "${DIR_MYPI_ROOT}/etc/mypi-auth/users.yml" ]]; then
    cat > "${DIR_MYPI_ROOT}/etc/mypi-auth/users.yml" <<__EOF__
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



PID="$(pgrep -f " --port 9101" || true)"
if [[ -n "${PID}" ]]; then
    kill -9 "${PID}"
fi

npm run dev
