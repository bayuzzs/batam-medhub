# Backend Workstream Instructions

Read `../tasks/backend.md` completely before editing. You own `backend/**` only.

- Treat `../specs/**`, `../providers/**`, `../mobile/**`, and `../docs/architecture/**` as read-only.
- Implement the core API exactly as specified in `../specs/openapi.yaml`.
- Consume provider operations exactly as specified in `../specs/provider-openapi.yaml`.
- Report contract blockers to the control-plane conversation; do not edit specifications.
- Use Gin, GORM, PostgreSQL, and `golang-migrate` SQL migrations without runtime AutoMigrate.
- Run the completion checks required by each phase in `../tasks/backend.md`.
