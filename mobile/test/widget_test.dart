// Basic widget smoke tests for the app entry point and auth flow routing.

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'package:mobile/main.dart';

/// Pumps the app inside a fresh [ProviderScope].
///
/// The router is now a Riverpod provider ([appRouterProvider]) and each test
/// builds its own scope, so auth state starts fresh (unauthenticated at
/// onboarding) and no cross-test reset is needed. The fake token store is
/// `shared_preferences`-backed, so it's mocked with empty initial values.
Future<void> _pumpApp(WidgetTester tester) async {
  SharedPreferences.setMockInitialValues({});
  await tester.pumpWidget(const ProviderScope(child: MyApp()));
  await tester.pumpAndSettle();
}

/// Logs in through the onboarding → login flow, landing on the chat screen
/// (the app's primary authenticated screen after login).
Future<void> _loginToChat(WidgetTester tester) async {
  await tester.tap(find.text('Get Started'));
  await tester.pumpAndSettle();
  await tester.enterText(find.byType(TextFormField).at(0), 'user@medhub.id');
  await tester.enterText(find.byType(TextFormField).at(1), 'secret123');
  await tester.tap(find.widgetWithText(FilledButton, 'Login'));
  await tester.pumpAndSettle();
}

void main() {
  testWidgets('App boots into onboarding screen', (WidgetTester tester) async {
    // Build our app and trigger a frame.
    await _pumpApp(tester);

    // Onboarding content is shown.
    expect(find.text('Your Health Journey, Orchestrated'), findsOneWidget);
    expect(find.widgetWithText(FilledButton, 'Get Started'), findsOneWidget);
  });

  testWidgets('Get Started routes to the login screen', (
    WidgetTester tester,
  ) async {
    await _pumpApp(tester);

    // Tap "Get Started" and let the route transition settle.
    await tester.tap(find.text('Get Started'));
    await tester.pumpAndSettle();

    // Login screen is shown.
    expect(find.text('Welcome Back!'), findsOneWidget);
    expect(find.text('Login to continue Batam MedHub'), findsOneWidget);
  });

  testWidgets('Login routes to the chat screen', (WidgetTester tester) async {
    await _pumpApp(tester);

    await _loginToChat(tester);

    // Chat screen is shown.
    expect(find.text('Hi, Name'), findsOneWidget);
    expect(find.text('Talk with your AI Assistant'), findsOneWidget);
  });

  testWidgets('Register routes to the chat screen', (
    WidgetTester tester,
  ) async {
    await _pumpApp(tester);

    // Onboarding → login → register.
    await tester.tap(find.text('Get Started'));
    await tester.pumpAndSettle();
    await tester.ensureVisible(find.widgetWithText(TextButton, 'Register'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(TextButton, 'Register'));
    await tester.pumpAndSettle();

    expect(find.text('Create Your Account'), findsOneWidget);

    // Fill the register form (full name, email, password, confirm).
    await tester.enterText(find.byType(TextFormField).at(0), 'Rina Tan');
    await tester.enterText(find.byType(TextFormField).at(1), 'user@medhub.id');
    await tester.enterText(find.byType(TextFormField).at(2), 'secret123');
    await tester.enterText(find.byType(TextFormField).at(3), 'secret123');
    await tester.ensureVisible(find.widgetWithText(FilledButton, 'Register'));
    await tester.pumpAndSettle();
    await tester.tap(find.widgetWithText(FilledButton, 'Register'));
    await tester.pumpAndSettle();

    // Authenticated → redirected to the chat shell.
    expect(find.text('Hi, Name'), findsOneWidget);
  });

  testWidgets('Logout returns to the login screen', (
    WidgetTester tester,
  ) async {
    await _pumpApp(tester);
    await _loginToChat(tester);

    // Open Profile from the chat top bar — shows the signed-in patient.
    await tester.tap(find.byKey(const Key('chat_profile_button')));
    await tester.pumpAndSettle();
    expect(find.text('Profile'), findsOneWidget);
    expect(find.text('Rina Tan'), findsOneWidget);

    // Sign out → router redirects back to login.
    await tester.tap(find.widgetWithText(FilledButton, 'Log Out'));
    await tester.pumpAndSettle();
    expect(find.text('Welcome Back!'), findsOneWidget);
  });

  testWidgets('Chat top bar navigates to History and back', (
    WidgetTester tester,
  ) async {
    await _pumpApp(tester);

    // Login through to the chat screen.
    await _loginToChat(tester);

    // The pinned top bar shows both actions.
    expect(find.byKey(const Key('chat_history_button')), findsOneWidget);
    expect(find.byKey(const Key('chat_profile_button')), findsOneWidget);

    // Tap History → history page is shown.
    await tester.tap(find.byKey(const Key('chat_history_button')));
    await tester.pumpAndSettle();
    expect(
      find.text('Your past medical journeys will appear here.'),
      findsOneWidget,
    );

    // Back returns to the chat screen.
    await tester.tap(find.byKey(const Key('page_back_button')));
    await tester.pumpAndSettle();
    expect(find.text('Hi, Name'), findsOneWidget);
  });

  testWidgets('Chat top bar opens Profile and back returns to chat', (
    WidgetTester tester,
  ) async {
    await _pumpApp(tester);

    await _loginToChat(tester);

    await tester.tap(find.byKey(const Key('chat_profile_button')));
    await tester.pumpAndSettle();
    expect(find.text('Rina Tan'), findsOneWidget);

    await tester.tap(find.byKey(const Key('page_back_button')));
    await tester.pumpAndSettle();
    expect(find.text('Hi, Name'), findsOneWidget);
  });
}
