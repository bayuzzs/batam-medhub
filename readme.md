# Batam MedHub 🏝️🩺

**One medical-trip request → one feasible, confirmable, and recoverable
cross-provider journey.**

Batam MedHub orchestrates a medical tourist's planned journey from **Singapore
to Batam** across hospital appointments, ferry travel, internal transport, and
accommodation — and keeps that confirmed journey synchronized when a provider
disruption (a ferry delay, a doctor's follow-up request) invalidates part of it.

> **Hackathon project.** The research & design phase was completed *before* the
> hackathon; the implementation was built *entirely during* it. The full stack
> is deployed at `https://api.bayumaulana.my.id`.

---

## The problem — why it matters (Batam ↔ Singapore)

Batam's new **International Health Tourism Special Economic Zone** brings
world-class hospitals within a one-hour ferry ride of Singapore. But the patient
experience is still fragmented across independent hospitals, ferry operators,
ground transport, hotels, currencies, time zones, and schedule changes.

A marketplace can show that each service **exists**. Batam MedHub proves they
**fit together** — under real constraints like ferry check-in cutoffs,
immigration buffers, appointment times, capacity, accessibility, and offer
expiry. Our demo encodes the proof directly: an 08:00 hospital slot is
infeasible because the earliest ferry arrival plus immigration, transfer, and
medical buffers reaches the hospital at 08:45 — the 10:00 appointment is the
earliest feasible anchor.

> We are **not** an AI doctor and **not** a booking marketplace. We are the
> orchestration layer that owns cross-provider feasibility and recovery.

---

## The solution — what we built

1. **Natural language → validated intent.** The patient types a request like
   *"basic check-up in Batam on 22 August for 1 person"*. Cloudflare Workers AI
   extracts a structured intent; the backend validates it against the active
   medical-service catalog and applies clinical guardrails.
2. **Deterministic planning.** The appointment anchors the journey. Hard
   constraints filter before ranking, and the planner returns **at most two
   explainable options** with an itemized timeline and per-leg pricing.
3. **One booking, four providers.** An orchestrated **booking saga** holds and
   confirms hospital, ferry, transport, and hotel with idempotency keys and
   automatic compensation.
4. **Immutable itinerary versions.** A confirmed journey is itinerary **v1**;
   it is never mutated.
5. **Disruption & recovery.** One generic pipeline ingests provider events,
   analyzes impact, generates recovery options, and activates itinerary **v2**
   (v1 stays `SUPERSEDED`).

---

## Models & AI

### AI running inside the product (from `.env`)

The backend `.env` configures a **Cloudflare Workers AI** intent extractor:

```dotenv
CLOUDFLARE_AI_MODEL=@cf/meta/llama-3.3-70b-instruct-fp8-fast
CLOUDFLARE_AI_BASE_URL=https://api.cloudflare.com/client/v4/accounts/<ACCOUNT_ID>/ai/v1
```

- **Production model:** **Llama 3.3 70B Instruct (FP8 fast)** on Cloudflare
  Workers AI (`@cf/meta/llama-3.3-70b-instruct-fp8-fast`), called via the
  OpenAI-compatible `/ai/v1` gateway.
- **Role:** converts natural-language trip requests into **structured intent**
  only.
- **Trust boundary:** output is treated as untrusted — schema-validated,
  enum-checked, catalog-verified, and guardrailed. The model never diagnoses,
  chooses treatment, invents availability, plans constraints, or books anything.
- **Resilience:** a deterministic rule-based extractor keeps the demo fully
  usable if the model or network is offline.

### AI used during development

The team used AI agents as development accelerators across all workstreams:

- **GPT 5.6**
- **OpenCode Big Pickle**
- **Gemini 3.7 High Reasoning**

See [docs/build-process.md](docs/build-process.md) for the full account of the
research & design phase (pre-hackathon) and the implementation phase
(during-hackathon).

---

## Architecture at a glance

```mermaid
flowchart LR
    MOBILE[Flutter Patient App]
    ACTORS[Provider Actors]
    AI[Cloudflare Workers AI]
    CORE_DB[(Core PostgreSQL)]
    subgraph PROVIDERS[providers/]
        HOSPITAL[Hospital] & FERRY[Ferry] & HOTEL[Hotel] & TRANSPORT[Transport]
        PD[(PostgreSQL: 4 isolated DBs)]
    end
    subgraph BACKEND[Go Core Orchestrator]
        API[Gin API] --> AUTH[Auth] & INTENT[Intent validation] & ADAPTERS[Adapters]
        INTENT --> PLANNER[Planner] --> BOOKING[Saga] --> JOURNEY[Itinerary versions]
        IMPACT[Disruption impact] --> PLANNER
    end
    MOBILE -->|HTTPS / JWT| API
    ACTORS -->|Provider key| API
    ADAPTERS <--> HOSPITAL & FERRY & HOTEL & TRANSPORT
    API --> AI
    API --> CORE_DB
```

- **Mobile** — Flutter app (Riverpod, GoRouter, Dio): auth, chat-based trip
  request, plan detail, active itinerary, profile, history.
- **Backend** — Go modular monolith (Gin, GORM, PostgreSQL): auth, intent
  validation, planning, booking saga, itinerary versioning, disruption recovery.
- **Providers** — four standalone Go services (hospital, ferry, hotel,
  transport) sharing one PostgreSQL server with **four isolated logical
  databases** and separate credentials. Services communicate **only over HTTP**.
- **Contracts** — two OpenAPI 3.1 contracts are the source of truth, with
  schema-validated golden payloads.

Full details: [docs/architecture.md](docs/architecture.md).

---

## Tech stack

| Layer | Technology |
| :--- | :--- |
| Mobile | Flutter (Dart) · Riverpod · GoRouter · Dio |
| Core backend | Go · Gin · GORM · PostgreSQL · golang-migrate |
| Provider mocks | Go (4 services, shared module) |
| Persistence | PostgreSQL — 1 core DB + 4 isolated provider DBs |
| AI | Cloudflare Workers AI · Llama 3.3 70B Instruct (FP8 fast) |
| Contracts | OpenAPI 3.1 (core + provider) + golden examples |
| Deployment | Docker Compose · nginx · Let's Encrypt HTTPS |

---

## Documentation

| Document | Contents |
| :--- | :--- |
| [docs/README.md](docs/README.md) | Documentation index — start here. |
| [docs/hackathon.md](docs/hackathon.md) | Submission narrative mapped to the judging rubric. |
| [docs/architecture.md](docs/architecture.md) | System architecture, flows, saga, recovery, trade-offs. |
| [docs/api.md](docs/api.md) | Human-readable API tour. |
| [docs/demo-script.md](docs/demo-script.md) | The 3–5 minute demo, step by step. |
| [docs/build-process.md](docs/build-process.md) | Research & design (pre-hackathon) and implementation (hackathon) + AI tooling. |
| [Project understanding](PROJECT_UNDERSTANDING.md) | Full design rationale. |
| [Domain model](docs/architecture/domain-model.md) · [State machines](docs/architecture/state-machines.md) · [ERD](docs/architecture/erd.md) | Deep architecture references. |
| [Backend guide](backend/README.md) | Setup, env vars, end-to-end smoke tests. |
| [Deployment guide](deploy/README.md) | Docker Compose + HTTPS deployment. |

---

## API documentation

- [Core API — `specs/openapi.yaml`](specs/openapi.yaml) — patient-facing backend
  API (auth, trip requests, plans, booking, itinerary, disruptions).
- [Provider API — `specs/provider-openapi.yaml`](specs/provider-openapi.yaml) —
  the backend ↔ provider protocol (search / hold / confirm / release).
- [Golden payloads — `specs/examples/`](specs/examples/) — schema-validated
  example requests and responses.
- [Human-readable API tour](docs/api.md).

Validate both contracts:

```bash
bash specs/validate.sh
```

---

## Quick start

```bash
# 1. Start provider infrastructure + PostgreSQL
docker compose -f ./providers/docker-compose.yml up -d

# 2. Apply core migrations
cd backend
export DATABASE_URL="postgres://provider_admin:provider_admin_dev_password@localhost:5432/core_db?sslmode=disable"
go run ./cmd/migrate up

# 3. Run the backend (optional: set Cloudflare credentials for the AI path)
export JWT_SIGNING_SECRET="12345678901234567890123456789012"
export DEMO_SECRET="demo_dev_secret"
go run ./cmd/api

# 4. Verify the whole flow
go run ./cmd/verify
```

Run the mobile app (from repo root):

```bash
task mobile:run:chrome      # or run:linux / run:android
```

See [backend/README.md](backend/README.md) and [deploy/README.md](deploy/README.md)
for full instructions.

---

## Repository layout

| Area | Responsibility |
| :--- | :--- |
| `mobile/` | Flutter patient application. |
| `backend/` | Go core orchestrator (auth, planning, saga, disruption recovery). |
| `providers/` | Four headless Go provider services. |
| `specs/` | OpenAPI contracts + golden examples. |
| `docs/` | Hackathon & architecture documentation. |
| `docs/architecture/` | Domain model, state machines, ERD. |
| `deploy/` | Production Docker Compose + nginx + HTTPS. |
| `tasks/` | Workstream briefs and integration gates. |

---

## Synthetic-data notice

All providers, patients, prices, schedules, availability, instructions, and
reservation references in this repository are **fictional hackathon data**.
Runtime examples must preserve `synthetic: true` and `source: MOCK` markers.
The mocks simulate real provider systems; they are not real operational data.

