import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'package:mobile/ui/auth/login_page.dart';
import 'package:mobile/ui/auth/onboarding_page.dart';
import 'package:mobile/ui/auth/register_page.dart';
import 'package:mobile/ui/chat/chat_page.dart';
import 'package:mobile/ui/core/main_shell.dart';
import 'package:mobile/ui/history/history_page.dart';
import 'package:mobile/ui/profile/profile_page.dart';

/// Route path constants, kept in one place so screens never hard-code
/// route strings.
abstract final class AppRoutes {
  static const String onboarding = '/';
  static const String login = '/login';
  static const String register = '/register';
  static const String chat = '/chat';
  static const String history = '/history';
  static const String profile = '/profile';
}

/// Central [GoRouter] configuration used by [MaterialApp.router].
///
/// Add new [GoRoute]s here as the app grows (e.g. authenticated home,
/// appointment detail, etc.) and expose typed helpers in [AppRouterX].
abstract final class AppRouter {
  static final GoRouter router = GoRouter(
    initialLocation: AppRoutes.onboarding,
    routes: [
      GoRoute(
        path: AppRoutes.onboarding,
        builder: (context, state) => const OnboardingPage(),
      ),
      GoRoute(
        path: AppRoutes.login,
        builder: (context, state) => const LoginPage(),
      ),
      GoRoute(
        path: AppRoutes.register,
        builder: (context, state) => const RegisterPage(),
      ),
      // Authenticated shell: the three bottom-nav destinations. Each branch
      // keeps its own state via an [IndexedStack].
      StatefulShellRoute.indexedStack(
        builder: (context, state, navigationShell) =>
            MainShell(navigationShell: navigationShell),
        branches: [
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: AppRoutes.history,
                builder: (context, state) => const HistoryPage(),
              ),
            ],
          ),
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: AppRoutes.chat,
                builder: (context, state) => const ChatPage(),
              ),
            ],
          ),
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: AppRoutes.profile,
                builder: (context, state) => const ProfilePage(),
              ),
            ],
          ),
        ],
      ),
    ],
  );
}

/// Typed navigation helpers.
///
/// Two behaviors only, matching how navigation is used in this app:
/// - `pushXxx()` — **push** the screen onto the navigation stack (the
///   previous screen stays below and back returns to it).
/// - `goXxx()` — **replace** the current screen with the target one.
///
/// Prefer these over calling `context.push('/login')` with raw strings, e.g.
/// `context.pushLogin()` or `context.goLogin()`.
extension AppRouterX on BuildContext {
  // Onboarding.
  void pushOnboarding() => push(AppRoutes.onboarding);
  void goOnboarding() => replace(AppRoutes.onboarding);

  // Login.
  void pushLogin() => push(AppRoutes.login);
  void goLogin() => replace(AppRoutes.login);

  // Register.
  void pushRegister() => push(AppRoutes.register);
  void goRegister() => replace(AppRoutes.register);

  // Chat.
  void pushChat() => push(AppRoutes.chat);
  void goChat() => replace(AppRoutes.chat);

  // History.
  void pushHistory() => push(AppRoutes.history);
  void goHistory() => replace(AppRoutes.history);

  // Profile.
  void pushProfile() => push(AppRoutes.profile);
  void goProfile() => replace(AppRoutes.profile);
}
