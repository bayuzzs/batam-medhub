# Mobile app — structure & conventions

Flutter app for **Batam MedHub**. Entry point: `lib/main.dart` (GoRouter + theme).
Display name: "Batam MedHub"; bundle/application ID: `id.medhub.batam.mobile`.

## Directory structure

- `lib/main.dart` — app entry point. Builds a `ProviderContainer`, awaits
  `authControllerProvider.notifier.restore()` before the first frame (so the
  app boots straight into the right screen without an auth flicker), then runs
  `MyApp` under an `UncontrolledProviderScope`. `MaterialApp.router` uses the
  `appRouterProvider`.
- `lib/application/` — app-level orchestration / global state (Riverpod). One folder per concern:
  - `lib/application/auth/` — auth state & wiring: `auth_controller.dart`, `providers.dart`
  - `lib/application/journey/` — journey DI wiring: `providers.dart`
    (`journeyRepositoryProvider` → `FakeJourneyRepository`/`JourneyRepositoryImpl`
    behind `kUseFakeBackend`; overridable in tests)
  - `lib/application/chat/` — chat conversation state: `chat_controller.dart`
    (`ChatController` + `ChatState`), `providers.dart` (`chatControllerProvider`)
- `lib/ui/` — feature-based UI. Group screens/widgets by feature, one folder per feature.
  - `lib/ui/<feature>/` — e.g. `lib/ui/auth/`, `lib/ui/chat/`, `lib/ui/history/`,
    `lib/ui/profile/`, `lib/ui/itinerary/`
  - `lib/ui/core/` — shared app-level things, grouped by concern:
    - `lib/ui/core/theme/` — design tokens: `app_theme.dart`, `app_colors.dart`,
      `app_spacing.dart`, `app_assets.dart`
    - `lib/ui/core/widgets/` — reusable widgets & validation:
      `app_container.dart`, `app_text_field.dart`, `app_validators.dart`,
      `primary_radial_gradient.dart`, `itenary_option_card.dart`, and others
    - `lib/ui/core/utils/` — formatting helpers: `money_formatter.dart`,
      `time_window_format.dart` (`formatWindow`/`formatWindowWithZone`),
      `provider_label.dart` (`formatProviderLabel`)
    - `lib/ui/core/navigation/` — routing: `app_router.dart`
- `lib/data/` — data layer (no UI imports here)
  - `lib/data/repository/` — data repositories (fetch/sync domain data, abstracts + impls):
    `journey_repository.dart` (abstract + `JourneyException`), `journey_repository_impl.dart`
    (Dio), `fake_journey_repository.dart` (deterministic in-memory demo, default)
  - `lib/data/service/` — services (API/network clients, platform integrations):
    `journey_api.dart` + `dio_journey_api.dart` (trip request → plans → confirm → itinerary)
- `lib/models/` — plain data models / DTOs: journey domain models `money.dart`,
  `time_window.dart`, `structured_intent.dart`, `trip_request.dart`, `plan_option.dart`,
  `journey.dart`, `medical_service.dart` (freezed + json_serializable)

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
  authenticated users are pulled from the auth screens into the **chat screen**;
  unauthenticated users are redirected to `/login`. `appRouterProvider` calls
  `router.refresh()` whenever auth state changes.
- There is **no bottom-nav shell**. The chat screen (`/chat`) is the app's
  primary authenticated screen; `History` (`/history`), `Profile` (`/profile`),
  and `Itinerary Journey` (`/itinerary`) are full-screen routes pushed on top of
  it. Chat owns its own `Scaffold` with a **pinned top bar**: History
  (top-left) and Profile (top-right), which call `context.pushHistory()` /
  `context.pushProfile()`. History/Profile/Itinerary each own their `Scaffold`
  and render a header with a back button (`context.pop()`) that returns to chat.
- Screens pushed on top of chat own their **full-screen `Scaffold`** (no bottom
  nav). In widget tests, pump them directly when they own a `Scaffold` (e.g.
  `MaterialApp(home: ChatPage())`).
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
  `Notifier<AuthState>`. Restores a persisted session on startup (idempotently;
  `main()` awaits it before the first frame), handles
  login/register/logout, and schedules an access-token refresh just before the JWT
  `exp` claim (`AuthSession.accessExpiresAt`, see the `AuthSessionExpiry` extension
  in `lib/models/auth_session.dart`). `refresh()` shares one in-flight future so
  concurrent callers (timer + Dio 401) only rotate the single-use refresh token
  once; on failure the user is signed out.
- **`providers.dart`** — DI switch. `kUseFakeBackend = true` (default) uses
  `FakeAuthRepository` + a `shared_preferences`-backed `SharedPreferencesTokenStore`,
  so the demo session survives a full app restart without a backend or platform
  secure storage. Set it to `false` to use the real Dio backend
  (`AuthRepositoryImpl` + `DioAuthApi`) + `SecureTokenStore`
  (`flutter_secure_storage`). Override `authRepositoryProvider`/`tokenStoreProvider`
  in tests to swap implementations (e.g. `InMemoryTokenStore`).
