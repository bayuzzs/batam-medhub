# Batam MedHub API Contracts

Status: core implementation contract v0.2; provider implementation contract v0.1

## Contracts

- `openapi.yaml` defines backend-owned patient auth/profile operations, the patient-facing orchestration API, and the provider-authenticated disruption-ingestion operation.
- `provider-openapi.yaml` defines the backend-to-provider protocol implemented by the hospital, ferry, hotel, and internal transport services.
- `examples/core/` and `examples/provider/` contain payload-only golden examples referenced by the contracts.

Both contracts use OpenAPI 3.1. Money is represented in integer minor units using each currency's ISO 4217 exponent and code. Schedule instants are UTC and retain IANA time zones. All example data is synthetic.

The core auth contract uses email/password registration, short-lived HS256 access JWTs, rotating opaque refresh tokens, and stored session revocation. The mobile client must keep tokens in secure storage, never log them, send both current credentials for a profile update, and replace both tokens returned by that update.

## Ownership

The control-plane conversation is the only writer of `specs/**`. Backend, provider, and mobile workers consume these files as read-only contracts. A worker must report a contract-change request before implementing behavior that is not represented here.

## Validation

```bash
bash specs/validate.sh
```

The script first lints both original contracts, then inlines every local
`externalValue` JSON fixture into temporary OpenAPI documents and lints those.
This second pass proves that golden payloads conform to their referenced
schemas; a normal OpenAPI lint pass checks only that external files resolve.

The recommended linter may warn that a repository license has not yet been declared and that health operations do not contain artificial `4XX` responses. Those warnings do not invalidate the contracts.
