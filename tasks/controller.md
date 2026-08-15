# Control-Plane Codex Brief

## Mission

Keep the backend, provider services, and mobile-facing contract aligned while two implementation workers run independently. The control plane is the only writer of shared API contracts.

## Write ownership

- `PROJECT_UNDERSTANDING.md`
- `docs/architecture/**`
- `specs/**`
- `tasks/**`
- root-level integration documentation when explicitly needed

Do not take over broad implementation inside `backend/**` or `providers/**` while their workers are active.

## Responsibilities

1. Freeze domain vocabulary, ownership, invariants, and state transitions.
2. Maintain `specs/openapi.yaml` and `specs/provider-openapi.yaml`.
3. Keep golden JSON examples aligned with both specifications.
4. Review contract-change requests from workers.
5. Reject direct database coupling between core and providers.
6. Check backend/provider implementations against contract behavior at integration gates.
7. Maintain a deterministic end-to-end demo sequence.

## Contract-change procedure

A worker must report:

- the affected operation or schema;
- the implementation blocker;
- the smallest proposed contract change;
- whether the change is breaking; and
- which examples and consumers are affected.

The controller then updates the appropriate OpenAPI document and examples, validates them, and tells both workers when to sync. Workers must not silently extend request or response payloads.

## Integration gates

### Gate 0 — Contract

- Both OpenAPI files validate.
- Core and provider database ownership is explicit.
- The itinerary-v1 and disruption-to-v2 flows can be traced using golden examples.

### Gate 1 — Runtime and persistence

- Core API and four provider APIs build and expose health endpoints.
- Core migrations run on an empty core database.
- One provider PostgreSQL server creates four logical databases with isolated credentials.
- All provider migrations and deterministic seeds succeed.

### Gate 2 — Provider protocol

- Every provider completes search, hold, idempotent retry, confirm, lookup, and release.
- Invalid integration credentials return `401`.
- Capacity, price, or state conflicts return `409`; expired offers or holds return `410`.
- A provider credential cannot access another provider database.

### Gate 3 — Journey v1

- A supported structured intent produces no more than two feasible plans.
- An impossible cross-border combination is rejected with deterministic reasons.
- Confirmation compensates prior holds when a later provider fails.
- An uncertain initial-confirmation compensation leaves the trip request in `MANUAL_REVIEW` and creates no journey.
- A successful confirmation creates immutable itinerary v1 for an `ACTIVE` journey.
- Full contents of active and superseded itinerary versions are queryable.

### Gate 4 — Language boundary

- Workers AI output is schema-validated and catalog-validated.
- Matched, clarification, unsupported-service, and out-of-scope results are deterministic after validation.
- Malformed model output fails safely without invoking planning or booking.

### Gate 5 — Recovery

- Provider identity comes from its secret, not request-body identity.
- Provider event type and target ownership are authorized against the authenticated provider and stated journey.
- Duplicate external events replay only for an identical body fingerprint; changed reuse returns `409`.
- Additional care includes hospital `instruction_reference` and books provider-only follow-up inventory through normal search, hold, and confirm operations.
- The active itinerary is impact-analysed and recovery options expose signed time and price deltas.
- Approved recovery creates itinerary v2 without deleting v1.

### Gate 6 — Demo

- One reset procedure restores all synthetic data.
- Setup and event-injection commands are documented.
- The complete three-minute scenario is repeatable.

## Definition of done

The repository can demonstrate a supported request, an unsupported service, a feasible confirmed journey, a provider-authored disruption, and a recovered itinerary while preserving provider boundaries and data history.
