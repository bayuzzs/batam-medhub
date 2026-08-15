import 'package:mobile/models/journey.dart';
import 'package:mobile/models/medical_service.dart';
import 'package:mobile/models/money.dart';
import 'package:mobile/models/plan_option.dart';
import 'package:mobile/models/structured_intent.dart';
import 'package:mobile/models/time_window.dart';
import 'package:mobile/models/trip_request.dart';

import 'journey_repository.dart';

/// In-memory [JourneyRepository] used while the core backend isn't
/// implemented.
///
/// Mirrors the golden examples in `specs/examples/core/` (all synthetic,
/// `source: MOCK`): creating a trip request returns `NEEDS_CLARIFICATION`
/// with one question; amending intent returns a matched request; generating
/// plans returns two options; confirming one creates an active journey.
///
/// A small [delay] simulates network latency; pass [Duration.zero] in tests.
class FakeJourneyRepository implements JourneyRepository {
  FakeJourneyRepository({this.delay = const Duration(milliseconds: 300)});

  /// Simulated network latency for each operation.
  final Duration delay;

  TripRequest? _tripRequest;
  PlanningResult? _plans;
  JourneyDetail? _journey;

  static const _currency = 'SGD';

  @override
  Future<TripRequestDetail> createTripRequest({
    required String prompt,
    required String locale,
  }) async {
    await Future<void>.delayed(delay);
    _tripRequest = _buildClarificationTripRequest(prompt);
    return TripRequestDetail(tripRequest: _tripRequest!, planOptions: const []);
  }

  @override
  Future<TripRequestDetail> amendIntent({
    required String tripRequestId,
    String? answer,
    IntentCorrections? corrections,
  }) async {
    await Future<void>.delayed(delay);
    if (_tripRequest == null || _tripRequest!.id != tripRequestId) {
      throw const JourneyException('Unknown trip request', code: 'NOT_FOUND');
    }
    _tripRequest = _buildMatchedTripRequest(answer: answer);
    return TripRequestDetail(tripRequest: _tripRequest!, planOptions: const []);
  }

  @override
  Future<PlanningResult> generatePlans({required String tripRequestId}) async {
    await Future<void>.delayed(delay);
    if (_tripRequest == null || _tripRequest!.id != tripRequestId) {
      throw const JourneyException('Unknown trip request', code: 'NOT_FOUND');
    }
    final trip = _tripRequest!.copyWith(
      status: TripRequestStatus.planReady,
      planningRevision: 1,
    );
    _tripRequest = trip;
    _plans = PlanningResult(
      tripRequest: trip,
      options: _buildPlanOptions(),
      noMatchReasons: const [],
      providerWarnings: const [],
    );
    return _plans!;
  }

  @override
  Future<JourneyDetail> confirmPlanOption({
    required String planOptionId,
  }) async {
    await Future<void>.delayed(delay);
    if (_plans == null || !_plans!.options.any((o) => o.id == planOptionId)) {
      throw const JourneyException('Unknown plan option', code: 'NOT_FOUND');
    }
    _journey = _buildActiveJourney();
    return _journey!;
  }

  @override
  Future<JourneyDetail> getJourneyItinerary({required String journeyId}) async {
    await Future<void>.delayed(delay);
    if (_journey == null || _journey!.journey.id != journeyId) {
      throw const JourneyException('Unknown journey', code: 'NOT_FOUND');
    }
    return _journey!;
  }

  @override
  Future<MedicalServiceListResponse> listMedicalServices() async {
    await Future<void>.delayed(delay);
    return const MedicalServiceListResponse(
      services: [
        MedicalService(
          code: 'MCU_BASIC',
          name: 'Basic Medical Check-up',
          category: 'PREVENTIVE_CHECKUP',
          description: 'Synthetic basic planned check-up package for the demo.',
          active: true,
          synthetic: true,
          source: 'MOCK',
        ),
        MedicalService(
          code: 'MCU_COMPREHENSIVE',
          name: 'Comprehensive Medical Check-up',
          category: 'PREVENTIVE_CHECKUP',
          description:
              'Synthetic comprehensive planned check-up package for the demo.',
          active: true,
          synthetic: true,
          source: 'MOCK',
        ),
        MedicalService(
          code: 'DENTAL_CHECKUP',
          name: 'Dental Check-up',
          category: 'DENTAL_SCREENING',
          description:
              'Synthetic planned dental screening package for the demo.',
          active: true,
          synthetic: true,
          source: 'MOCK',
        ),
        MedicalService(
          code: 'EYE_SCREENING',
          name: 'Eye Screening',
          category: 'VISION_SCREENING',
          description:
              'Synthetic planned vision screening package for the demo.',
          active: true,
          synthetic: true,
          source: 'MOCK',
        ),
      ],
    );
  }

