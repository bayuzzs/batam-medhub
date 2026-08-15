# Mobile app — structure & conventions

Flutter app for **Batam MedHub**. Entry point: `lib/main.dart` (GoRouter + theme).
Display name: "Batam MedHub"; bundle/application ID: `id.medhub.batam.mobile`.

## Directory structure

- `lib/main.dart` — app entry point, `MaterialApp.router` (router lives in core)
- `lib/ui/` — feature-based UI. Group screens/widgets by feature, one folder per feature.
  - `lib/ui/<feature>/` — e.g. `lib/ui/auth/`, `lib/ui/home/`
  - `lib/ui/core/` — shared app-level things: `app_theme.dart`, `app_colors.dart`,
    `app_spacing.dart`, `app_assets.dart`, `app_container.dart`,
    `primary_radial_gradient.dart`, `app_text_field.dart`, `app_validators.dart`,
    `app_router.dart`, `app_bottom_nav.dart`, `main_shell.dart`, and other shared
    widgets
- `lib/data/` — data layer (no UI imports here)
  - `lib/data/repository/` — data repositories (fetch/sync domain data, abstracts + impls)
  - `lib/data/service/` — services (API/network clients, platform integrations)
- `lib/models/` — plain data models / DTOs

## Naming conventions

- Feature files live in `lib/ui/<feature>/` and are named `<feature>.dart`
  (e.g. `lib/ui/auth/login_page.dart`, `lib/ui/home/home_page.dart`).
- Keep one screen/widget per file, snake_case file names, PascalCase class names
  matching the file (e.g. `login_page.dart` → `LoginPage`).

## Data flow

`lib/ui/` → `lib/data/repository/` → `lib/data/service/` → `lib/models/`
Repositories expose models from `lib/models/`; services handle transport/plumbing.

## Routing

- All routes live centrally in **`lib/ui/core/app_router.dart`**:
  - `AppRoutes` — route path constants (e.g. `AppRoutes.login`).
  - `AppRouter.router` — the `GoRouter` config wired into `MaterialApp.router` in `main.dart`.
  - `AppRouterX` — typed navigation helpers on `BuildContext`, e.g.
    `context.pushLogin()` or `context.goLogin()`.
- The three bottom-nav destinations (History / New Itinerary / Profile) live in a
  `StatefulShellRoute.indexedStack` with one `StatefulShellBranch` per tab. The
  shell Scaffold + `AppBottomNav` live in **`lib/ui/core/main_shell.dart`**
  (`MainShell`), which takes a `StatefulNavigationShell` and calls
  `goBranch` on tab select. Each branch keeps its own state (chat scroll,
  form input, etc.) across tab switches.
- Feature pages that are shell destinations return **body content only** (no
  `Scaffold`/`AppBottomNav`/`_selectTab`) — the shell owns those. When pumping
  such a page in a widget test, wrap it in a `Scaffold` (e.g.
  `MaterialApp(home: Scaffold(body: ChatPage()))`).
- When adding a route: add its path to `AppRoutes`, a `GoRoute` to `AppRouter.router`,
  and typed helpers to `AppRouterX`.
- Helper semantics (two behaviors only):
  - `pushXxx()` — **push** onto the navigation stack (back returns to previous screen).
  - `goXxx()` — **replace** the current screen with the target.
- **Always navigate via the typed helpers** (`context.pushLogin()`, …) — never call
  `context.push('/login')` with raw strings or import `go_router` directly in feature pages.

## Reuse before you create

- **Reuse existing components** whenever possible. Before creating a new widget,
  screen, theme value, color, asset, model, repository, or service, check whether
  something already exists in `lib/ui/core/`, `lib/models/`, `lib/data/`, or the
  relevant `lib/ui/<feature>/` folder.
- If a component already exists (e.g. an app bar, button, or theme constant),
  **use it instead of duplicating** it or adding a near-identical variant.
- Prefer extending/updating the shared component in `lib/ui/core/` over copying
  it into a feature folder.

## Validation

- Centralize form validation in **`lib/ui/core/app_validators.dart`**
  (`AppValidators`). Use `AppValidators.email`, `AppValidators.password`, etc.
  as the `validator` on `AppTextField`/`TextFormField` — don't inline
  validation rules or messages in pages.

## Verification

- Prefer **widget tests** (in `test/`) to verify UI behavior instead of manually
  running the app. Run them with `flutter test`.
- When checking that a screen renders or that navigation/routing works, write or
  extend a widget test (e.g. pump `MyApp`, assert on visible text, tap a button,
  then assert the resulting screen) rather than launching the app to eyeball it.
- Use `flutter analyze` to confirm there are no lint/compile issues before
  finishing UI changes.
