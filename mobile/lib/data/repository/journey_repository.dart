import 'package:mobile/models/journey.dart';
import 'package:mobile/models/medical_service.dart';
import 'package:mobile/models/plan_option.dart';
import 'package:mobile/models/structured_intent.dart';

/// Thrown when a journey-orchestration operation fails (e.g. unsupported
/// service, planning failure, network error).
class JourneyException implements Exception {
  const JourneyException(this.message, {this.code});

  /// Human-readable message for display in the chat UI.
  final String message;

  /// Machine code from the API error envelope, when available.
  final String? code;

  @override
  String toString() => 'JourneyException($code): $message';
}

/// Contract for journey orchestration, agnostic of the transport.
///
/// The real implementation ([JourneyRepositoryImpl]) talks to the core API
/// over [JourneyApi]; [FakeJourneyRepository] is the in-memory stand-in used
/// while the backend isn't implemented. Both implement this interface so they
/// can be swapped via dependency injection (see
/// `application/journey/providers.dart`).
abstract class JourneyRepository {
  /// Create a trip request and extract structured intent.
  Future<TripRequestDetail> createTripRequest({
    required String prompt,
    required String locale,
  });

  /// Answer a clarification question and/or correct extracted intent.
  Future<TripRequestDetail> amendIntent({
    required String tripRequestId,
    String? answer,
    IntentCorrections? corrections,
  });

  /// Generate at most two feasible cross-provider plans.
  Future<PlanningResult> generatePlans({required String tripRequestId});

  /// Confirm one plan option, creating the journey.
  Future<JourneyDetail> confirmPlanOption({required String planOptionId});

  /// Retrieve the active itinerary and version history for a journey.
  Future<JourneyDetail> getJourneyItinerary({required String journeyId});

  /// List active medical services the platform can orchestrate.
  Future<MedicalServiceListResponse> listMedicalServices();
}
