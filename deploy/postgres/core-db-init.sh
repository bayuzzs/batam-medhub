#!/bin/sh
# One-shot helper that creates the core backend database (core_db) if it does
# not already exist. The provider init script (providers/postgres/init) only
# creates the four provider databases, so the core DB is created here.
set -eu

if psql -h postgres -U provider_admin -d postgres -tAc \
    "SELECT 1 FROM pg_database WHERE datname='core_db'" | grep -q 1; then
  echo "core_db already exists; skipping creation."
  exit 0
fi

psql -h postgres -U provider_admin -d postgres \
  -v ON_ERROR_STOP=1 \
  -c "CREATE DATABASE core_db OWNER provider_admin;"

echo "core_db created."
