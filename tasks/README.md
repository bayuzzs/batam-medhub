# Batam MedHub — Codex Workstreams

This directory contains the execution briefs for three non-overlapping roles:

| Role | Brief | Write ownership |
|---|---|---|
| Control plane | [`controller.md`](controller.md) | Contracts, architecture documents, integration decisions, and task briefs |
| Backend worker | [`backend.md`](backend.md) | `backend/**` |
| Provider worker | [`providers.md`](providers.md) | `providers/**` |

## Recommended execution model

Keep the current conversation as the control plane. Run the backend and provider workers in two separate Codex CLI or extension sessions, each in a separate Git worktree. Do not run two workers in the same checkout.

The repository-level and workstream-level `AGENTS.md` files automatically route each Codex session to the applicable brief when the session starts in the relevant worktree and directory.

Create the worktrees only after the contract pack has been reviewed and committed:

```bash
git worktree add ../batam-medhub-backend -b feat/backend
git worktree add ../batam-medhub-providers -b feat/providers
```

The backend worker is the implementation owner of `backend/**`, but it does not edit `specs/**`. The provider worker treats `specs/provider-openapi.yaml` as read-only. Contract changes are requested from the control-plane conversation and applied in one place before either worker continues with the affected operation.

Start each CLI session from its owned subtree so Codex loads both the root and nested instructions:

```bash
cd ../batam-medhub-backend
codex --cd backend

cd ../batam-medhub-providers
codex --cd providers
```

Paste the corresponding copy-paste prompt from `tasks/backend.md` or `tasks/providers.md` as the first task in that session.

## Integration order

1. Freeze and commit the architecture and OpenAPI v0.1 contracts.
2. Start both workers from that same commit.
3. Integrate runtime, migrations, health checks, and deterministic seeds.
4. Integrate provider search and reservation operations.
5. Complete planning and itinerary v1 before adding Workers AI.
6. Add Workers AI at the language boundary.
7. Complete disruption recovery and itinerary v2.
8. Harden reset, setup, documentation, and the demo script.

Automated test suites are deferred by project decision. OpenAPI validation, builds, migrations on empty databases, Compose health checks, and repeatable manual smoke flows remain completion gates.
