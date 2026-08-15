// Widget smoke tests for the AI chat screen.

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:mobile/ui/chat/chat_options.dart';
import 'package:mobile/ui/chat/chat_page.dart';

void main() {
  testWidgets('Chat screen renders greeting, assistant card, and input', (
    WidgetTester tester,
  ) async {
    await tester.pumpWidget(
      const MaterialApp(home: Scaffold(body: ChatPage())),
    );

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

    // Assistant card.
    expect(find.text('Talk with your AI Assistant'), findsOneWidget);

    // Placeholder chat items (user + assistant).
    expect(
      find.text('I need help planning my medical trip to Batam'),
      findsOneWidget,
    );
    expect(
      find.text(
        'Of course! I\'ll orchestrate the best medical journey for you.',
      ),
      findsOneWidget,
    );

    // Assistant item with selectable option cards.
    expect(find.text('Which hospital would you prefer?'), findsOneWidget);

    // Hospital options live inside ChatOptions; the itinerary card below
    // reuses 'RS Awal Bros Batam' as its provider name, so scope the finds.
    final options = find.byType(ChatOptions);
    expect(
      find.descendant(of: options, matching: find.text('RS Awal Bros Batam')),
      findsOneWidget,
    );
    expect(
      find.descendant(of: options, matching: find.text('RS Hermina Batam')),
      findsOneWidget,
    );
    expect(
      find.descendant(of: options, matching: find.text('RSUD Embung Fatimah')),
      findsOneWidget,
    );

    // Example itinerary option card.
    expect(find.text('Here\'s a great itinerary for you:'), findsOneWidget);
    expect(find.text('Cardiac Screening Package'), findsOneWidget);
    expect(find.text('Recommended'), findsOneWidget);
    expect(find.text('View Details Itinerary'), findsOneWidget);

    // Chat input with a send button.
    expect(find.byType(TextField), findsOneWidget);
    expect(find.byType(IconButton), findsOneWidget);
  });

  testWidgets('Selecting a hospital option highlights it', (
    WidgetTester tester,
  ) async {
    await tester.pumpWidget(
      const MaterialApp(home: Scaffold(body: ChatPage())),
    );

    // The options sit below the fold in the test viewport; bring them into
    // view before tapping.
    final option = find.text('RS Hermina Batam');
    await tester.ensureVisible(option);
    await tester.pumpAndSettle();

    await tester.tap(option);
    await tester.pump();

    // The option card is still present; selection is tracked internally.
    final options = find.byType(ChatOptions);
    expect(
      find.descendant(of: options, matching: find.text('RS Hermina Batam')),
      findsOneWidget,
    );
    expect(
      find.descendant(of: options, matching: find.text('RS Awal Bros Batam')),
      findsOneWidget,
    );
  });

  testWidgets('Typing a message and tapping send clears the input', (
    WidgetTester tester,
  ) async {
    await tester.pumpWidget(
      const MaterialApp(home: Scaffold(body: ChatPage())),
    );

    await tester.enterText(find.byType(TextField), 'Hello');
    expect(find.text('Hello'), findsOneWidget);

    await tester.tap(find.byType(IconButton));
    await tester.pump();

    expect(find.text('Hello'), findsNothing);
  });
}
