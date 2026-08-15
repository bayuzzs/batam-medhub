#!/usr/bin/env bash
# Reset and reseed all four provider databases deterministically using golang-migrate.
set -Eeuo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROVIDERS_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROVIDERS_DIR"

echo "=== Batam MedHub Providers: Deterministic Database Reset ==="

# 1. Clean up active records if tables exist
echo "--> Cleaning existing transaction tables..."
docker compose exec -T postgres psql -U hospital_user -d hospital_db -c "
  DELETE FROM reservations;
  DELETE FROM holds;
  DELETE FROM idempotency_records;
" >/dev/null 2>&1 || true

docker compose exec -T postgres psql -U ferry_user -d ferry_db -c "
  DELETE FROM reservations;
  DELETE FROM holds;
  DELETE FROM idempotency_records;
" >/dev/null 2>&1 || true

docker compose exec -T postgres psql -U hotel_user -d hotel_db -c "
  DELETE FROM reservations;
  DELETE FROM hold_nights;
  DELETE FROM holds;
  DELETE FROM idempotency_records;
" >/dev/null 2>&1 || true

docker compose exec -T postgres psql -U transport_user -d transport_db -c "
  DELETE FROM reservations;
  DELETE FROM hold_assignments;
  DELETE FROM holds;
  DELETE FROM idempotency_records;
" >/dev/null 2>&1 || true

# 2. Roll down all migrations
echo "--> Rolling down migrations across all 4 databases..."
docker compose run --rm hospital-migrate -path=/migrations -database="postgres://hospital_user:hospital_dev_password@postgres:5432/hospital_db?sslmode=disable" down -all
docker compose run --rm ferry-migrate -path=/migrations -database="postgres://ferry_user:ferry_dev_password@postgres:5432/ferry_db?sslmode=disable" down -all
docker compose run --rm hotel-migrate -path=/migrations -database="postgres://hotel_user:hotel_dev_password@postgres:5432/hotel_db?sslmode=disable" down -all
docker compose run --rm transport-migrate -path=/migrations -database="postgres://transport_user:transport_dev_password@postgres:5432/transport_db?sslmode=disable" down -all

# 3. Apply schema & golden seed migrations
echo "--> Applying schema and seed migrations across all 4 databases..."
docker compose run --rm hospital-migrate -path=/migrations -database="postgres://hospital_user:hospital_dev_password@postgres:5432/hospital_db?sslmode=disable" up
docker compose run --rm ferry-migrate -path=/migrations -database="postgres://ferry_user:ferry_dev_password@postgres:5432/ferry_db?sslmode=disable" up
docker compose run --rm hotel-migrate -path=/migrations -database="postgres://hotel_user:hotel_dev_password@postgres:5432/hotel_db?sslmode=disable" up
docker compose run --rm transport-migrate -path=/migrations -database="postgres://transport_user:transport_dev_password@postgres:5432/transport_db?sslmode=disable" up

echo "=== All 4 provider databases successfully reset and seeded! ==="
