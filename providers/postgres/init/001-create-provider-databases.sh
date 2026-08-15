#!/usr/bin/env bash
set -Eeuo pipefail

hospital_db_password="${HOSPITAL_DB_PASSWORD:-hospital_dev_password}"
ferry_db_password="${FERRY_DB_PASSWORD:-ferry_dev_password}"
hotel_db_password="${HOTEL_DB_PASSWORD:-hotel_dev_password}"
transport_db_password="${TRANSPORT_DB_PASSWORD:-transport_dev_password}"

psql \
  --set=ON_ERROR_STOP=1 \
  --set=hospital_db_password="$hospital_db_password" \
  --set=ferry_db_password="$ferry_db_password" \
  --set=hotel_db_password="$hotel_db_password" \
  --set=transport_db_password="$transport_db_password" \
  --username "$POSTGRES_USER" \
  --dbname "${POSTGRES_DB:-postgres}" <<'SQL'
CREATE ROLE hospital_user
  LOGIN PASSWORD :'hospital_db_password'
  NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT NOREPLICATION NOBYPASSRLS;
CREATE ROLE ferry_user
  LOGIN PASSWORD :'ferry_db_password'
  NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT NOREPLICATION NOBYPASSRLS;
CREATE ROLE hotel_user
  LOGIN PASSWORD :'hotel_db_password'
  NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT NOREPLICATION NOBYPASSRLS;
CREATE ROLE transport_user
  LOGIN PASSWORD :'transport_db_password'
  NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT NOREPLICATION NOBYPASSRLS;

REVOKE ALL PRIVILEGES ON DATABASE postgres FROM PUBLIC;
REVOKE ALL PRIVILEGES ON DATABASE template1 FROM PUBLIC;

CREATE DATABASE hospital_db OWNER hospital_user;
CREATE DATABASE ferry_db OWNER ferry_user;
CREATE DATABASE hotel_db OWNER hotel_user;
CREATE DATABASE transport_db OWNER transport_user;

REVOKE ALL PRIVILEGES ON DATABASE hospital_db FROM PUBLIC;
REVOKE ALL PRIVILEGES ON DATABASE ferry_db FROM PUBLIC;
REVOKE ALL PRIVILEGES ON DATABASE hotel_db FROM PUBLIC;
REVOKE ALL PRIVILEGES ON DATABASE transport_db FROM PUBLIC;

GRANT CONNECT, TEMPORARY ON DATABASE hospital_db TO hospital_user;
GRANT CONNECT, TEMPORARY ON DATABASE ferry_db TO ferry_user;
GRANT CONNECT, TEMPORARY ON DATABASE hotel_db TO hotel_user;
GRANT CONNECT, TEMPORARY ON DATABASE transport_db TO transport_user;

\connect hospital_db
REVOKE ALL PRIVILEGES ON SCHEMA public FROM PUBLIC;
ALTER SCHEMA public OWNER TO hospital_user;
GRANT USAGE, CREATE ON SCHEMA public TO hospital_user;

\connect ferry_db
REVOKE ALL PRIVILEGES ON SCHEMA public FROM PUBLIC;
ALTER SCHEMA public OWNER TO ferry_user;
GRANT USAGE, CREATE ON SCHEMA public TO ferry_user;

\connect hotel_db
REVOKE ALL PRIVILEGES ON SCHEMA public FROM PUBLIC;
ALTER SCHEMA public OWNER TO hotel_user;
GRANT USAGE, CREATE ON SCHEMA public TO hotel_user;

\connect transport_db
REVOKE ALL PRIVILEGES ON SCHEMA public FROM PUBLIC;
ALTER SCHEMA public OWNER TO transport_user;
GRANT USAGE, CREATE ON SCHEMA public TO transport_user;
SQL