  // --- Fixture builders (mirror specs/examples/core/*.json) ----------------

  TripRequest _buildClarificationTripRequest(String prompt) {
    return TripRequest(
      id: 'trip-000002',
      status: TripRequestStatus.needsInput,
      intent: StructuredIntent(
        schemaVersion: '1.0',
        resolution: IntentResolution.needsClarification,
        intentCategory: 'PREVENTIVE_CHECKUP',
        requestedServiceText: prompt,
        serviceCode: null,
        candidateServiceCodes: const ['MCU_BASIC', 'MCU_COMPREHENSIVE'],
        originPort: 'HARBOURFRONT_SG',
        dateWindow: null,
        patientCount: 1,
        companionCount: 0,
        stayType: null,
        budget: null,
        preferences: const IntentPreferences(
          language: 'en',
          hotelTier: null,
          accessibility: [],
        ),
        missingFields: const ['service_code', 'date_window'],
        clarificationQuestion:
            'Would you like the basic or comprehensive check-up, and what '
            'date would you prefer?',
        outOfScopeReason: null,
        unsupportedReason: null,
      ),
      planningRevision: 0,
      referenceCurrency: _currency,
      journeyId: null,
      createdAt: DateTime.utc(2026, 8, 15, 8, 5),
      updatedAt: DateTime.utc(2026, 8, 15, 8, 5, 1),
    );
  }

  TripRequest _buildMatchedTripRequest({String? answer}) {
    return TripRequest(
      id: 'trip-000001',
      status: TripRequestStatus.planning,
      intent: StructuredIntent(
        schemaVersion: '1.0',
        resolution: IntentResolution.matched,
        intentCategory: 'PREVENTIVE_CHECKUP',
        requestedServiceText: answer ?? 'basic medical check-up',
        serviceCode: 'MCU_BASIC',
        candidateServiceCodes: const [],
        originPort: 'HARBOURFRONT_SG',
        dateWindow: DateWindow(
          from: DateTime.utc(2026, 8, 22),
          to: DateTime.utc(2026, 8, 22),
        ),
        patientCount: 1,
        companionCount: 1,
        stayType: StayType.sameDay,
        budget: const Money(amountMinor: 40000, currency: 'SGD'),
        preferences: const IntentPreferences(
          language: 'en',
          hotelTier: null,
          accessibility: [],
        ),
        missingFields: const [],
        clarificationQuestion: null,
        outOfScopeReason: null,
        unsupportedReason: null,
      ),
      planningRevision: 0,
      referenceCurrency: _currency,
      journeyId: null,
      createdAt: DateTime.utc(2026, 8, 15, 8, 0),
      updatedAt: DateTime.utc(2026, 8, 15, 8, 0, 1),
    );
  }

