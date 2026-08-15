# Batam MedHub — Hackathon Submission

This document is the narrative behind the submission, mapped to the official
judging rubric. It is intended to make the judges' job easy: each section below
answers one rubric criterion directly.

**One-line pitch:** *One medical-trip request becomes one feasible, confirmable,
and recoverable cross-provider journey for a medical tourist travelling from
Singapore to Batam.*

---

## 1. Problem Understanding & Relevance — 20%

### The Batam-specific challenge

Batam's new **International Health Tourism Special Economic Zone** — backed by
multi-billion-rupiah investment and world-class hospital ambitions — is designed
to attract **Singaporean medical tourists**. The clinical infrastructure is
coming; what is missing is the **patient experience glue** across independent
providers, booking systems, currencies, time zones, and schedule changes.

We identified the core failure the challenge implies:

> A marketplace can show that a hospital slot, a ferry, a transfer, and a hotel
> each *exist*. Batam MedHub proves they **fit together** — under real
> constraints such as ferry check-in cutoffs, immigration buffers, appointment
> times, capacity, accessibility, offer expiry, and cross-border time zones.

Concrete example encoded in our demo data: a hospital 08:00 slot is infeasible
because the earliest ferry arrival plus immigration, transfer, and medical
safety buffers puts the patient at the hospital at 08:45. A 10:00 appointment is
the earliest feasible anchor. That single example shows we understand the
**orchestration gap**, not just the booking gap.

### Cross-border relevance (Batam ↔ Singapore)

- A ~1-hour ferry connects Singapore (HarbourFront) to Batam (Batam Centre) —
  same-day medical visits are realistic.
- Currency differences (SGD/IDR), time-zone differences (WIB/SGT), and
  immigration/ferry cutoffs make cross-border journeys genuinely harder to
  coordinate than domestic ones — all modeled.
- Our recovery demo (a doctor's follow-up request that invalidates the return
  ferry) is a real scenario a same-day medical tourist can face.

---

## 2. Technical Execution & Engineering Quality — 25%

### Recommended tech stack, used properly

| Layer | Technology |
| :--- | :--- |
| Mobile | **Flutter** (Dart) — Riverpod, GoRouter, Dio, freezed |
| Core backend | **Go** — Gin, GORM, PostgreSQL, `golang-migrate` |
| Provider mocks | **Go** — four standalone services sharing one module |
| Persistence | **PostgreSQL** — 1 core DB + 4 isolated provider DBs |
| AI | **Cloudflare Workers AI** — Llama 3.3 70B (OpenAI-compatible gateway) |
| Contracts | **OpenAPI 3.1** (core + provider), golden JSON examples, lint/validate |
| Deployment | **Docker Compose**, nginx reverse proxy, Let's Encrypt HTTPS |

### Engineering quality highlights

- **Contract-first development.** Two OpenAPI contracts are the source of
  truth; golden payloads are linted against them (`bash specs/validate.sh`).
- **Clean separation of ownership.** Core, providers, and mobile are separate
  workstreams that communicate only over HTTP; no service reads another's
  database.
- **Production-aware auth.** Argon2id hashing, short-lived HS256 access JWTs,
  rotating opaque refresh tokens (only hashes stored), session revocation,
  idempotent logout, rate limiting.
- **Distributed-consistency patterns.** An orchestration-based **booking saga**
  with holds → confirms → compensation, plus **idempotency keys** on every
  mutation — double-clicks and retries are safe.
- **Immutable itinerary versioning.** v1 is never mutated; recovery activates
  v2 and marks v1 `SUPERSEDED`.
- **Deterministic, explainable planning.** Hard constraints filter before
  scoring; at most two plan options; explicit "best fit for your travel
  preferences" language (never clinical superiority claims).
- **AI trust boundary.** Model output is treated as untrusted: schema-validated,
  enum-checked, catalog-verified, guardrailed (no diagnosis/triage/treatment),
  with a deterministic offline fallback.
- **Exact money & time.** Integer minor units + ISO currency; UTC instants with
  IANA zones.
- **Verification.** Builds, empty-database migrations, contract validation,
  health checks, and repeatable end-to-end smoke flows (`backend/cmd/verify`,
  `providers/smoke.sh`); Flutter widget tests cover routing, chat state, plan
  detail, and the active itinerary.

