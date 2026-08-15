// Basic widget smoke tests for the app entry point and auth flow routing.

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:mobile/main.dart';

/// Pumps the app inside a fresh [ProviderScope].
///
/// The router is now a Riverpod provider ([appRouterProvider]) and each test
/// builds its own scope, so auth state starts fresh (unauthenticated at
/// onboarding) and no cross-test reset is needed.
Future<void> _pumpApp(WidgetTester tester) async {
  await tester.pumpWidget(const ProviderScope(child: MyApp()));
  await tester.pumpAndSettle();
}

/// Logs in through the onboarding → login flow, landing on the chat screen
/// (the first shell destination shown after login).
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

  testWidgets('Register routes to the chat screen', (WidgetTester tester) async {
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

  testWidgets('Logout returns to the login screen', (WidgetTester tester) async {
    await _pumpApp(tester);
    await _loginToChat(tester);

    // Open the Profile tab — shows the signed-in patient.
    await tester.tap(find.byKey(const Key('app_bottom_nav_item_2')));
    await tester.pumpAndSettle();
    expect(find.text('Profile'), findsOneWidget);
    expect(find.text('Rina Tan'), findsOneWidget);

    // Sign out → router redirects back to login.
    await tester.tap(find.widgetWithText(FilledButton, 'Log Out'));
    await tester.pumpAndSettle();
    expect(find.text('Welcome Back!'), findsOneWidget);
  });

  testWidgets('Bottom nav switches from chat to History', (
    WidgetTester tester,
  ) async {
    await _pumpApp(tester);

    // Login through to the chat screen.
    await _loginToChat(tester);

    // Bottom nav shows the three destinations (label-less, keyed items).
    expect(find.byKey(const Key('app_bottom_nav_item_0')), findsOneWidget);
    expect(find.byKey(const Key('app_bottom_nav_item_1')), findsOneWidget);
    expect(find.byKey(const Key('app_bottom_nav_item_2')), findsOneWidget);

    // Tap History → history page is shown.
    await tester.tap(find.byKey(const Key('app_bottom_nav_item_0')));
    await tester.pumpAndSettle();
    expect(
      find.text('Your past medical journeys will appear here.'),
      findsOneWidget,
    );
  });

  testWidgets('Bottom nav is a floating pill, not full-screen height', (
    WidgetTester tester,
  ) async {
    await _pumpApp(tester);

    // Login through to the chat screen (the shell with the nav bar).
    await _loginToChat(tester);

    // The nav bar must shrink-wrap its content, not stretch to the screen.
    // (Scaffold passes a full-screen-height constraint to the
    // bottomNavigationBar slot; a Center would expand to fill it.)
    final pill = find.byKey(const Key('app_bottom_nav_pill'));
    final pillSize = tester.getSize(pill);
    final screenSize = tester.getSize(find.byType(Scaffold).first);

    // Fit-to-content width: the pill is narrower than the screen.
    expect(pillSize.width, lessThan(screenSize.width));
    // Floating height: nowhere near full screen.
    expect(pillSize.height, lessThan(screenSize.height / 2));
  });

  testWidgets('Selected circle animates to the tapped item', (
    WidgetTester tester,
  ) async {
    await _pumpApp(tester);
    await _loginToChat(tester);

    // The circle starts over the active tab (chat = item 1).
    final indicator = find.byKey(const Key('app_bottom_nav_indicator'));
    final item0 = find.byKey(const Key('app_bottom_nav_item_0'));
    final item1 = find.byKey(const Key('app_bottom_nav_item_1'));
    expect(tester.getCenter(indicator).dx, tester.getCenter(item1).dx);

    // Tap History — mid-animation the circle is between the two items.
    await tester.tap(item0);
    await tester.pump();
    await tester.pump(const Duration(milliseconds: 100));
    final midDx = tester.getCenter(indicator).dx;
    expect(midDx, greaterThan(tester.getCenter(item0).dx));
    expect(midDx, lessThan(tester.getCenter(item1).dx));

    // Once settled it sits over the tapped item.
    await tester.pumpAndSettle();
    expect(tester.getCenter(indicator).dx, tester.getCenter(item0).dx);
  });
}
