// Widget smoke tests for the itinerary option card.

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:mobile/ui/core/widgets/itenary_option_card.dart';

void main() {
  Widget wrap(Widget child) {
    return MaterialApp(
      home: Scaffold(body: SingleChildScrollView(child: child)),
    );
  }

  testWidgets('renders banner, provider info, chips, price, and action', (
    WidgetTester tester,
  ) async {
    var detailsTapped = false;

    await tester.pumpWidget(
      wrap(
        ItenaryOptionCard(
          imageUrl: 'assets/images/barelang.png',
          providerName: 'RS Awal Bros Batam',
          serviceName: 'Cardiac Screening Package',
          location: 'Batu Aji · 5 km',
          appointment: 'Tomorrow, 09:00',
          rating: 4.8,
          reviewCount: 214,
          duration: '3 days',
          price: 'IDR 4.500.000',
          onViewDetails: () => detailsTapped = true,
        ),
      ),
    );

    // Recommended banner.
    expect(find.text('Recommended'), findsOneWidget);

    // Provider info.
    expect(find.text('RS Awal Bros Batam'), findsOneWidget);
    expect(find.text('Cardiac Screening Package'), findsOneWidget);
    expect(find.textContaining('Batu Aji'), findsOneWidget);
    expect(find.textContaining('Appointment · Tomorrow'), findsOneWidget);

    // Metadata chips.
    expect(find.text('4.8 (214 reviews)'), findsOneWidget);
    expect(find.text('English Speaking'), findsOneWidget);
    expect(find.text('3 days'), findsOneWidget);

    // Price section.
    expect(find.text('Estimated Total'), findsOneWidget);
    expect(
      find.text('Includes Ferry · Transport · Medical · Hotel'),
      findsOneWidget,
    );
    expect(find.text('IDR 4.500.000/person'), findsOneWidget);

    // Details action fires the callback.
    await tester.tap(find.text('View Itinerary Details'));
    await tester.pump();
    expect(detailsTapped, isTrue);
  });

  testWidgets('hides the Recommended banner when recommended is false', (
    WidgetTester tester,
  ) async {
    await tester.pumpWidget(
      wrap(
        ItenaryOptionCard(
          imageUrl: 'https://invalid.example/hospital.png',
          providerName: 'RS Hermina Batam',
          serviceName: 'General Check Up',
          location: 'Batam Center · 7 km',
          appointment: 'Today, 14:00',
          rating: 4.5,
          reviewCount: 96,
          duration: '2 days',
          price: 'IDR 950.000',
          recommended: false,
        ),
      ),
    );

    expect(find.text('Recommended'), findsNothing);
    expect(find.text('RS Hermina Batam'), findsOneWidget);
  });
}