### Agent patterns & tool-calling

We used **AI agents as development accelerators** across three parallel
workstreams (see [build-process.md](build-process.md)), and in the product the
LLM performs **structured intent extraction** — the safe, high-value "tool
call" — while every decision that affects the patient remains deterministic.

### Robustness & production awareness

- Timeouts and independent failure isolation on provider calls (one provider
  failing doesn't erase another's results).
- CORS, health (`/healthz`) and readiness (`/readyz`) endpoints, Dockerized
  deployment with per-provider least-privilege DB roles.
- A real deployed endpoint: `https://api.bayumaulana.my.id` (see `deploy/`).

---

## 3. Innovation & Creativity — 20%

- **The orchestration layer is the product.** Not an AI doctor, not a
  marketplace — a *journey orchestrator* that owns cross-provider feasibility
  and recovery. This is a novel position in the medical-tourism stack.
- **Disruption-driven recovery.** The system doesn't just book; it continuously
  protects the journey. One generic pipeline handles hospital, ferry, hotel,
  and transport disruptions and replans the whole trip — including
  provider-requested additional care.
- **Explainable AI integration.** The LLM is used exactly where language
  understanding adds value (free-text → structured intent) and explicitly kept
  away from clinical and booking decisions — a defensible, safe AI pattern.
- **Real constraint engineering.** Ferry cutoffs + immigration buffers +
  transfer time + medical arrival buffers + cross-border time zones are encoded
  as hard constraints; the demo proves feasibility mathematically, not by
  ranking.
- **Cross-border currency & timezone correctness.** Exact money arithmetic and
  explicit IANA zones show genuine attention to the Batam/Singapore context.
- **Offline-resilient demo.** A deterministic fallback keeps the presentation
  running even if the model or network is unavailable.

---

## 4. Impact & Feasibility — 20%

### Real-world value

- Addresses a real, funded policy push (Batam's health-tourism SEZ) and a real
  friction: fragmented cross-border medical journeys.
- Reduces the cognitive and coordination burden on patients and their
  caregivers — one conversation instead of four bookings and constant replanning.
- Gives hospitals a practical way to attract international patients without
  owning travel logistics.
- A trustworthy recovery flow de-risks same-day cross-border medical visits.

### Scalability & deployability

- Providers are clean HTTP boundaries, so real hospital/ferry/hotel systems can
  replace the mocks with the same contract.
- The core is a modular monolith that can grow; the saga pattern extends to
  more providers and async transports.
- The full stack runs on a single VPS with Docker Compose + nginx + HTTPS and
  is already deployed at `api.bayumaulana.my.id`.
- Deterministic seeds and a one-shot demo reset make the system repeatable and
  sustainable for judging and future demos.

### Sustainability & trust

- The AI guardrails (no diagnosis, no invented availability, human approval for
  every booking and every recovery) are exactly what a production medical
  product needs — this is a credible path to real adoption, not a toy.

---

## 5. Presentation & Demo — 15%

A repeatable **3–5 minute demo** is documented in [demo-script.md](demo-script.md).
In short:

1. Reset demo state → register a synthetic patient.
2. Type a natural-language trip request (e.g. *"basic check-up in Batam on 22
   August for 1 person"*).
3. Watch AI intent extraction → clarification → **two plan options**.
4. Approve a plan → the multi-provider **booking saga** confirms hospital,
   ferry, transfer, hotel → active itinerary v1.
5. Submit a **hospital disruption** (additional care) → recovery options → approve
   → itinerary **v2** activated, v1 superseded.

Talking points cover the rubric explicitly, and the team can explain trade-offs
(modular monolith vs microservices, deterministic planner vs agent, sync saga vs
message broker).

---

## 6. Team collaboration

The build ran as **three parallel workstreams** coordinated by a control plane
(see [build-process.md](build-process.md)):

- **Backend worker** — auth, intent, planning, booking saga, disruption engine.
- **Provider worker** — hospital, ferry, hotel, transport services.
- **Mobile worker** — Flutter patient app.
- **Control plane** — shared contracts (`specs/**`), architecture, integration.

Each workstream shipped in phases (backend B1–B10, providers P1–P8, mobile
features) with integration gates, and the branches were merged into `main`
continuously.
