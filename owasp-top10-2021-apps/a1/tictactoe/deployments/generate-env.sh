#!/bin/sh
#
# This script creates new environment variables every time it runs.
#

DATABASE="a5db"
DATABASE_USER=User$RANDOM$RANDOM
DATABASE_PASSWORD=Pass$RANDOM$RANDOM
# Generate a cryptographically secure SECRET when possible
# Prefer openssl, then od/hexdump, then base64 from /dev/urandom; fall back to previous behavior if none available
if command -v openssl >/dev/null 2>&1; then
  SECRET=$(openssl rand -hex 32)
elif [ -r /dev/urandom ] && command -v od >/dev/null 2>&1; then
  SECRET=$(od -An -N32 -v -tx1 /dev/urandom | tr -d ' \n')
elif [ -r /dev/urandom ] && command -v head >/dev/null 2>&1 && command -v base64 >/dev/null 2>&1; then
  SECRET=$(head -c 32 /dev/urandom | base64 | tr -d '\n')
else
  # Fallback for environments without crypto tools (preserves previous behavior)
  SECRET=$RANDOM$RANDOM
fi
PORT=3306
DATABASE_HOST="mysqldb"

echo "#" > deployments/.dockers.env
echo "# This file is auto generated and contains all environment variables needed by Tic-Tac-Toe's database" >> deployments/.dockers.env
echo "#" >> deployments/.dockers.env
echo "SECRET=$SECRET" >> deployments/.dockers.env
echo "DB_PASSWORD=$DATABASE_PASSWORD" >> deployments/.dockers.env
echo "MYSQL_PASSWORD=$DATABASE_PASSWORD" >> deployments/.dockers.env
echo "DB_DATABASE=$DATABASE" >> deployments/.dockers.env
echo "MYSQL_DATABASE=$DATABASE" >> deployments/.dockers.env
echo "DB_USER=$DATABASE_USER" >> deployments/.dockers.env
echo "MYSQL_USER=$DATABASE_USER" >> deployments/.dockers.env
echo "MYSQL_ROOT_PASSWORD=$DATABASE_PASSWORD" >> deployments/.dockers.env
echo "DB_HOST=$DATABASE_HOST" >> deployments/.dockers.env
echo "DB_PORT=$PORT" >> deployments/.dockers.env