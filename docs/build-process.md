# Batam MedHub — Build Process & AI Tooling

This page is an honest account of **how** Batam MedHub was built and **which AI
tools** were used. Two important facts:

1. **Research & design was executed before the hackathon.**
2. **The implementation was executed fully during the hackathon.**

---

## 1. Phase A — Research & design (before the hackathon)

Before the hackathon, the team did the thinking that makes the demo precise:

- Studied the challenge brief and the Batam/Singapore context (Batam's
  International Health Tourism Special Economic Zone, Singapore medical-tourist
  demand, ferry/immigration constraints, currency and time-zone differences).
- Explored the fragmented patient experience across booking, records, travel,
  post-care, and language boundaries, and **narrowed the scope** to the highest
  value, most provable slice: **cross-provider journey orchestration**.
- Decided on the architecture: a Go core orchestrator, four HTTP provider mocks,
  PostgreSQL persistence, a Flutter client, and AI used only for structured
  intent extraction.
- Fixed non-negotiables: contract-first OpenAPI, no shared databases, exact
  money/time handling, deterministic planning, AI trust boundaries, and
  synthetic-only data.
- Wrote the architecture documentation that is now in `docs/architecture/`.

Everything below that is **code** was built during the hackathon.

---

## 2. Phase B — Implementation (during the hackathon)

The implementation ran as **three parallel workstreams** coordinated by a
control plane, with continuous integration onto `main`:

| Workstream | Scope | Shipped as |
| :--- | :--- | :--- |
| **Backend worker** | Auth & sessions, intent validation, provider adapters, deterministic planner & ranking, booking saga, itinerary versioning, disruption & recovery engine, demo reset. | Phases B1–B10, all merged & verified |
| **Provider worker** | Four standalone Go services (hospital, ferry, hotel, transport) with isolated PostgreSQL DBs, deterministic seeds, smoke suite, and reset. | Phases P1–P8, all merged & verified |
| **Mobile worker** | Flutter patient app: auth, chat-based trip request, plan detail, booking, active itinerary, profile, history. | Features merged via the `home-screen` branch |
| **Control plane** | Shared OpenAPI contracts, golden payloads, architecture docs, integration gates, and this documentation. | Continuous |

Key gates that were completed end-to-end:

- **Gate 0** — architecture + both OpenAPI contracts valid.
- **Gates 1–2** — backend runtime/persistence and provider platform foundation;
  patient auth, sessions, profile, FX conversion, catalog.
- **Gate 3** — all four provider services complete and verified.
- **Gate 4** — trip planning & ranking engine.
- **Gate 5** — multi-provider booking saga & itinerary v1.
- **Gate 6** — AI intent extraction with guardrails and offline fallback.
- **Gate 7** — disruption & recovery engine (itinerary v2).
- **Gate 8** — demo hardening, `POST /v1/demo/reset`, documentation, deployment.

The project is also **deployed**: the full stack runs on a VPS behind nginx +
HTTPS at `https://api.bayumaulana.my.id` (see `deploy/README.md`).

---

## 3. AI used during development

AI agents were used as development accelerators across all workstreams — for
architecture, code generation, contract design, test writing, review, and
documentation. The primary AI tooling during development was:

- **GPT 5.6**
- **OpenCode Big Pickle**
- **Gemini 3.7 High Reasoning**

These were used to design the OpenAPI contracts, generate and review Go backend
and provider code, build the Flutter UI, write deterministic tests, and produce
this documentation — while the team reviewed and integrated every change.

> Note the distinction: the AI used **during development** (above) is not the
> same as the AI model that runs **inside the product**. See the next section.

---

## 4. AI used inside the product (read from `.env`)

The product's runtime AI is configured through the backend environment
(`backend/.env.example`, mirrored in `deploy/.env.example`):

```dotenv
# Cloudflare Workers AI Intent Extractor
CLOUDFLARE_ACCOUNT_ID=<account>
CLOUDFLARE_API_TOKEN=<token>
CLOUDFLARE_AI_MODEL=@cf/meta/llama-3.3-70b-instruct-fp8-fast
CLOUDFLARE_AI_BASE_URL=https://api.cloudflare.com/client/v4/accounts/<ACCOUNT_ID>/ai/v1
CLOUDFLARE_AI_TIMEOUT_SECONDS=15
```

**The production model is `@cf/meta/llama-3.3-70b-instruct-fp8-fast`
(Llama 3.3 70B Instruct, FP8 fast variant) served on Cloudflare Workers AI**,
called through the OpenAI-compatible `/ai/v1` gateway (`/chat/completions`).

Model usage in the product is deliberately narrow:

- Converts the patient's natural-language trip request into a **structured
  intent** (service, date window, patient count, stay type, budget,
  preferences, missing fields).
- The backend treats the output as **untrusted**: it validates the schema,
  enums, types, dates, and ranges; verifies the service code against the active
  PostgreSQL catalog; and applies clinical guardrails (no diagnosis, triage,
  treatment selection, or emergency planning).
- If the model or network is unavailable, a **deterministic rule-based
  extractor** keeps the demo fully usable offline.

The model never invents availability, never plans constraints, and never books —
all of that stays deterministic in the Go orchestrator.
