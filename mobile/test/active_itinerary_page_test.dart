// Widget tests for the active itinerary screen (`ActiveItineraryPage`).

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:mobile/models/journey.dart';
import 'package:mobile/models/money.dart';
import 'package:mobile/models/plan_option.dart';
import 'package:mobile/models/time_window.dart';
import 'package:mobile/ui/itinerary/active_itinerary_page.dart';

/// Builds a confirmed journey matching the fake backend's active itinerary.
JourneyDetail _buildDetail() {
  return JourneyDetail(
    journey: Journey(
      id: 'journey-000001',
      tripRequestId: 'trip-000001',
      status: JourneyStatus.active,
      activeItineraryVersion: 1,
      createdAt: DateTime.utc(2026, 8, 15, 8, 10),
      updatedAt: DateTime.utc(2026, 8, 15, 8, 10),
    ),
    activeItinerary: ItineraryVersion(
      id: 'itinerary-v1',
      journeyId: 'journey-000001',
      version: 1,
      status: ItineraryVersionStatus.active,
      totalPrice: const PriceSummary(
        sourceTotals: [Money(amountMinor: 25190, currency: 'SGD')],
        displayTotal: Money(amountMinor: 25190, currency: 'SGD'),
        estimated: true,
      ),
      items: [
        ItineraryItem(
          id: 'item-ferry-out',
          itemType: ItemType.ferryOutbound,
          providerId: 'ferry-demo-01',
          title: 'HarbourFront to Batam Centre',
          status: ItineraryItemStatus.confirmed,
          timeWindow: TimeWindow(
            startsAt: DateTime.utc(2026, 8, 21, 23, 30),
            endsAt: DateTime.utc(2026, 8, 22, 0, 40),
            startTimeZone: 'Asia/Singapore',
            endTimeZone: 'Asia/Jakarta',
          ),
          operationalNotes: const [],
          synthetic: true,
          source: 'MOCK',
        ),
        ItineraryItem(
          id: 'item-hospital',
          itemType: ItemType.hospitalAppointment,
          providerId: 'hospital-demo-01',
          title: 'Basic Medical Check-up',
          status: ItineraryItemStatus.confirmed,
          timeWindow: TimeWindow(
            startsAt: DateTime.utc(2026, 8, 22, 3),
            endsAt: DateTime.utc(2026, 8, 22, 5),
            startTimeZone: 'Asia/Jakarta',
            endTimeZone: 'Asia/Jakarta',
          ),
          operationalNotes: const [
            'Follow only the preparation instructions supplied by the hospital.',
          ],
          synthetic: true,
          source: 'MOCK',
        ),
        ItineraryItem(
          id: 'item-ferry-return',
          itemType: ItemType.ferryReturn,
          providerId: 'ferry-demo-01',
          title: 'Batam Centre to HarbourFront',
          status: ItineraryItemStatus.confirmed,
          timeWindow: TimeWindow(
            startsAt: DateTime.utc(2026, 8, 22, 7, 30),
            endsAt: DateTime.utc(2026, 8, 22, 8, 40),
            startTimeZone: 'Asia/Jakarta',
            endTimeZone: 'Asia/Singapore',
          ),
          operationalNotes: const [],
          synthetic: true,
          source: 'MOCK',
        ),
      ],
      createdAt: DateTime.utc(2026, 8, 15, 8, 10),
    ),
    itineraryVersions: [
      ItineraryVersionSummary(
        id: 'itinerary-v1',
        version: 1,
        status: ItineraryVersionStatus.active,
        createdAt: DateTime.utc(2026, 8, 15, 8, 10),
      ),
    ],
  );
}

void main() {
  testWidgets('renders the active itinerary journey', (
    WidgetTester tester,
  ) async {
    tester.view.physicalSize = const Size(800, 1600);
    tester.view.devicePixelRatio = 1.0;
    addTearDown(tester.view.reset);

    await tester.pumpWidget(
      MaterialApp(home: ActiveItineraryPage(detail: _buildDetail())),
    );

    // Header + status banner.
    expect(find.text('My Itinerary'), findsOneWidget);
    expect(find.text('Journey active'), findsOneWidget);
    expect(find.textContaining('Booking ref · journey-000001'), findsOneWidget);

    // Summary card: service, provider, total price.
    expect(find.text('Basic Medical Check-up'), findsWidgets);
    expect(find.text('Hospital Demo 01'), findsWidgets);
    expect(find.text('SGD 251.90'), findsOneWidget);

    // Confirmed journey timeline.
    expect(find.text('Confirmed journey, step by step'), findsOneWidget);
    expect(find.text('HarbourFront to Batam Centre'), findsOneWidget);
    expect(find.text('Batam Centre to HarbourFront'), findsOneWidget);
    expect(find.textContaining('Ferry outbound'), findsOneWidget);
    expect(find.textContaining('Hospital appointment'), findsOneWidget);
    expect(find.textContaining('Ferry return'), findsOneWidget);

    // Confirmed status chips (one per booked leg).
    expect(find.text('Confirmed'), findsNWidgets(3));

    // Operational note on the hospital leg.
    expect(
      find.textContaining('preparation instructions supplied by the hospital'),
      findsOneWidget,
    );

    // Footer note.
    expect(find.textContaining("your device's local time"), findsOneWidget);
  });
}
