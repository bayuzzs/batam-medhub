# Batam MedHub

Batam MedHub orchestrates a medical tourist's planned journey from Singapore to Batam across hospital appointments, ferry travel, internal transport, and accommodation. It creates a feasible cross-provider itinerary and coordinates recovery when a provider disruption invalidates part of that journey.

> Project status: core OpenAPI/architecture v0.2 and provider OpenAPI v0.1 are defined. Core/provider implementation has not started.

## Why it exists

Batam's International Health Tourism Special Economic Zone can provide world-class clinical infrastructure while the patient experience remains fragmented across independent providers, booking systems, currencies, time zones, and schedule changes. Batam MedHub focuses on the orchestration gap: one request becomes one confirmable and recoverable journey.

Cloudflare Workers AI extracts structured intent from the patient's prompt. The Go backend treats that output as untrusted, validates it against the active medical-service catalog, and performs deterministic planning and booking. AI does not diagnose, choose treatment, invent availability, or confirm reservations.

## Repository

| Area | Responsibility | Owner |
|---|---|---|
| `mobile/` | Patient application | Mobile teammate; out of scope for these workstreams |
| `backend/` | Patient auth/profile, core API, intent validation, planning, booking saga, itinerary versions, and disruption recovery | Backend worker |
| `providers/` | Four headless provider APIs: hospital, ferry, hotel, and internal transport | Provider worker |
| `specs/` | Core and provider OpenAPI contracts plus golden examples | Control plane |
| `docs/architecture/` | Domain model, state machines, and logical ERDs | Control plane |
| `tasks/` | Non-overlapping Codex execution briefs and integration gates | Control plane |

The four provider applications use one PostgreSQL server with four isolated logical databases and credentials: `hospital_db`, `ferry_db`, `hotel_db`, and `transport_db`. The core backend uses a separate PostgreSQL datastore. Services communicate only through HTTP.

## Contract and architecture

- [Project understanding](PROJECT_UNDERSTANDING.md)
- [Domain model](docs/architecture/domain-model.md)
- [State machines](docs/architecture/state-machines.md)
- [Logical ERD](docs/architecture/erd.md)
- [Core API](specs/openapi.yaml)
- [Provider integration API](specs/provider-openapi.yaml)
- [Golden payloads](specs/examples/)

Validate both OpenAPI contracts with:

```bash
bash specs/validate.sh
```

## Implementation workstreams

Keep one control-plane conversation for shared contracts and run two implementation workers from separate Git worktrees:

- [Backend worker brief](tasks/backend.md)
- [Provider worker brief](tasks/providers.md)
- [Worktree and integration guide](tasks/README.md)

Commit this contract baseline before creating the worktrees so every worker starts from the same source of truth.

## Contracted hackathon slice

- Backend-owned registration, login, rotating refresh sessions, logout, and patient profile currency preference.
- Catalog-driven planned medical services with explicit unsupported and out-of-scope results.
- Natural-language prompt to validated structured intent.
- Deterministic appointment-anchored planning with cross-border time and capacity constraints.
- Static reference-currency conversion while preserving source prices.
- Idempotent provider search, hold, confirm, lookup, and release operations.
- Immutable itinerary versions and provider-authenticated disruption recovery.
- Persisted synthetic provider data and a deterministic demo reset.

Password recovery, email verification, MFA, social login, medical records, post-care, multilingual product behavior, payments, real providers, and production compliance are later-phase work.

## Synthetic-data notice

All providers, patients, prices, schedules, availability, instructions, and reservation references in this repository are fictional hackathon data. Runtime examples must preserve `synthetic: true` and `source: MOCK` markers.
