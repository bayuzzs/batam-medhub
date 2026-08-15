# Batam MedHub Provider Platform

This directory contains the four standalone, headless Go/Gin mock provider services for the Batam MedHub platform:

1. **Hospital Provider (`hospital`)**: Manages medical check-up (MCU) services, appointments, doctor consultation slots, holds, and reservations.
2. **Ferry Provider (`ferry`)**: Manages cross-strait sailings, seat capacity, holds, and reservations between Singapore and Batam.
3. **Hotel Provider (`hotel`)**: Manages room types, multi-night stay inventory, atomic date-range holds, and reservations.
4. **Transport Provider (`transport`)**: Manages ground vehicles, route availability slots, holds with required booking requirements, and reservations.

---

## 1. Architecture and Runtime Topology

All four provider services run as standalone Go binaries in isolated containers managed by `docker-compose.yml`. They connect to a shared PostgreSQL 17 server hosting **four strictly isolated logical databases**.

```
+-------------------------------------------------------------------------+
|                       Batam MedHub Provider Platform                    |
|                                                                         |
|  +--------------------+  +--------------------+  +--------------------+ |
|  | Hospital (:8081)   |  | Ferry (:8082)      |  | Hotel (:8083)      | |
|  | DB: hospital_db    |  | DB: ferry_db       |  | DB: hotel_db       | |
|  | User: hospital_user|  | User: ferry_user   |  | User: hotel_user   | |
|  +---------+----------+  +---------+----------+  +---------+----------+ |
|            |                       |                       |            |
|            |             +---------+----------+            |            |
|            |             | Transport (:8084)  |            |            |
|            |             | DB: transport_db   |            |            |
|            |             | User: transport_user            |            |
|            |             +---------+----------+            |            |
|            |                       |                       |            |
|  +---------v-----------------------v-----------------------v----------+ |
|  |                 PostgreSQL 17 (port 5432)                          | |
|  |  [hospital_db]      [ferry_db]      [hotel_db]      [transport_db] | |
|  +--------------------------------------------------------------------+ |
+-------------------------------------------------------------------------+
```

### Provider Identity and Ports

| Service | Port | Provider ID | Provider Type | Database | Owner User | Default Integration Key |
|---|---|---|---|---|---|---|
| Hospital | `8081` | `hospital-demo-01` | `HOSPITAL` | `hospital_db` | `hospital_user` | `hospital_dev_secret` |
| Ferry | `8082` | `ferry-demo-01` | `FERRY` | `ferry_db` | `ferry_user` | `ferry_dev_secret` |
| Hotel | `8083` | `hotel-demo-01` | `HOTEL` | `hotel_db` | `hotel_user` | `hotel_dev_secret` |
| Transport | `8084` | `transport-demo-01` | `TRANSPORT` | `transport_db` | `transport_user` | `transport_dev_secret` |

---

## 2. Database Isolation and Security

- Each database owner (`hospital_user`, `ferry_user`, `hotel_user`, `transport_user`) has access **only** to its assigned database.
- Public privileges on databases and schemas are revoked during PostgreSQL initialization (`postgres/init/001-create-provider-databases.sh`).
- No provider service can connect to or query any other provider database.
- Services never share database transactions, models, or state across boundaries; all inter-service communication from the core backend is over HTTP.

---

## 3. Configuration and Environment Variables

Configuration is loaded from environment variables. An example file is available at `.env.example`.

```bash
cp .env.example .env
```

Key environment variables:
- `PROVIDER_POSTGRES_PASSWORD`: Superuser password for PostgreSQL (`provider_admin`).
- `HOSPITAL_DB_PASSWORD`, `FERRY_DB_PASSWORD`, `HOTEL_DB_PASSWORD`, `TRANSPORT_DB_PASSWORD`: Individual database passwords.
- `HOSPITAL_INTEGRATION_KEY`, `FERRY_INTEGRATION_KEY`, `HOTEL_INTEGRATION_KEY`, `TRANSPORT_INTEGRATION_KEY`: Shared integration secrets passed via the `X-Integration-Key` header.
- `DATABASE_URL`: Connection string passed to each provider container.
- `PORT`: HTTP port inside the container (`8080`).

---

## 4. Startup and Teardown

### Starting Services with Docker Compose
```bash
cd providers
docker compose up -d --build
```

