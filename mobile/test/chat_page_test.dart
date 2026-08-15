// Widget smoke tests for the AI chat screen.

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'package:mobile/application/journey/providers.dart';
import 'package:mobile/data/repository/fake_journey_repository.dart';
import 'package:mobile/ui/chat/chat_page.dart';
import 'package:mobile/ui/chat/plan_option_card.dart';

/// Pumps the chat screen with an instant fake journey backend so the
/// conversation can be driven end-to-end without real timers. Uses a tall
/// viewport so lazy `ListView.builder` keeps every message built on screen.
Future<void> _pumpChat(WidgetTester tester) async {
  tester.view.physicalSize = const Size(800, 1600);
  tester.view.devicePixelRatio = 1.0;
  addTearDown(tester.view.reset);

  SharedPreferences.setMockInitialValues({});
  await tester.pumpWidget(
    ProviderScope(
      overrides: [
        journeyRepositoryProvider.overrideWithValue(
          FakeJourneyRepository(delay: Duration.zero),
        ),
      ],
      child: const MaterialApp(home: Scaffold(body: ChatPage())),
    ),
  );
  await tester.pumpAndSettle();
}

void main() {
  testWidgets('Chat screen renders greeting, assistant card, and input', (
    WidgetTester tester,
  ) async {
    await _pumpChat(tester);

    // Greeting texts.
    expect(find.text('Hi, Name'), findsOneWidget);
    expect(
      find.text('Let\'s orchestrate your medical journey.'),
      findsOneWidget,
    );
    expect(
      find.text(
        'One request, One connected journey across care, travel and stay.',
      ),
      findsOneWidget,
    );

    // Assistant card and the seeded assistant greeting bubble.
    expect(find.text('Talk with your AI Assistant'), findsOneWidget);
    expect(
      find.textContaining('for example a same-day check-up in Batam'),
      findsOneWidget,
    );

    // Pinned top bar: History (left) and Profile (right).
    expect(find.byKey(const Key('chat_history_button')), findsOneWidget);
    expect(find.byKey(const Key('chat_profile_button')), findsOneWidget);

    // Chat input with a send button.
    expect(find.byType(TextField), findsOneWidget);
    expect(find.byKey(const Key('chat_send_button')), findsOneWidget);
  });

  testWidgets('Typing a message and tapping send clears the input', (
    WidgetTester tester,
  ) async {
    await _pumpChat(tester);

    final field = find.byType(TextField);
    await tester.enterText(field, 'Hello');
    await tester.tap(find.byKey(const Key('chat_send_button')));
    await tester.pumpAndSettle();

    // The input is cleared, and the typed text became a user bubble.
    final textField = tester.widget<TextField>(field);
    expect(textField.controller!.text, isEmpty);
    expect(find.text('Hello'), findsOneWidget);
  });

  testWidgets('Sending a request walks the trip-request conversation', (
    WidgetTester tester,
  ) async {
    await _pumpChat(tester);

    // Ask for help → the assistant asks a clarification question.
    await tester.enterText(
      find.byType(TextField),
      'I need a check-up in Batam',
    );
    await tester.tap(find.byKey(const Key('chat_send_button')));
    await tester.pumpAndSettle();

    expect(find.text('I need a check-up in Batam'), findsOneWidget);
    expect(
      find.textContaining('basic or comprehensive check-up'),
      findsOneWidget,
    );

    // Answer the clarification → plan cards appear.
    await tester.enterText(find.byType(TextField), 'Basic, next Friday please');
    await tester.tap(find.byKey(const Key('chat_send_button')));
    await tester.pumpAndSettle();

    expect(find.byType(PlanOptionCard), findsNWidgets(2));
    // Chat cards lead to the detail screen; booking happens only there.
    expect(find.text('View Details'), findsNWidgets(2));
    expect(find.text('Book this journey'), findsNothing);
    expect(find.text('Recommended'), findsOneWidget);
  });
}
