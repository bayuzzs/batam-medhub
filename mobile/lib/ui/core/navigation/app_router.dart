import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:mobile/application/auth/auth_controller.dart';
import 'package:mobile/application/auth/providers.dart';
import 'package:mobile/ui/auth/login_page.dart';
import 'package:mobile/ui/auth/onboarding_page.dart';
import 'package:mobile/ui/auth/register_page.dart';
import 'package:mobile/ui/chat/chat_page.dart';
import 'package:mobile/ui/core/navigation/main_shell.dart';
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
/// A Riverpod provider so the router can react to auth state:
/// - While a session is being restored on startup, redirects are paused
///   (`AuthStatus.restoring`) so the app doesn't flicker to login.
/// - Authenticated users are pulled out of the auth screens into the shell.
/// - Unauthenticated users are redirected to login away from the shell.
///
/// Add new [GoRoute]s here as the app grows (e.g. appointment detail, etc.)
/// and expose typed helpers in [AppRouterX].
final appRouterProvider = Provider<GoRouter>((ref) {
  final router = GoRouter(
    initialLocation: AppRoutes.onboarding,
    redirect: (context, state) {
      final auth = ref.read(authControllerProvider);
      final location = state.matchedLocation;
      final isAuthScreen =
          location == AppRoutes.onboarding ||
          location == AppRoutes.login ||
          location == AppRoutes.register;

      // Wait for the startup restore before redirecting anywhere.
      if (auth.status == AuthStatus.restoring) {
        return null;
      }
      if (auth.isAuthenticated && isAuthScreen) {
        return AppRoutes.chat;
      }
      if (auth.status == AuthStatus.unauthenticated && !isAuthScreen) {
        return AppRoutes.login;
      }
      return null;
    },
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

  // Re-evaluate redirects whenever auth state changes.
  ref.listen(authControllerProvider, (_, _) => router.refresh());

  return router;
});

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
