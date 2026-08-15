import 'package:freezed_annotation/freezed_annotation.dart';

import 'money.dart';
import 'time_window.dart';

part 'structured_intent.freezed.dart';
part 'structured_intent.g.dart';

/// How the backend resolved a trip request's structured intent.
enum IntentResolution {
  @JsonValue('MATCHED')
  matched,
  @JsonValue('NEEDS_CLARIFICATION')
  needsClarification,
  @JsonValue('UNSUPPORTED_SERVICE')
  unsupportedService,
  @JsonValue('OUT_OF_SCOPE')
  outOfScope,
}

/// Stay duration preference extracted from the patient's request.
enum StayType {
  @JsonValue('SAME_DAY')
  sameDay,
  @JsonValue('OVERNIGHT')
  overnight,
  @JsonValue('FLEXIBLE')
  flexible,
}

/// Non-critical patient preferences from the extracted intent.
@freezed
abstract class IntentPreferences with _$IntentPreferences {
  const factory IntentPreferences({
    String? language,
    @JsonKey(name: 'hotel_tier') String? hotelTier,
    List<String>? accessibility,
  }) = _IntentPreferences;

  factory IntentPreferences.fromJson(Map<String, dynamic> json) =>
      _$IntentPreferencesFromJson(json);
}

/// Structurally strict output from the language boundary
/// (`StructuredIntent` schema). Resolution-specific invariants are enforced
/// by the backend; nullable values are never filled with facts the patient
/// didn't state.
@freezed
abstract class StructuredIntent with _$StructuredIntent {
  const factory StructuredIntent({
    @JsonKey(name: 'schema_version') required String schemaVersion,
    required IntentResolution resolution,
    @JsonKey(name: 'intent_category') String? intentCategory,
    @JsonKey(name: 'requested_service_text')
    required String requestedServiceText,
    @JsonKey(name: 'service_code') String? serviceCode,
    @JsonKey(name: 'candidate_service_codes')
    required List<String> candidateServiceCodes,
    @JsonKey(name: 'origin_port') String? originPort,
    @JsonKey(name: 'date_window') DateWindow? dateWindow,
    @JsonKey(name: 'patient_count') int? patientCount,
    @JsonKey(name: 'companion_count') int? companionCount,
    @JsonKey(name: 'stay_type') StayType? stayType,
    Money? budget,
    required IntentPreferences preferences,
    @JsonKey(name: 'missing_fields') required List<String> missingFields,
    @JsonKey(name: 'clarification_question') String? clarificationQuestion,
    @JsonKey(name: 'out_of_scope_reason') String? outOfScopeReason,
    @JsonKey(name: 'unsupported_reason') String? unsupportedReason,
  }) = _StructuredIntent;

  factory StructuredIntent.fromJson(Map<String, dynamic> json) =>
      _$StructuredIntentFromJson(json);
}

/// Corrections the patient supplies to fix extracted intent
/// (`IntentCorrections` schema). At least one property must be present.
@freezed
abstract class IntentCorrections with _$IntentCorrections {
  const factory IntentCorrections({
    @JsonKey(name: 'service_code') String? serviceCode,
    @JsonKey(name: 'origin_port') String? originPort,
    @JsonKey(name: 'date_window') DateWindow? dateWindow,
    @JsonKey(name: 'patient_count') int? patientCount,
    @JsonKey(name: 'companion_count') int? companionCount,
    @JsonKey(name: 'stay_type') StayType? stayType,
    Money? budget,
    IntentPreferences? preferences,
  }) = _IntentCorrections;

  factory IntentCorrections.fromJson(Map<String, dynamic> json) =>
      _$IntentCorrectionsFromJson(json);
}

/// Body for `PATCH /v1/trip-requests/{id}/intent`: either an [answer] to the
/// clarification question and/or [corrections] to the extracted intent. At
/// least one of the two must be supplied.
@freezed
abstract class AmendIntentRequest with _$AmendIntentRequest {
  const factory AmendIntentRequest({
    String? answer,
    IntentCorrections? corrections,
  }) = _AmendIntentRequest;

  factory AmendIntentRequest.fromJson(Map<String, dynamic> json) =>
      _$AmendIntentRequestFromJson(json);
}
