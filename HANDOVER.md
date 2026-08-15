# Batam MedHub Control-Plane Handover

Last updated: 2026-08-15 (Asia/Jakarta)

## Purpose

This document is the resume point for the next control-plane agent. The control
plane reviews backend and provider checkpoints, owns shared contracts, and only
integrates worker branches after the relevant gate passes.

Product code, API names, fixtures, user-facing copy, and repository
documentation must remain in English. The project owner may converse in English
or Indonesian.

## Read First

Read these files completely before making decisions or edits:

1. `AGENTS.md`
2. `PROJECT_UNDERSTANDING.md`
3. `docs/architecture/domain-model.md`
4. `docs/architecture/state-machines.md`
5. `docs/architecture/erd.md`
6. `specs/openapi.yaml`
7. `specs/provider-openapi.yaml`
8. `tasks/controller.md`
9. The workstream brief and nested `AGENTS.md` for any code being reviewed

The research repository at `/home/bayuzzs/Documents/bsh-research` was supplied
only as background. Do not modify it.

## Ownership Boundaries

- `mobile/**`: owned by the user's partner. Treat it as read-only unless the
  user explicitly assigns mobile work.
- `backend/**`: backend worker implementation scope.
- `providers/**`: provider worker implementation scope.
- `specs/**`, `docs/architecture/**`, `tasks/**`: control-plane scope.
- Workers must request contract changes; they must not edit shared contracts.
- Core and providers communicate only through HTTP. They never read each
  other's databases.

### Current Git Snapshot

The snapshot below was verified immediately after the Backend B10 integration.

| Worktree | Branch | HEAD | State |
|---|---|---|---|
| `/home/bayuzzs/Project/batam-medhub` | `main` | `46fb814fe0cb20e98d98d281db24608c2a93c7f9` | Clean and synchronized with `origin/main` |
| `/home/bayuzzs/Project/batam-medhub-backend` | `feat/backend-b10` | `46fb814fe0cb20e98d98d281db24608c2a93c7f9` | Integrated into `main` (Phases B1–B10 Complete) |
| `/home/bayuzzs/Project/batam-medhub-providers` | `feat/providers-p7-p8` | `2e3a8e4d92f0d1973ab64b1ccf98a705a9283a1f` | Integrated into `main` (Phases P1–P8 Complete) |

Already integrated and pushed commits on `main`:

- `a28e5ee` — backend runtime and persistence foundation (Gate 1)
- `4c16914` — provider Gate 1 foundation
- `b6427da` — immutable itinerary-history trigger fix
- `45e21e6`, `e0804f2`, `885a6b0` — backend B3 trust-boundary primitives & schema checks (Gate 2)
- `a707d80`, `d7b4b36` — provider P3 hospital reference implementation & strict validation (Gate 3 Hospital)
- `2cd30a7` — combined B3+P3 integration merge commit
- `297becf` — provider P4 ferry reference service implementation (Gate 3 Ferry)
- `9167e44` — backend B4 structured-intent & trip-request planning flow (Gate 4 Foundation)
- `cd2c5ed` — backend B4 integration merge commit
- `ee22029` — backend B5 provider HTTP client adapters & multi-provider aggregator
- `6d88d70` — provider P5 hotel reference service implementation (Gate 3 Hotel)
- `4667048` — combined B5+P5 integration merge commit
- `fb146bc` — backend B6 package planning and ranking engine
- `11b5ad6` — provider P6 transport reference service implementation (Gate 3 Transport)
- `41913a8` — combined B6+P6 integration merge commit
- `2e3a8e4` — provider P7/P8 smoke verification suite, deterministic reset, and documentation
- `a9c8584` — backend B7 multi-provider booking saga, compensations, and journey tracking
- `5993305` — backend B7 integration merge commit
- `7bf639c`, `e6cf0f4` — backend B8 Cloudflare Workers AI intent extraction with guardrails and offline fallback
- `27da600`, `9ca8e92`, `2310aec` — backend B9 disruption & recovery engine with versioned itinerary activation
- `4f7fcb2`, `12df843`, `46fb814` — backend B10 demo hardening, reset endpoint, documentation, and env template

## Completed Baseline

- Gate 0 architecture and two OpenAPI contracts are present and valid.
- Gate 1 is integrated on `main` (backend and provider runtime foundation, PostgreSQL schemas, deterministic seeds).
- Gate 2 is integrated on `main` (patient auth, sessions, profile, rate-limiting, FX conversion, catalog).
- Gate 3 (Provider Platform) is **100% COMPLETE & VERIFIED** on `main` across all 4 mock providers (Hospital, Ferry, Hotel, Transport).
- Gate 4 (Trip Planning & Ranking Engine) is **100% COMPLETE & INTEGRATED** on `main`.
- Gate 5 (Multi-Provider Booking Saga & Itinerary v1) is **100% COMPLETE & INTEGRATED** on `main`.
- Gate 6 (AI Intent Extraction) is **100% COMPLETE & INTEGRATED** on `main`.
- Gate 7 (Disruption & Recovery Engine) is **100% COMPLETE & INTEGRATED** on `main`.
- Gate 8 (Demo Hardening & Documentation) is **100% COMPLETE & INTEGRATED** on `main`:
  - `POST /v1/demo/reset` with authentication (`X-Demo-Secret`), dynamic table truncate cascade, reference preservation, and golden registration replay;
  - Comprehensive English `backend/README.md` and `backend/.env.example`;
  - Repeatable end-to-end verification suite in `backend/cmd/verify/main.go`.

## Status of Workstreams

### 1. Backend Workstream: 100% COMPLETE

All assigned phases in `tasks/backend.md` (B1 through B10) are 100% complete, verified, documented, and merged to `main`.

### 2. Provider Platform Workstream: 100% COMPLETE

All assigned phases in `tasks/providers.md` (P1 through P8) are 100% complete, verified, documented, and merged to `main`.

### 3. Mobile Workstream

The mobile application (`mobile/**`) is ready for integration against the complete core backend API (`http://localhost:8080`) and provider mock services.

## Non-Negotiable Decisions

- Go, Gin, GORM, PostgreSQL, and `golang-migrate`; never runtime AutoMigrate.
- One provider PostgreSQL server with four logical databases and isolated
  credentials.
- Core PostgreSQL is separate from all provider databases.
- Patient authentication uses short-lived HS256 access JWTs and rotating opaque
  refresh tokens whose SHA-256 hashes alone are stored.
- PostgreSQL uniqueness and row locks are authoritative for concurrent state.
- Money uses integer ISO 4217 minor units and exact arithmetic, never binary
  floating point.
- IDR uses ISO 4217 minor-unit exponent 2 in this contract. Do not change the
  existing IDR fixtures to exponent 0.
- Static FX metadata is returned with both source and display money.
- Cloudflare Workers AI will be called directly by the backend in B8. Model
  output remains untrusted and cannot diagnose, select treatment, invent
  provider facts, plan constraints, or book resources.
- Provider disruptions are later submitted manually to the core provider-auth
  route; no provider dashboard is required.
- Automated test suites remain deferred by owner decision, but focused checks,
  builds, migrations, contract validation, and deterministic smoke flows are
  required at every checkpoint.

## Known Repository Note

`mobile/.env` was already tracked as an empty file by the mobile workstream.
Do not put secrets in it and do not modify mobile-owned files as part of backend
or provider integration.
