// Widget tests for the itemized plan detail screen (`PlanDetailPage`).

import 'dart:convert';
import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'package:mobile/application/chat/providers.dart';
import 'package:mobile/application/journey/providers.dart';
import 'package:mobile/data/repository/fake_journey_repository.dart';
import 'package:mobile/models/journey.dart';
import 'package:mobile/models/plan_option.dart';
import 'package:mobile/ui/chat/chat_page.dart';
import 'package:mobile/ui/core/navigation/app_router.dart';
import 'package:mobile/ui/itinerary/active_itinerary_page.dart';
import 'package:mobile/ui/itinerary/plan_detail_page.dart';

/// Loads the single plan option from the `plan-result.json` fixture.
PlanOption _loadOption() {
  final json =
      jsonDecode(File('test/fixtures/core/plan-result.json').readAsStringSync())
          as Map<String, dynamic>;
  final result = PlanningResult.fromJson(json.cast<String, dynamic>());
  return result.options.first;
}

/// Pumps the chat screen under a real [GoRouter] with the `/chat` and `/plan`
/// routes, backed by an instant fake journey repository.
Future<void> _pumpChatRouter(WidgetTester tester) async {
  tester.view.physicalSize = const Size(800, 1600);
  tester.view.devicePixelRatio = 1.0;
  addTearDown(tester.view.reset);
  SharedPreferences.setMockInitialValues({});

  final router = GoRouter(
    initialLocation: AppRoutes.chat,
    routes: [
      GoRoute(
        path: AppRoutes.chat,
        builder: (context, state) => const ChatPage(),
      ),
      GoRoute(
        path: AppRoutes.planDetail,
        builder: (context, state) =>
            PlanDetailPage(option: state.extra! as PlanOption),
      ),
      GoRoute(
        path: AppRoutes.activeItinerary,
        builder: (context, state) =>
            ActiveItineraryPage(detail: state.extra! as JourneyDetail),
      ),
    ],
  );

  await tester.pumpWidget(
    ProviderScope(
      overrides: [
        journeyRepositoryProvider.overrideWithValue(
          FakeJourneyRepository(delay: Duration.zero),
        ),
      ],
      child: MaterialApp.router(routerConfig: router),
    ),
  );
  await tester.pumpAndSettle();
}

/// Walks the chat conversation until the two plan cards appear.
Future<void> _driveToPlanCards(WidgetTester tester) async {
  await tester.enterText(find.byType(TextField), 'I need a check-up in Batam');
  await tester.tap(find.byKey(const Key('chat_send_button')));
  await tester.pumpAndSettle();
  await tester.enterText(find.byType(TextField), 'Basic, next Friday please');
  await tester.tap(find.byKey(const Key('chat_send_button')));
  await tester.pumpAndSettle();
}

void main() {
  testWidgets('renders the itemized plan details', (WidgetTester tester) async {
    await tester.pumpWidget(
      ProviderScope(
        child: MaterialApp(home: PlanDetailPage(option: _loadOption())),
      ),
    );

    // Header.
    expect(find.text('Journey Plan'), findsOneWidget);

    // Summary card: service, provider, recommended banner, total price.
    expect(find.text('Basic Medical Check-up'), findsWidgets);
    expect(find.text('Hospital Demo 01'), findsWidgets);
    expect(find.text('Recommended'), findsOneWidget);
    expect(find.text('SGD 251.90'), findsOneWidget);

    // Itemized timeline.
    expect(find.text('Your journey, step by step'), findsOneWidget);
    expect(find.text('HarbourFront to Batam Centre'), findsOneWidget);
    expect(find.text('Batam Centre to HarbourFront'), findsOneWidget);
    expect(find.textContaining('Ferry outbound'), findsOneWidget);
    expect(find.textContaining('Hospital appointment'), findsOneWidget);
    expect(find.textContaining('Arrival buffer'), findsOneWidget);

    // Per-leg price (ferry SGD 50.00 and hospital SGD 126.58).
    expect(find.text('SGD 50.00'), findsNWidgets(2));
    expect(find.text('SGD 126.58'), findsOneWidget);

    // Planner explanation + book action.
    expect(find.text('Why this plan'), findsOneWidget);
    expect(
      find.textContaining('ferry arrives with 140 minutes'),
      findsOneWidget,
    );
    expect(find.text('Book this journey'), findsOneWidget);
  });

  testWidgets('View Details on a chat plan card opens the plan detail', (
    WidgetTester tester,
  ) async {
    await _pumpChatRouter(tester);
    await _driveToPlanCards(tester);

    expect(find.text('View Details'), findsNWidgets(2));

    await tester.tap(find.text('View Details').first);
    await tester.pumpAndSettle();

    expect(find.byType(PlanDetailPage), findsOneWidget);
    expect(find.text('Journey Plan'), findsOneWidget);
    expect(find.text('Your journey, step by step'), findsOneWidget);
    expect(find.textContaining('Basic Medical Check-up'), findsWidgets);
  });

  testWidgets('Booking from the plan detail opens the active itinerary', (
    WidgetTester tester,
  ) async {
    await _pumpChatRouter(tester);
    await _driveToPlanCards(tester);

    // Open the first plan's detail screen and book it.
    await tester.tap(find.text('View Details').first);
    await tester.pumpAndSettle();
    expect(find.byType(PlanDetailPage), findsOneWidget);

    await tester.tap(find.text('Book this journey'));
    await tester.pumpAndSettle();

    // Lands on the active itinerary — not the chat or plan detail.
    expect(find.byType(ActiveItineraryPage), findsOneWidget);
    expect(find.byType(PlanDetailPage), findsNothing);
    expect(find.text('My Itinerary'), findsOneWidget);
    expect(find.text('Journey active'), findsOneWidget);
    expect(find.textContaining('Booking ref · journey-000001'), findsOneWidget);

    // The journey was confirmed in the chat controller.
    final container = ProviderScope.containerOf(
      tester.element(find.byType(ActiveItineraryPage)),
    );
    expect(container.read(chatControllerProvider).journey, isNotNull);
  });
}