Compose will:
1. Start `postgres` and run the initial role/database creation script.
2. Run database migration containers (`hospital-migrate`, `ferry-migrate`, `hotel-migrate`, `transport-migrate`) to apply initial schema and golden seed data.
3. Start the four provider web services on ports `8081`–`8084` with automatic health checks.

### Stopping Services
```bash
cd providers
docker compose down
```

To stop and remove data volumes:
```bash
docker compose down -v
```

---

## 5. Migrations and Seed Data

- Database schemas are managed exclusively via `golang-migrate` SQL migrations located in `migrations/<service>/`.
- Runtime GORM `AutoMigrate` is **strictly disabled**.
- Migrations include both table creation (`000001_create_schema.up.sql`) and golden seed datasets (`000002_seed.up.sql` / `000002_seed_golden_data.up.sql`).
- Rollback migrations (`.down.sql`) allow clean teardown and fresh re-creation without affecting other databases.

---

## 6. Deterministic Database Reset

To reset all four databases back to clean, freshly-seeded states:

```bash
cd providers
./scripts/reset.sh
```

Or run directly via docker compose:
```bash
# Example for hospital:
docker compose run --rm hospital-migrate down -all
docker compose run --rm hospital-migrate up
```

---

## 7. Contract Invariants and Implementation Details

- **Strict Bounded JSON Decoding**: Rejects unknown fields (`DisallowUnknownFields`), limits payload size (1MB max), and rejects trailing data.
- **Strict UTC Timestamps**: Timestamps must be RFC 3339 formatted instants ending with uppercase `'Z'` (e.g. `2026-08-22T01:25:00Z`).
- **Money Minor Units**: All monetary values use integer minor units (e.g. `15000000` IDR = IDR 150,000; `2500` SGD = SGD 25.00) with a 3-letter ISO 4217 currency code.
- **Concurrency & Idempotency Serialization**:
  - `POST /v1/holds`, `POST /v1/holds/:id/confirm`, `POST /v1/holds/:id/release`, and `POST /v1/reservations/:id/release` acquire PostgreSQL transaction advisory locks via `SELECT pg_advisory_xact_lock(hashtext(?))` before querying the idempotency table.
  - Concurrent identical requests serialize and replay the same `201 Created` response with `Idempotency-Replayed: true` and zero duplicate capacity consumption.
  - Reusing an idempotency key with a mutated payload returns `409 IDEMPOTENCY_CONFLICT`.
- **Capacity Locking**: Real-time seat, room, or slot capacity is locked using `SELECT ... FOR UPDATE` within the transaction to prevent overselling.
- **Credential-Safe Recovery**: Uncaught panics are caught by `SafeRecoveryMiddleware`, returning a standardized `INTERNAL_ERROR` envelope without exposing stack traces or `X-Integration-Key` secrets in logs.

---

## 8. Smoke Verification

A comprehensive smoke test script tests all contract requirements, concurrency idempotency, lifecycle states, negative validation, and database isolation across all 4 services:

```bash
cd providers
./smoke.sh
```

To run with an automatic database reset before testing:
```bash
./smoke.sh --reset
```

### Verified Test Cases:
1. `GET /healthz` status, database status, and `X-Request-ID` tracing/replacement.
2. `401 AUTHENTICATION_FAILED` on missing or invalid `X-Integration-Key`.
3. `403 PROVIDER_IDENTITY_MISMATCH` on mismatched payload `provider_id`.
4. Offer searches returning deterministic golden offers and exact pricing.
5. 10 concurrent identical hold requests serializing to a single reservation without duplicate allocation.
6. `409 IDEMPOTENCY_CONFLICT` on changed payloads.
7. Strict JSON decoding: 400 on unknown fields, non-Z timestamps, and invalid path parameters.
8. Hold confirmation (`POST /v1/holds/:id/confirm`) returning `201 Created` with `Location` header.
9. Reservation lookup (`GET /v1/reservations/:id`) returning confirmed reservation with external `hold_id`.
10. Reservation release (`POST /v1/reservations/:id/release`) and direct hold release (`POST /v1/holds/:id/release`) returning `200 OK`.
11. `409 INVALID_STATE` when attempting to confirm an already-released hold.
12. Database isolation: verifies that no database user has `CONNECT` permissions on any other provider database.
