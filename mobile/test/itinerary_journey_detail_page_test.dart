// Widget smoke tests for the Itinerary Journey Detail screen.

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:mobile/ui/itinerary/itinerary_journey_detail_page.dart';

void main() {
  Widget wrap() {
    return const MaterialApp(
      home: ItineraryJourneyDetailPage(
        imageUrl: 'assets/images/barelang.png',
        providerName: 'RS Awal Bros Batam',
        serviceName: 'Cardiac Screening Package',
        location: 'Batu Aji · 5 km',
        appointment: 'Tomorrow, 09:00',
        duration: '3 days',
        price: 'IDR 4.5jt',
      ),
    );
  }

  testWidgets('renders hospital card, timeline, summary, and CTA', (
    WidgetTester tester,
  ) async {
    await tester.pumpWidget(wrap());

    // Hospital card.
    expect(find.text('Itinerary Journey'), findsOneWidget);
    // The provider name shows in the hospital card and again as the timeline
    // appointment location, so expect at least one.
    expect(find.text('RS Awal Bros Batam'), findsWidgets);
    expect(find.text('Cardiac Screening Package'), findsWidgets);
    expect(find.textContaining('Batu Aji'), findsWidgets);
    expect(
      find.textContaining('Appointment · Tomorrow, 09:00'),
      findsOneWidget,
    );

    // Journey timeline: date header + a step with time/activity/location.
    expect(find.text('Day 1 · Tomorrow'), findsOneWidget);
    expect(find.text('06:30'), findsOneWidget);
    expect(find.text('Ferry departure'), findsOneWidget);
    expect(find.textContaining('HarbourFront Ferry Terminal'), findsOneWidget);

    // Itinerary summary.
    expect(find.text('Itinerary Summary'), findsOneWidget);
    expect(find.text('Total Duration'), findsOneWidget);
    expect(find.text('3 days'), findsOneWidget);
    expect(find.text('Cost Estimate'), findsOneWidget);
    expect(find.text('IDR 4.5jt/person'), findsOneWidget);

    // Included items + CTA.
    expect(find.text('Ferry'), findsOneWidget);
    expect(find.text('Hotel'), findsOneWidget);
    expect(find.text('Choose This Itinerary'), findsOneWidget);
  });

  testWidgets('tapping choose shows a confirmation snackbar', (
    WidgetTester tester,
  ) async {
    await tester.pumpWidget(wrap());

    await tester.tap(find.text('Choose This Itinerary'));
    await tester.pump();

    expect(
      find.text('Itinerary selected! Booking is coming soon.'),
      findsOneWidget,
    );
  });
}
