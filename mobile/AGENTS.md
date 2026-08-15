# Mobile app — structure & conventions

Flutter app for **Batam MedHub**. Entry point: `lib/main.dart` (GoRouter + theme).
Display name: "Batam MedHub"; bundle/application ID: `id.medhub.batam.mobile`.

## Directory structure

- `lib/main.dart` — app entry point. Wraps `MyApp` (a `ConsumerWidget`) in `ProviderScope`; `MaterialApp.router` uses the `appRouterProvider`.
- `lib/application/` — app-level orchestration / global state (Riverpod). One folder per concern:
  - `lib/application/auth/` — auth state & wiring: `auth_controller.dart`, `providers.dart`
- `lib/ui/` — feature-based UI. Group screens/widgets by feature, one folder per feature.
  - `lib/ui/<feature>/` — e.g. `lib/ui/auth/`, `lib/ui/home/`
  - `lib/ui/core/` — shared app-level things, grouped by concern:
    - `lib/ui/core/theme/` — design tokens: `app_theme.dart`, `app_colors.dart`,
      `app_spacing.dart`, `app_assets.dart`
    - `lib/ui/core/widgets/` — reusable widgets & validation:
      `app_container.dart`, `app_text_field.dart`, `app_validators.dart`,
      `primary_radial_gradient.dart`, `itenary_option_card.dart`, and others
    - `lib/ui/core/navigation/` — routing & shell: `app_router.dart`,
      `main_shell.dart`, `app_bottom_nav.dart`
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

- All routes live centrally in **`lib/ui/core/navigation/app_router.dart`**:
  - `AppRoutes` — route path constants (e.g. `AppRoutes.login`).
  - `appRouterProvider` — the `GoRouter` config as a Riverpod provider, wired into
    `MaterialApp.router` in `main.dart`.
  - `AppRouterX` — typed navigation helpers on `BuildContext`, e.g.
    `context.pushLogin()` or `context.goLogin()`.
- The router is auth-aware: its `redirect` reads `authControllerProvider`. While
  `AuthStatus.restoring` no redirect happens (no flicker to login on boot);
  authenticated users are pulled from the auth screens into the shell; unauthenticated
  users are redirected to `/login` away from the shell. `appRouterProvider` calls
  `router.refresh()` whenever auth state changes.
- The three bottom-nav destinations (History / New Itinerary / Profile) live in a
  `StatefulShellRoute.indexedStack` with one `StatefulShellBranch` per tab. The
  shell Scaffold + `AppBottomNav` live in **`lib/ui/core/navigation/main_shell.dart`**
  (`MainShell`), which takes a `StatefulNavigationShell` and calls
  `goBranch` on tab select. Each branch keeps its own state (chat scroll,
  form input, etc.) across tab switches.
- Feature pages that are shell destinations return **body content only** (no
  `Scaffold`/`AppBottomNav`/`_selectTab`) — the shell owns those. When pumping
  such a page in a widget test, wrap it in a `Scaffold` (e.g.
  `MaterialApp(home: Scaffold(body: ChatPage()))`).
- When adding a route: add its path to `AppRoutes`, a `GoRoute` to `appRouterProvider`,
  and typed helpers to `AppRouterX`.
- Helper semantics (two behaviors only):
  - `pushXxx()` — **push** onto the navigation stack (back returns to previous screen).
  - `goXxx()` — **replace** the current screen with the target.
- **Always navigate via the typed helpers** (`context.pushLogin()`, …) — never call
  `context.push('/login')` with raw strings or import `go_router` directly in feature pages.

## Authentication

Auth lives in `lib/application/auth/` (state/wiring), `lib/data/` (repositories,
services) and `lib/models/` (session/profile DTOs).

- **`AuthController`** (`lib/application/auth/auth_controller.dart`) — a Riverpod
  `Notifier<AuthState>`. Restores a persisted session on startup, handles
  login/register/logout, and schedules an access-token refresh just before the JWT
  `exp` claim (`AuthSession.accessExpiresAt`, see the `AuthSessionExpiry` extension
  in `lib/models/auth_session.dart`). `refresh()` shares one in-flight future so
  concurrent callers (timer + Dio 401) only rotate the single-use refresh token
  once; on failure the user is signed out.
- **`providers.dart`** — DI switch. `kUseFakeBackend = true` (default) uses
  `FakeAuthRepository` + `InMemoryTokenStore` (no backend / platform storage). Set
  it to `false` to use the real Dio backend (`AuthRepositoryImpl` + `DioAuthApi`) +
  `SecureTokenStore`. Override `authRepositoryProvider`/`tokenStoreProvider` in
  tests to swap implementations.
- **`lib/data/service/auth_interceptor.dart`** — Dio interceptor that attaches the
  Bearer access token to non-auth requests and transparently refreshes + retries
  once on a `401`.
- Fake mode accepts any well-formed credentials; documented demo user is
  `rina.tan@example.test` / `Demo-Only-Password-2026!`.
- `LoginPage`/`RegisterPage` call `authControllerProvider.notifier.login/register`;
  on success the router redirects to `/chat`. `ProfilePage` shows the session's
  patient profile and a Log Out action.

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

- Centralize form validation in **`lib/ui/core/widgets/app_validators.dart`**
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
