// Basic widget smoke tests for the app entry point and auth flow routing.

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:mobile/main.dart';

void main() {
  testWidgets('App boots into onboarding screen', (WidgetTester tester) async {
    // Build our app and trigger a frame.
    await tester.pumpWidget(const MyApp());

    // Onboarding content is shown.
    expect(find.text('Your Health Journey, Orchestrated'), findsOneWidget);
    expect(find.widgetWithText(FilledButton, 'Get Started'), findsOneWidget);
  });

  testWidgets('Get Started routes to the login screen', (
    WidgetTester tester,
  ) async {
    await tester.pumpWidget(const MyApp());

    // Tap "Get Started" and let the route transition settle.
    await tester.tap(find.text('Get Started'));
    await tester.pumpAndSettle();

    // Login screen is shown.
    expect(find.text('Welcome Back!'), findsOneWidget);
    expect(find.text('Login to continue Batam MedHub'), findsOneWidget);
  });
}
