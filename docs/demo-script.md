# Batam MedHub — Demo Script

A repeatable **3–5 minute** demo that hits every rubric criterion. It can run
against the deployed endpoint (`https://api.bayumaulana.my.id`) or the local
stack, and it still works if the AI model is offline (the backend falls back to
the deterministic extractor).

> Full copy-paste `curl` commands for each step live in
> [`backend/README.md`](../backend/README.md). The Flutter app wraps these same
> calls behind a chat UI.

---

## Setup (before the clock starts)

- Have the mobile app open on the **login screen** (or the backend reachable).
- Run one optional reset so the demo state is golden:
  `POST /v1/demo/reset` with `{"scenario":"DEFAULT","confirm":true}` and the
  demo key header.

---

## Act 1 — The problem & the patient (≈45s)

**On screen:** a Singapore-based patient who needs a health screening in Batam.

**Say:**

> "Batam's new Health Tourism SEZ has world-class hospitals one ferry ride from
> Singapore — but booking a hospital, ferry, transfer, and hotel still means
> four separate systems, two currencies, two time zones, and endless replanning
> when something changes. Batam MedHub turns one sentence into one confirmed,
> recoverable journey."

Register/log in as the demo patient.

**Rubric hook:** *Problem Understanding & Relevance.*

---

## Act 2 — Natural language → structured intent (≈60s)

**On screen:** the chat screen.

**Say (or type):**

> "I need a basic medical check-up in Batam on 22 August 2026 for 1 person."

**Show:**

1. The request is sent to the backend.
2. **Cloudflare Workers AI (Llama 3.3 70B)** extracts a structured intent.
3. The backend validates it — and either asks one focused clarification or
   matches it to a catalog service (`MCU_BASIC`).
4. Explain the trust boundary: *"The model never decides anything medical or
   books anything — it only understands language. The orchestrator validates
   everything."*

**Rubric hooks:** *Technical Execution (AI patterns, tool-calling), Innovation.*

---

## Act 3 — Planning: proving the pieces fit (≈60s)

**On screen:** two plan options with an itemized timeline.

**Say:**

> "The hospital appointment is the anchor. The planner checks hard constraints
> before ranking: ferry arrival + immigration buffer + transfer + medical
> buffer must all fit before the appointment; the return ferry must fit after;
> capacity and accessibility must pass. That's why an 08:00 slot is impossible
> but 10:00 works — this is orchestration, not a search result."

**Show:** the rank-1 option's itemized timeline (hospital, outbound ferry,
transfer, return ferry, hotel when overnight) with per-leg prices in SGD.

**Rubric hooks:** *Technical Execution (determinism), Innovation (constraint
engineering).*

---

## Act 4 — Booking: one journey, four providers (≈45s)

**On screen:** "Book this journey" → booking status → **active itinerary v1**.

**Say:**

> "Approving runs an orchestrated saga: hold hospital, hold ferry, hold
> transport, hold hotel, then confirm each — with compensation if anything
> fails, and idempotency keys so a double-tap can't double-book."

**Show:** the active itinerary with booking reference and per-leg status chips.

**Rubric hooks:** *Technical Execution (saga, idempotency, robustness).*

---

## Act 5 — Disruption & recovery (the wow moment) (≈75s)

**On screen:** the active itinerary, then a disruption arrives.

**Say:**

> "An hour later, the hospital sends an event: the doctor requests 90 minutes
> of additional observation. That invalidates the original return ferry. The
> system doesn't ask the patient to re-plan from scratch — it computes what's
> impacted, generates at most two recovery options with price and time deltas,
> asks for approval on logistics only, holds the replacements, confirms them,
> and only then releases the superseded reservations."

**Show:** recovery option → approve → itinerary **v2** active, v1 `SUPERSEDED`
(added / changed / removed items visible).

**Rubric hooks:** *Innovation, Impact, Technical Execution.*

---

## Act 6 — Close (≈30s)

**Say:**

> "Deployable: the whole stack runs on one VPS behind HTTPS at
> `api.bayumaulana.my.id`. Safe by design: the AI can't diagnose or invent
> availability, money and time are exact, and every journey change is
> versioned and auditable. One request, one journey, one recovery."

**Rubric hooks:** *Impact & Feasibility, Presentation.*

---

## Fallback plan (if AI / network is offline)

- The backend automatically uses the **deterministic rule extractor**, so the
  demo still produces a structured intent.
- Keep a pre-recorded clip of the AI path as backup if needed.

## Common pitfalls to avoid

- Don't claim clinical superiority — say "best fit for your travel
  preferences."
- Don't say the mocks are real providers — everything is visibly synthetic
  (`synthetic: true`, `source: MOCK`).
- Don't overclaim production compliance — the data is hackathon fixtures.
