// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'trip_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_TripRequest _$TripRequestFromJson(Map<String, dynamic> json) => _TripRequest(
  id: json['id'] as String,
  status: $enumDecode(_$TripRequestStatusEnumMap, json['status']),
  intent: StructuredIntent.fromJson(json['intent'] as Map<String, dynamic>),
  planningRevision: (json['planning_revision'] as num).toInt(),
  referenceCurrency: json['reference_currency'] as String,
  journeyId: json['journey_id'] as String?,
  createdAt: DateTime.parse(json['created_at'] as String),
  updatedAt: DateTime.parse(json['updated_at'] as String),
);

Map<String, dynamic> _$TripRequestToJson(_TripRequest instance) =>
    <String, dynamic>{
      'id': instance.id,
      'status': _$TripRequestStatusEnumMap[instance.status]!,
      'intent': instance.intent,
      'planning_revision': instance.planningRevision,
      'reference_currency': instance.referenceCurrency,
      'journey_id': instance.journeyId,
      'created_at': instance.createdAt.toIso8601String(),
      'updated_at': instance.updatedAt.toIso8601String(),
    };

const _$TripRequestStatusEnumMap = {
  TripRequestStatus.draft: 'DRAFT',
  TripRequestStatus.parsingIntent: 'PARSING_INTENT',
  TripRequestStatus.needsInput: 'NEEDS_INPUT',
  TripRequestStatus.unsupportedService: 'UNSUPPORTED_SERVICE',
  TripRequestStatus.outOfScope: 'OUT_OF_SCOPE',
  TripRequestStatus.planning: 'PLANNING',
  TripRequestStatus.noMatch: 'NO_MATCH',
  TripRequestStatus.planReady: 'PLAN_READY',
  TripRequestStatus.confirming: 'CONFIRMING',
  TripRequestStatus.active: 'ACTIVE',
  TripRequestStatus.confirmationFailed: 'CONFIRMATION_FAILED',
  TripRequestStatus.manualReview: 'MANUAL_REVIEW',
};

_CreateTripRequest _$CreateTripRequestFromJson(Map<String, dynamic> json) =>
    _CreateTripRequest(
      prompt: json['prompt'] as String,
      locale: json['locale'] as String,
    );

Map<String, dynamic> _$CreateTripRequestToJson(_CreateTripRequest instance) =>
    <String, dynamic>{'prompt': instance.prompt, 'locale': instance.locale};