- **`lib/data/service/auth_interceptor.dart`** — Dio interceptor that attaches the
  Bearer access token to non-auth requests and transparently refreshes + retries
  once on a `401`.
- Fake mode accepts any well-formed credentials; documented demo user is
  `rina.tan@example.test` / `Demo-Only-Password-2026!`.
- `LoginPage`/`RegisterPage` call `authControllerProvider.notifier.login/register`;
  on success the router redirects to `/chat`. `ProfilePage` shows the session's
  patient profile and a Log Out action.

## Journey orchestration & chat

Journey state lives in `lib/application/chat/` + `lib/application/journey/`; the
UI lives in `lib/ui/chat/` and models in `lib/models/`.

- **`ChatController`** (`lib/application/chat/chat_controller.dart`) — a Riverpod
  `Notifier<ChatState>`. `ChatState` holds the message list (`List<ChatMessage>`),
  the current `TripRequestDetail`/`TripRequest`, and `isSending`. `send(text)`:
  appends a user bubble, creates a trip request, and routes on the
  `IntentResolution` returned by the fake backend — `needsClarification` asks a
  follow-up question (subsequent `send()`s call `answerClarification` until the
  intent `matched`), `matched` generates plans, and unsupported/out-of-scope
  services return the reason. Selecting a plan option calls
  `confirmPlanOption`, which shows a "Booking…" status and then a confirmation
  message.
- **`FakeJourneyRepository`** (`lib/data/repository/fake_journey_repository.dart`)
  is the deterministic demo backend used by default (see `kUseFakeBackend` in
  `lib/application/journey/providers.dart`): it returns a needs-clarification
  trip request, a matched intent, two plan options (recommended rank 1, total
  SGD 251.90), and a confirmed active journey. It throws
  `JourneyException(message, code:)` for unknown ids. Override
  `journeyRepositoryProvider` in tests with `FakeJourneyRepository(delay:
  Duration.zero)`.
- **Models**: money uses integer minor units + ISO currency
  (`Money`/`ConvertedMoney`/`PriceSummary`, `fx_rate` is a String);
  `DateWindow`/`TimeWindow` use UTC instants with IANA zone strings
  (`startTimeZone`). Note: freezed/json_serializable use the default
  `explicitToJson: false`, so `toJson()` emits nested objects — round-trips in
  tests go through `jsonDecode(jsonEncode(toJson()))`. Date-only fixture strings
  (e.g. `"2026-08-22"`) parse as **local** `DateTime`, not UTC.
- **UI**: `lib/ui/chat/chat_page.dart` is a `ConsumerStatefulWidget` that renders
  the live conversation from `chatControllerProvider` (`ChatItem` bubbles, a
  `_TypingIndicator`, and `PlanOptionCard`s via `_PlanOptionsList`).
  `lib/ui/chat/plan_option_card.dart` shows the service title, provider,
  time window, formatted total (`MoneyFormatter`, `lib/ui/core/utils/money_formatter.dart`),
  explanation, and a "View Details" action that pushes
  `lib/ui/itinerary/plan_detail_page.dart` (`/plan`, `context.pushPlanDetail(option)`,
  option passed via `state.extra`); the rank-1 option carries a
  "Recommended" banner. Booking happens only on the detail screen (chat cards
  never confirm directly). `PlanDetailPage` renders the hospital
  summary, an **itemized journey timeline** (each leg's time, type, provider,
  per-leg price and operational notes via `MoneyFormatter`/`formatWindow`/`formatProviderLabel`),
  the planner's explanation, and a pinned "Book this journey" action that
  confirms the option and then opens `lib/ui/itinerary/active_itinerary_page.dart`
  (`/active-itinerary`, `context.goActiveItinerary(journey)`) — the patient
  lands on their confirmed journey, **not** back in the chat. `ActiveItineraryPage`
  renders the `JourneyDetail` returned by `confirmPlanOption`: a status banner
  ("Journey active" + booking ref), the hospital summary + total, and the
  active itinerary's booked legs as a timeline with `ItineraryItemStatus` chips.
  Item-type labels/icons are shared via `lib/ui/core/utils/item_type_presentation.dart`.
  `chat_options.dart` holds shared chat UI helpers. Times are shown in
  device-local time; money in the display currency (`.display` on
  `ConvertedMoney`).

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
- Journey-related tests: `test/journey_models_test.dart` (model round-trips +
  fixtures in `test/fixtures/core/`), `test/fake_journey_repository_test.dart`
  (full request→plans→confirm flow + errors),
  `test/chat_controller_test.dart` (conversation state machine + error bubble),
  `test/chat_page_test.dart` (renders greeting, typing + send, clarification,
  plan cards, and booking confirmation),
  `test/plan_detail_page_test.dart` (itemized plan rendering, View Details
  navigation from a chat plan card, and booking from the detail screen).
