# Batam MedHub — Documentation

This folder is the single source of documentation for the Batam MedHub
hackathon submission. It is written for judges, reviewers, and future
maintainers: every page explains **what** the project does, **why** it matters
for Batam and Singapore, **how** it is built, and **how** to demo it.

> New here? Start with the [README at the repository root](../readme.md) for a
> 2-minute summary, then come back to this index for the full story.

---

## Quick navigation

| Document | What it covers |
| :--- | :--- |
| [Hackathon submission](hackathon.md) | Problem understanding, relevance to Batam ↔ Singapore, innovation, impact, feasibility, presentation plan — mapped to the judging rubric. |
| [Architecture](architecture.md) | System context, components, request flows, booking saga, disruption & recovery, state machines, data model, and engineering trade-offs. |
| [API documentation](api.md) | Human-readable tour of the core and provider APIs, with links to the OpenAPI contracts and golden examples. |
| [Build process & AI](build-process.md) | How the project was built: research & design before the hackathon, full implementation during it, and the AI tooling used. |
| [Demo script](demo-script.md) | A repeatable 3–5 minute demo: register → request → plan → book → disrupt → recover, with talking points per judging criterion. |

---

## Reference documentation (kept in the repo)

- [Project understanding](../PROJECT_UNDERSTANDING.md) — full design decisions and assumptions behind the product.
- [Domain model](architecture/domain-model.md) — entities, bounded contexts, and ownership rules.
- [State machines](architecture/state-machines.md) — trip planning and disruption lifecycles.
- [Logical ERD](architecture/erd.md) — the core and provider data model.
- [Core API contract](../specs/openapi.yaml) — the patient-facing backend API (OpenAPI 3.1).
- [Provider API contract](../specs/provider-openapi.yaml) — the backend ↔ provider protocol (OpenAPI 3.1).
- [Golden payloads](../specs/examples/) — schema-validated example requests and responses.
- [Backend guide](../backend/README.md) — local setup, environment variables, and end-to-end smoke tests.
- [Deployment guide](../deploy/README.md) — Docker Compose + nginx + HTTPS deployment.

---

## One-paragraph summary

**Batam MedHub** turns one medical-trip request into one feasible, confirmable,
and recoverable cross-provider journey for a medical tourist travelling from
**Singapore to Batam**. It orchestrates a hospital appointment together with
ferry travel, internal transport, and accommodation — checking real constraints
such as ferry check-in cutoffs, immigration buffers, appointment times,
capacity, accessibility, and currency — and when any provider disrupts the
confirmed journey (a ferry delay, a doctor's follow-up request), it
automatically replans and recovers the entire itinerary.

The system is built from four mock **provider services** (hospital, ferry,
hotel, transport), a Go **core orchestrator** that owns planning, booking, and
recovery, a **Flutter mobile app** for the patient, and an **AI intent
extractor** (Cloudflare Workers AI) that converts natural language into
validated structured intent. The AI never diagnoses, never invents
availability, and never books anything by itself — the orchestrator stays
deterministic and explainable.
