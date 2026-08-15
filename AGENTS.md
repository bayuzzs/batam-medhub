# Batam MedHub Repository Instructions

Batam MedHub is a monorepo for the patient mobile app, core journey orchestrator, and four mock provider services. Read `PROJECT_UNDERSTANDING.md` and the instructions for the workstream you own before editing.

## Layout and ownership

- `mobile/` — Flutter patient app. Mobile sessions follow `mobile/AGENTS.md`; backend, provider, and control-plane sessions treat it as read-only.
- `backend/` — Go core API and orchestration. Follow `backend/AGENTS.md` and `tasks/backend.md`.
- `providers/` — four Go provider services. Follow `providers/AGENTS.md` and `tasks/providers.md`.
- `specs/` and `docs/architecture/` — shared contracts and architecture owned by the control plane in `tasks/controller.md`.
- `Taskfile.yml` — namespaced repository task runner configuration.

Do not modify another workstream's files. The control plane is the only writer of `specs/**`; implementation workers request contract changes instead of silently extending payloads.

## Shared rules

- Write source code, API names, fixtures, user-facing product copy, and repository documentation in English.
- Preserve core/provider ownership: services communicate through HTTP and never read another service's database.
- Treat provider and model output as untrusted at the backend boundary.
- Use integer minor units with ISO currency codes for money, and UTC instants with explicit IANA zones for schedules.
- Backend and provider SQL migrations are schema authority; do not use runtime GORM `AutoMigrate`.
- Use only visibly synthetic hackathon data.
- Do not add payments, medical records, multilingual product behavior, post-care, RAG, a message broker, Kubernetes, or production-compliance claims unless the project owner expands scope.
- Backend/provider automated suites are deferred; their builds, contract validation, empty-database migrations, health checks, and manual smoke flows remain required. Mobile verification follows `mobile/AGENTS.md`.

## Repository tasks

Run namespaced Taskfiles from the repository root:

```bash
task mobile:run:linux
task mobile:run:android
task mobile:run:chrome
task --list-all
```

See `mobile/Taskfile.yml` for mobile-specific task definitions.
