import 'package:mobile/models/journey.dart';
import 'package:mobile/models/medical_service.dart';
import 'package:mobile/models/plan_option.dart';
import 'package:mobile/models/structured_intent.dart';

/// Contract for the core API's journey-orchestration operations
/// (`/v1/trip-requests`, `/v1/plan-options`, `/v1/journeys`,
/// `/v1/medical-services`). Agnostic of the transport; [DioJourneyApi] is the
/// real HTTP implementation.
abstract class JourneyApi {
  /// `POST /v1/trip-requests` — create a trip request and extract structured
  /// intent from the patient's natural-language [prompt].
  Future<TripRequestDetail> createTripRequest({
    required String prompt,
    required String locale,
  });

  /// `PATCH /v1/trip-requests/{tripRequestId}/intent` — answer a
  /// clarification question and/or correct the extracted intent.
  Future<TripRequestDetail> amendIntent({
    required String tripRequestId,
    String? answer,
    IntentCorrections? corrections,
  });

  /// `POST /v1/trip-requests/{tripRequestId}/plans` — generate at most two
  /// feasible cross-provider plans.
  Future<PlanningResult> generatePlans({required String tripRequestId});

  /// `POST /v1/plan-options/{planOptionId}/confirm` — approve and confirm one
  /// plan option, creating the journey.
  Future<JourneyDetail> confirmPlanOption({required String planOptionId});

  /// `GET /v1/journeys/{journeyId}/itinerary` — retrieve the active itinerary
  /// and immutable version history.
  Future<JourneyDetail> getJourneyItinerary({required String journeyId});

  /// `GET /v1/medical-services` — list active medical services the platform
  /// can orchestrate.
  Future<MedicalServiceListResponse> listMedicalServices();
}
