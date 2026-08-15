# Batam MedHub — monorepo

Monorepo for the **Batam MedHub** project.

## Layout

- `mobile/` — the Flutter mobile app (see `mobile/AGENTS.md` for its structure & conventions)
- `Taskfile.yml` — repo task runner config (go-task)
- `readme.md` — placeholder

## Tasks

Taskfiles are namespaced by directory; run from the repo root:

```
task mobile:run:linux     # flutter run -d linux
task mobile:run:android   # flutter run -d android
task mobile:run:chrome    # flutter run -d chrome
task --list-all           # list all tasks
```

See `mobile/Taskfile.yml` for the mobile-specific task definitions.

## Mobile app

The Flutter app lives in `mobile/`. For its directory conventions, naming rules,
and data flow (`lib/ui/` → `lib/data/repository/` → `lib/data/service/` →
`lib/models/`), read **`mobile/AGENTS.md`** before editing code in `mobile/`.
