// Tests for the fake journey repository's in-memory state machine, which
// mirrors the golden examples in `specs/examples/core/`.

import 'package:flutter_test/flutter_test.dart';

import 'package:mobile/data/repository/fake_journey_repository.dart';
import 'package:mobile/data/repository/journey_repository.dart';
import 'package:mobile/models/journey.dart';
import 'package:mobile/models/plan_option.dart';
import 'package:mobile/models/structured_intent.dart';
import 'package:mobile/models/trip_request.dart';

void main() {
  group('FakeJourneyRepository', () {
    final repository = FakeJourneyRepository(delay: Duration.zero);

    test('create → amend → plans → confirm walks the state machine', () async {
      // 1. Create a trip request → asks for clarification.
      final created = await repository.createTripRequest(
        prompt: 'I need a medical check-up in Batam',
        locale: 'en',
      );
      expect(created.tripRequest.id, 'trip-000002');
      expect(created.tripRequest.status, TripRequestStatus.needsInput);
      expect(
        created.tripRequest.intent.resolution,
        IntentResolution.needsClarification,
      );
      expect(created.tripRequest.intent.clarificationQuestion, isNotNull);

      // 2. Answer the clarification → matched, ready for planning.
      final answered = await repository.amendIntent(
        tripRequestId: created.tripRequest.id,
        answer: 'Basic check-up, next Friday',
      );
      expect(answered.tripRequest.id, 'trip-000001');
      expect(answered.tripRequest.status, TripRequestStatus.planning);
      expect(answered.tripRequest.intent.resolution, IntentResolution.matched);
      expect(answered.tripRequest.intent.serviceCode, 'MCU_BASIC');

      final plans = await repository.generatePlans(
        tripRequestId: answered.tripRequest.id,
      );
      expect(plans.tripRequest.id, 'trip-000001');
      expect(plans.tripRequest.status, TripRequestStatus.planReady);
      expect(plans.options.length, 2);
      expect(plans.options.first.rank, 1);
      expect(plans.options.first.status, PlanOptionStatus.proposed);
      expect(plans.options.first.totalPrice.displayTotal.amountMinor, 25190);

      // 3. Confirm the top option → active journey with a full itinerary.
      final journey = await repository.confirmPlanOption(
        planOptionId: plans.options.first.id,
      );
      expect(journey.journey.id, 'journey-000001');
      expect(journey.journey.status, JourneyStatus.active);
      expect(journey.activeItinerary.items, isNotEmpty);
      expect(
        journey.activeItinerary.items.where(
          (i) => i.itemType == ItemType.hospitalAppointment,
        ),
        isNotEmpty,
      );
    });

    test('getJourneyItinerary returns the confirmed journey', () async {
      await repository.createTripRequest(prompt: 'check-up', locale: 'en');
      final answered = await repository.amendIntent(
        tripRequestId: 'trip-000002',
        answer: 'Basic',
      );
      final plans = await repository.generatePlans(
        tripRequestId: answered.tripRequest.id,
      );
      await repository.confirmPlanOption(planOptionId: plans.options.first.id);

      final journey = await repository.getJourneyItinerary(
        journeyId: 'journey-000001',
      );
      expect(journey.journey.id, 'journey-000001');
    });

    test('listMedicalServices returns the four synthetic services', () async {
      final response = await repository.listMedicalServices();
      expect(response.services.length, 4);
      expect(response.services.map((s) => s.code), [
        'MCU_BASIC',
        'MCU_COMPREHENSIVE',
        'DENTAL_CHECKUP',
        'EYE_SCREENING',
      ]);
      expect(response.services.every((s) => s.synthetic), isTrue);
      expect(response.services.every((s) => s.source == 'MOCK'), isTrue);
    });

    test('amendIntent rejects an unknown trip request', () async {
      expect(
        () => repository.amendIntent(
          tripRequestId: 'trip-unknown',
          answer: 'Basic',
        ),
        throwsA(
          isA<JourneyException>().having((e) => e.code, 'code', 'NOT_FOUND'),
        ),
      );
    });
  });
}
