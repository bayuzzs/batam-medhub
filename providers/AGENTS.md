# Provider Workstream Instructions

Read `../tasks/providers.md` completely before editing. You own `providers/**` only.

- Treat `../specs/**`, `../backend/**`, `../mobile/**`, and `../docs/architecture/**` as read-only.
- Implement `../specs/provider-openapi.yaml` exactly.
- Report contract blockers to the control-plane conversation; do not edit specifications.
- Run four standalone Go/Gin services with one provider PostgreSQL server and four logical databases.
- Use GORM and independent `golang-migrate` SQL histories without runtime AutoMigrate.
- Do not build provider UI or disruption callbacks.
- Run the completion checks required by each phase in `../tasks/providers.md`.
