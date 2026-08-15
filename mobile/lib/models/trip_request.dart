import 'package:freezed_annotation/freezed_annotation.dart';

import 'structured_intent.dart';

part 'trip_request.freezed.dart';
part 'trip_request.g.dart';

/// Lifecycle of a planning request before it becomes a journey
/// (`TripRequestStatus` schema).
enum TripRequestStatus {
  @JsonValue('DRAFT')
  draft,
  @JsonValue('PARSING_INTENT')
  parsingIntent,
  @JsonValue('NEEDS_INPUT')
  needsInput,
  @JsonValue('UNSUPPORTED_SERVICE')
  unsupportedService,
  @JsonValue('OUT_OF_SCOPE')
  outOfScope,
  @JsonValue('PLANNING')
  planning,
  @JsonValue('NO_MATCH')
  noMatch,
  @JsonValue('PLAN_READY')
  planReady,
  @JsonValue('CONFIRMING')
  confirming,
  @JsonValue('ACTIVE')
  active,
  @JsonValue('CONFIRMATION_FAILED')
  confirmationFailed,
  @JsonValue('MANUAL_REVIEW')
  manualReview,
}

/// A patient's planning request before confirmation (`TripRequest` schema).
@freezed
abstract class TripRequest with _$TripRequest {
  const factory TripRequest({
    required String id,
    required TripRequestStatus status,
    required StructuredIntent intent,
    @JsonKey(name: 'planning_revision') required int planningRevision,
    @JsonKey(name: 'reference_currency') required String referenceCurrency,
    @JsonKey(name: 'journey_id') String? journeyId,
    @JsonKey(name: 'created_at') required DateTime createdAt,
    @JsonKey(name: 'updated_at') required DateTime updatedAt,
  }) = _TripRequest;

  factory TripRequest.fromJson(Map<String, dynamic> json) =>
      _$TripRequestFromJson(json);
}

/// Request body for `POST /v1/trip-requests`: the patient's natural-language
/// request plus an input-language hint.
@freezed
abstract class CreateTripRequest with _$CreateTripRequest {
  const factory CreateTripRequest({
    required String prompt,
    required String locale,
  }) = _CreateTripRequest;

  factory CreateTripRequest.fromJson(Map<String, dynamic> json) =>
      _$CreateTripRequestFromJson(json);
}
