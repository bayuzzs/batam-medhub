import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'package:mobile/ui/auth/login_page.dart';
import 'package:mobile/ui/auth/onboarding_page.dart';
import 'package:mobile/ui/auth/register_page.dart';

/// Route path constants, kept in one place so screens never hard-code
/// route strings.
abstract final class AppRoutes {
  static const String onboarding = '/';
  static const String login = '/login';
  static const String register = '/register';
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
}
