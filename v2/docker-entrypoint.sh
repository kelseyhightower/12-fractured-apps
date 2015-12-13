#!/bin/sh
set -e

datadir=${APP_DATADIR:="/var/lib/data"}
host=${APP_HOST:="127.0.0.1"}
port=${APP_PORT:="3306"}
username=${APP_USERNAME:=""}
password=${APP_PASSWORD:=""}
database=${APP_DATABASE:=""}

cat <<EOF > /etc/config.json
{
  "datadir": "${datadir}",
  "host": "${host}",
  "port": "${port}",
  "username": "${username}",
  "password": "${password}",
  "database": "${database}"
}
EOF

mkdir -p ${APP_DATADIR}

exec "/app"