  List<PlanOption> _buildPlanOptions() {
    final ferryPrice = ConvertedMoney(
      source: Money(amountMinor: 5000, currency: 'SGD'),
      display: Money(amountMinor: 5000, currency: 'SGD'),
      fxRate: '1.000000',
      fxSource: 'DEMO_STATIC_2026_08',
      fxEffectiveAt: _fxDate,
      estimated: true,
    );
    final transferPrice = ConvertedMoney(
      source: Money(amountMinor: 15000000, currency: 'IDR'),
      display: Money(amountMinor: 1266, currency: 'SGD'),
      fxRate: '0.0000843882',
      fxSource: 'DEMO_STATIC_2026_08',
      fxEffectiveAt: _fxDate,
      estimated: true,
    );
    final hospitalPrice = ConvertedMoney(
      source: Money(amountMinor: 150000000, currency: 'IDR'),
      display: Money(amountMinor: 12658, currency: 'SGD'),
      fxRate: '0.0000843882',
      fxSource: 'DEMO_STATIC_2026_08',
      fxEffectiveAt: _fxDate,
      estimated: true,
    );

    final items = <PlanItem>[
      PlanItem(
        id: 'plan-item-ferry-out',
        itemType: ItemType.ferryOutbound,
        providerId: 'ferry-demo-01',
        externalOfferId: 'ferry-offer-hf-btm-20260822-0730',
        title: 'HarbourFront to Batam Centre',
        timeWindow: TimeWindow(
          startsAt: _ferryOutStart,
          endsAt: _ferryOutEnd,
          startTimeZone: 'Asia/Singapore',
          endTimeZone: 'Asia/Jakarta',
        ),
        originCode: 'HARBOURFRONT_SG',
        destinationCode: 'BATAM_CENTRE_ID',
        price: ferryPrice,
        offerExpiresAt: _offerExpiry,
        operationalNotes: ['Check in at least 60 minutes before departure.'],
        synthetic: true,
        source: 'MOCK',
      ),
      PlanItem(
        id: 'plan-item-arrival-buffer',
        itemType: ItemType.arrivalBuffer,
        providerId: null,
        externalOfferId: null,
        title: 'Immigration and arrival buffer',
        timeWindow: TimeWindow(
          startsAt: _arrivalBufferStart,
          endsAt: _arrivalBufferEnd,
          startTimeZone: 'Asia/Jakarta',
          endTimeZone: 'Asia/Jakarta',
        ),
        originCode: 'BATAM_CENTRE_ID',
        destinationCode: 'BATAM_CENTRE_ID',
        price: null,
        offerExpiresAt: null,
        operationalNotes: [],
        synthetic: true,
        source: 'MOCK',
      ),
      PlanItem(
        id: 'plan-item-transfer-out',
        itemType: ItemType.transportPickup,
        providerId: 'transport-demo-01',
        externalOfferId: 'transport-offer-btm-hospital-20260822-0825',
        title: 'Terminal pickup to hospital',
        timeWindow: TimeWindow(
          startsAt: _transferOutStart,
          endsAt: _transferOutEnd,
          startTimeZone: 'Asia/Jakarta',
          endTimeZone: 'Asia/Jakarta',
        ),
        originCode: 'BATAM_CENTRE_ID',
        destinationCode: 'HOSPITAL_DEMO_ID',
        price: transferPrice,
        offerExpiresAt: _offerExpiry,
        operationalNotes: [
          'Meet the driver at the signed arrival pickup point.',
        ],
        synthetic: true,
        source: 'MOCK',
      ),
      PlanItem(
        id: 'plan-item-hospital',
        itemType: ItemType.hospitalAppointment,
        providerId: 'hospital-demo-01',
        externalOfferId: 'hospital-offer-mcu-basic-20260822-1000',
        title: 'Basic Medical Check-up',
        timeWindow: TimeWindow(
          startsAt: _hospitalStart,
          endsAt: _hospitalEnd,
          startTimeZone: 'Asia/Jakarta',
          endTimeZone: 'Asia/Jakarta',
        ),
        originCode: 'HOSPITAL_DEMO_ID',
        destinationCode: 'HOSPITAL_DEMO_ID',
        price: hospitalPrice,
        offerExpiresAt: _offerExpiry,
        operationalNotes: [
          'Follow only the preparation instructions supplied by the hospital.',
        ],
        synthetic: true,
        source: 'MOCK',
      ),
      PlanItem(
        id: 'plan-item-transfer-return',
        itemType: ItemType.transportDropoff,
        providerId: 'transport-demo-01',
        externalOfferId: 'transport-offer-hospital-btm-20260822-1230',
        title: 'Hospital to terminal',
        timeWindow: TimeWindow(
          startsAt: _transferReturnStart,
          endsAt: _transferReturnEnd,
          startTimeZone: 'Asia/Jakarta',
          endTimeZone: 'Asia/Jakarta',
        ),
        originCode: 'HOSPITAL_DEMO_ID',
        destinationCode: 'BATAM_CENTRE_ID',
        price: transferPrice,
        offerExpiresAt: _offerExpiry,
        operationalNotes: [],
        synthetic: true,
        source: 'MOCK',
      ),
      PlanItem(
        id: 'plan-item-departure-buffer',
        itemType: ItemType.departureBuffer,
        providerId: null,
        externalOfferId: null,
        title: 'Return ferry check-in buffer',
        timeWindow: TimeWindow(
          startsAt: _departureBufferStart,
          endsAt: _departureBufferEnd,
          startTimeZone: 'Asia/Jakarta',
          endTimeZone: 'Asia/Jakarta',
        ),
        originCode: 'BATAM_CENTRE_ID',
        destinationCode: 'BATAM_CENTRE_ID',
        price: null,
        offerExpiresAt: null,
        operationalNotes: [],
        synthetic: true,
        source: 'MOCK',
      ),
      PlanItem(
        id: 'plan-item-ferry-return',
        itemType: ItemType.ferryReturn,
        providerId: 'ferry-demo-01',
        externalOfferId: 'ferry-offer-btm-hf-20260822-1430',
        title: 'Batam Centre to HarbourFront',
        timeWindow: TimeWindow(
          startsAt: _ferryReturnStart,
          endsAt: _ferryReturnEnd,
          startTimeZone: 'Asia/Jakarta',
          endTimeZone: 'Asia/Singapore',
        ),
        originCode: 'BATAM_CENTRE_ID',
        destinationCode: 'HARBOURFRONT_SG',
        price: ferryPrice,
        offerExpiresAt: _offerExpiry,
        operationalNotes: ['Check in at least 30 minutes before departure.'],
        synthetic: true,
        source: 'MOCK',
      ),
    ];

    const total = PriceSummary(
      sourceTotals: [
        Money(amountMinor: 10000, currency: 'SGD'),
        Money(amountMinor: 180000000, currency: 'IDR'),
      ],
      displayTotal: Money(amountMinor: 25190, currency: 'SGD'),
      estimated: true,
    );

    return [
      PlanOption(
        id: 'plan-000001',
        tripRequestId: 'trip-000001',
        planningRevision: 1,
        rank: 1,
        status: PlanOptionStatus.proposed,
        expiresAt: _offerExpiry,
        explanation: const [
          'The ferry arrives with 140 minutes for immigration, transfer, and '
              'the medical arrival buffer.',
          'Every required provider has capacity for the patient and companion.',
          'The return sailing leaves after the appointment and terminal '
              'cutoff buffer.',
        ],
        items: items,
        totalPrice: total,
      ),
      PlanOption(
        id: 'plan-000002',
        tripRequestId: 'trip-000001',
        planningRevision: 1,
        rank: 2,
        status: PlanOptionStatus.proposed,
        expiresAt: _offerExpiry,
        explanation: const [
          'A later morning ferry still fits the appointment window with a '
              'tighter buffer.',
        ],
        items: items,
        totalPrice: total,
      ),
    ];
  }

