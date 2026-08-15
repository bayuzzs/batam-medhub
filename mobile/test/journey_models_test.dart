// Model tests for the journey data layer, parsed from the golden fixtures
// copied from `specs/examples/core/` into `test/fixtures/core/`.

import 'dart:convert';
import 'dart:io';

import 'package:flutter_test/flutter_test.dart';

import 'package:mobile/models/journey.dart';
import 'package:mobile/models/medical_service.dart';
import 'package:mobile/models/plan_option.dart';
import 'package:mobile/models/structured_intent.dart';
import 'package:mobile/models/trip_request.dart';

Map<String, dynamic> _load(String name) =>
    jsonDecode(File('test/fixtures/core/$name').readAsStringSync())
        as Map<String, dynamic>;

/// Serializes [value] through JSON and back to a plain map. Nested freezed
/// objects have `explicitToJson: false`, so `toJson()` emits object graphs
/// that only become plain maps after `jsonEncode`; this mirrors how the app
/// actually sends payloads over the wire.
Map<String, dynamic> _toJsonMap(Object value) =>
    (jsonDecode(jsonEncode(value)) as Map).cast<String, dynamic>();

void main() {
  group('TripRequestDetail from golden fixtures', () {
    test('matched-trip-request parses a MATCHED trip with no options', () {
      final detail = TripRequestDetail.fromJson(
        _load('matched-trip-request.json'),
      );
      final trip = detail.tripRequest;
      expect(trip.id, 'trip-000001');
      expect(trip.status, TripRequestStatus.planning);
      expect(trip.intent.resolution, IntentResolution.matched);
      expect(trip.intent.serviceCode, 'MCU_BASIC');
      expect(trip.intent.budget!.amountMinor, 40000);
      expect(trip.intent.budget!.currency, 'SGD');
      expect(trip.intent.dateWindow!.from, DateTime(2026, 8, 22));
      expect(detail.planOptions, isEmpty);
    });

    test(
      'needs-clarification-trip-request parses missing fields + question',
      () {
        final detail = TripRequestDetail.fromJson(
          _load('needs-clarification-trip-request.json'),
        );
        final trip = detail.tripRequest;
        expect(trip.id, 'trip-000002');
        expect(trip.status, TripRequestStatus.needsInput);
        expect(trip.intent.resolution, IntentResolution.needsClarification);
        expect(trip.intent.candidateServiceCodes, [
          'MCU_BASIC',
          'MCU_COMPREHENSIVE',
        ]);
        expect(
          trip.intent.missingFields,
          containsAll(['service_code', 'date_window']),
        );
        expect(
          trip.intent.clarificationQuestion,
          contains('basic or comprehensive'),
        );
        expect(detail.planOptions, isEmpty);
      },
    );
  });

  group('PlanningResult from golden fixture', () {
    test('plan-result parses the option with legs and price totals', () {
      final result = PlanningResult.fromJson(_load('plan-result.json'));
      expect(result.tripRequest.id, 'trip-000001');
      expect(result.tripRequest.status, TripRequestStatus.planReady);
      expect(result.options.length, 1);

      final option = result.options.first;
      expect(option.rank, 1);
      expect(option.status, PlanOptionStatus.proposed);
      expect(option.totalPrice.sourceTotals, isNotEmpty);
      expect(option.totalPrice.displayTotal.currency, 'SGD');
      expect(option.totalPrice.displayTotal.amountMinor, greaterThan(0));

      // First leg is the outbound ferry; the plan also contains a hospital
      // appointment with a provider and time zone.
      expect(option.items.first.itemType, ItemType.ferryOutbound);
      final hospital = option.items.firstWhere(
        (i) => i.itemType == ItemType.hospitalAppointment,
      );
      expect(hospital.providerId, 'hospital-demo-01');
      expect(hospital.timeWindow.startTimeZone, 'Asia/Jakarta');
      expect(hospital.price!.display.amountMinor, 12658);
    });
  });

  group('JourneyDetail from golden fixture', () {
    test('active-journey-v1 parses journey, active itinerary, and history', () {
      final detail = JourneyDetail.fromJson(_load('active-journey-v1.json'));
      expect(detail.journey.id, 'journey-000001');
      expect(detail.journey.status, JourneyStatus.active);
      expect(detail.activeItinerary.id, 'itinerary-000001-v1');
      expect(detail.activeItinerary.version, 1);
      expect(detail.activeItinerary.items, isNotEmpty);
      expect(detail.activeItinerary.totalPrice.displayTotal.currency, 'SGD');
      expect(detail.itineraryVersions, isNotEmpty);
      expect(
        detail.itineraryVersions.first.status,
        ItineraryVersionStatus.active,
      );
    });
  });

  group('MedicalServiceListResponse from golden fixture', () {
    test('medical-services parses the four synthetic services', () {
      final response = MedicalServiceListResponse.fromJson(
        _load('medical-services.json'),
      );
      expect(response.services.map((s) => s.code), [
        'MCU_BASIC',
        'MCU_COMPREHENSIVE',
        'DENTAL_CHECKUP',
        'EYE_SCREENING',
      ]);
      expect(response.services.first.active, isTrue);
      expect(response.services.first.synthetic, isTrue);
      expect(response.services.first.source, 'MOCK');
    });
  });

  group('Round trips (toJson → fromJson)', () {
    test('TripRequestDetail survives a round trip', () {
      final original = TripRequestDetail.fromJson(
        _load('needs-clarification-trip-request.json'),
      );
      expect(
        TripRequestDetail.fromJson(_toJsonMap(original.toJson())),
        original,
      );
    });

    test('PlanningResult survives a round trip', () {
      final original = PlanningResult.fromJson(_load('plan-result.json'));
      expect(PlanningResult.fromJson(_toJsonMap(original.toJson())), original);
    });

    test('JourneyDetail survives a round trip', () {
      final original = JourneyDetail.fromJson(_load('active-journey-v1.json'));
      expect(JourneyDetail.fromJson(_toJsonMap(original.toJson())), original);
    });

    test('MedicalServiceListResponse survives a round trip', () {
      final original = MedicalServiceListResponse.fromJson(
        _load('medical-services.json'),
      );
      expect(
        MedicalServiceListResponse.fromJson(_toJsonMap(original.toJson())),
        original,
      );
    });
  });
}