  JourneyDetail _buildActiveJourney() {
    final journey = Journey(
      id: 'journey-000001',
      tripRequestId: 'trip-000001',
      status: JourneyStatus.active,
      activeItineraryVersion: 1,
      createdAt: DateTime.utc(2026, 8, 15, 8, 10),
      updatedAt: DateTime.utc(2026, 8, 15, 8, 10),
    );

    final active = ItineraryVersion(
      id: 'itinerary-v1',
      journeyId: journey.id,
      version: 1,
      status: ItineraryVersionStatus.active,
      basedOnDisruptionId: null,
      totalPrice: const PriceSummary(
        sourceTotals: [
          Money(amountMinor: 10000, currency: 'SGD'),
          Money(amountMinor: 180000000, currency: 'IDR'),
        ],
        displayTotal: Money(amountMinor: 25190, currency: 'SGD'),
        estimated: true,
      ),
      items: _buildItineraryItems(),
      createdAt: DateTime.utc(2026, 8, 15, 8, 10),
    );

    return JourneyDetail(
      journey: journey,
      activeItinerary: active,
      itineraryVersions: [
        ItineraryVersionSummary(
          id: 'itinerary-v1',
          version: 1,
          status: ItineraryVersionStatus.active,
          basedOnDisruptionId: null,
          createdAt: _createdAt,
        ),
      ],
    );
  }

  List<ItineraryItem> _buildItineraryItems() {
    return [
      ItineraryItem(
        id: 'item-ferry-out',
        itemType: ItemType.ferryOutbound,
        providerId: 'ferry-demo-01',
        externalReservationId: 'reservation-ferry-out-001',
        title: 'HarbourFront to Batam Centre',
        status: ItineraryItemStatus.confirmed,
        timeWindow: TimeWindow(
          startsAt: _ferryOutStart,
          endsAt: _ferryOutEnd,
          startTimeZone: 'Asia/Singapore',
          endTimeZone: 'Asia/Jakarta',
        ),
        originCode: 'HARBOURFRONT_SG',
        destinationCode: 'BATAM_CENTRE_ID',
        price: null,
        operationalNotes: [],
        synthetic: true,
        source: 'MOCK',
      ),
      ItineraryItem(
        id: 'item-hospital',
        itemType: ItemType.hospitalAppointment,
        providerId: 'hospital-demo-01',
        externalReservationId: 'reservation-hospital-001',
        title: 'Basic Medical Check-up',
        status: ItineraryItemStatus.confirmed,
        timeWindow: TimeWindow(
          startsAt: _hospitalStart,
          endsAt: _hospitalEnd,
          startTimeZone: 'Asia/Jakarta',
          endTimeZone: 'Asia/Jakarta',
        ),
        originCode: 'HOSPITAL_DEMO_ID',
        destinationCode: 'HOSPITAL_DEMO_ID',
        price: null,
        operationalNotes: [],
        synthetic: true,
        source: 'MOCK',
      ),
      ItineraryItem(
        id: 'item-ferry-return',
        itemType: ItemType.ferryReturn,
        providerId: 'ferry-demo-01',
        externalReservationId: 'reservation-ferry-return-001',
        title: 'Batam Centre to HarbourFront',
        status: ItineraryItemStatus.confirmed,
        timeWindow: TimeWindow(
          startsAt: _ferryReturnStart,
          endsAt: _ferryReturnEnd,
          startTimeZone: 'Asia/Jakarta',
          endTimeZone: 'Asia/Singapore',
        ),
        originCode: 'BATAM_CENTRE_ID',
        destinationCode: 'HARBOURFRONT_SG',
        price: null,
        operationalNotes: [],
        synthetic: true,
        source: 'MOCK',
      ),
    ];
  }

  static final _fxDate = DateTime.utc(2026, 8, 1);
  static final _createdAt = DateTime.utc(2026, 8, 15, 8, 10);
  static final _offerExpiry = DateTime.utc(2026, 8, 20, 12);
  static final _ferryOutStart = DateTime.utc(2026, 8, 21, 23, 30);
  static final _ferryOutEnd = DateTime.utc(2026, 8, 22, 0, 40);
  static final _arrivalBufferStart = DateTime.utc(2026, 8, 22, 0, 40);
  static final _arrivalBufferEnd = DateTime.utc(2026, 8, 22, 1, 25);
  static final _transferOutStart = DateTime.utc(2026, 8, 22, 1, 25);
  static final _transferOutEnd = DateTime.utc(2026, 8, 22, 2, 10);
  static final _hospitalStart = DateTime.utc(2026, 8, 22, 3);
  static final _hospitalEnd = DateTime.utc(2026, 8, 22, 5);
  static final _transferReturnStart = DateTime.utc(2026, 8, 22, 5, 30);
  static final _transferReturnEnd = DateTime.utc(2026, 8, 22, 6, 15);
  static final _departureBufferStart = DateTime.utc(2026, 8, 22, 6, 15);
  static final _departureBufferEnd = DateTime.utc(2026, 8, 22, 7);
  static final _ferryReturnStart = DateTime.utc(2026, 8, 22, 7, 30);
  static final _ferryReturnEnd = DateTime.utc(2026, 8, 22, 8, 40);
}
