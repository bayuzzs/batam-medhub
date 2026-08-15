// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'structured_intent.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_IntentPreferences _$IntentPreferencesFromJson(Map<String, dynamic> json) =>
    _IntentPreferences(
      language: json['language'] as String?,
      hotelTier: json['hotel_tier'] as String?,
      accessibility: (json['accessibility'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList(),
    );

Map<String, dynamic> _$IntentPreferencesToJson(_IntentPreferences instance) =>
    <String, dynamic>{
      'language': instance.language,
      'hotel_tier': instance.hotelTier,
      'accessibility': instance.accessibility,
    };

_StructuredIntent _$StructuredIntentFromJson(Map<String, dynamic> json) =>
    _StructuredIntent(
      schemaVersion: json['schema_version'] as String,
      resolution: $enumDecode(_$IntentResolutionEnumMap, json['resolution']),
      intentCategory: json['intent_category'] as String?,
      requestedServiceText: json['requested_service_text'] as String,
      serviceCode: json['service_code'] as String?,
      candidateServiceCodes: (json['candidate_service_codes'] as List<dynamic>)
          .map((e) => e as String)
          .toList(),
      originPort: json['origin_port'] as String?,
      dateWindow: json['date_window'] == null
          ? null
          : DateWindow.fromJson(json['date_window'] as Map<String, dynamic>),
      patientCount: (json['patient_count'] as num?)?.toInt(),
      companionCount: (json['companion_count'] as num?)?.toInt(),
      stayType: $enumDecodeNullable(_$StayTypeEnumMap, json['stay_type']),
      budget: json['budget'] == null
          ? null
          : Money.fromJson(json['budget'] as Map<String, dynamic>),
      preferences: IntentPreferences.fromJson(
        json['preferences'] as Map<String, dynamic>,
      ),
      missingFields: (json['missing_fields'] as List<dynamic>)
          .map((e) => e as String)
          .toList(),
      clarificationQuestion: json['clarification_question'] as String?,
      outOfScopeReason: json['out_of_scope_reason'] as String?,
      unsupportedReason: json['unsupported_reason'] as String?,
    );

Map<String, dynamic> _$StructuredIntentToJson(_StructuredIntent instance) =>
    <String, dynamic>{
      'schema_version': instance.schemaVersion,
      'resolution': _$IntentResolutionEnumMap[instance.resolution]!,
      'intent_category': instance.intentCategory,
      'requested_service_text': instance.requestedServiceText,
      'service_code': instance.serviceCode,
      'candidate_service_codes': instance.candidateServiceCodes,
      'origin_port': instance.originPort,
      'date_window': instance.dateWindow,
      'patient_count': instance.patientCount,
      'companion_count': instance.companionCount,
      'stay_type': _$StayTypeEnumMap[instance.stayType],
      'budget': instance.budget,
      'preferences': instance.preferences,
      'missing_fields': instance.missingFields,
      'clarification_question': instance.clarificationQuestion,
      'out_of_scope_reason': instance.outOfScopeReason,
      'unsupported_reason': instance.unsupportedReason,
    };

const _$IntentResolutionEnumMap = {
  IntentResolution.matched: 'MATCHED',
  IntentResolution.needsClarification: 'NEEDS_CLARIFICATION',
  IntentResolution.unsupportedService: 'UNSUPPORTED_SERVICE',
  IntentResolution.outOfScope: 'OUT_OF_SCOPE',
};

const _$StayTypeEnumMap = {
  StayType.sameDay: 'SAME_DAY',
  StayType.overnight: 'OVERNIGHT',
  StayType.flexible: 'FLEXIBLE',
};

_IntentCorrections _$IntentCorrectionsFromJson(Map<String, dynamic> json) =>
    _IntentCorrections(
      serviceCode: json['service_code'] as String?,
      originPort: json['origin_port'] as String?,
      dateWindow: json['date_window'] == null
          ? null
          : DateWindow.fromJson(json['date_window'] as Map<String, dynamic>),
      patientCount: (json['patient_count'] as num?)?.toInt(),
      companionCount: (json['companion_count'] as num?)?.toInt(),
      stayType: $enumDecodeNullable(_$StayTypeEnumMap, json['stay_type']),
      budget: json['budget'] == null
          ? null
          : Money.fromJson(json['budget'] as Map<String, dynamic>),
      preferences: json['preferences'] == null
          ? null
          : IntentPreferences.fromJson(
              json['preferences'] as Map<String, dynamic>,
            ),
    );

Map<String, dynamic> _$IntentCorrectionsToJson(_IntentCorrections instance) =>
    <String, dynamic>{
      'service_code': instance.serviceCode,
      'origin_port': instance.originPort,
      'date_window': instance.dateWindow,
      'patient_count': instance.patientCount,
      'companion_count': instance.companionCount,
      'stay_type': _$StayTypeEnumMap[instance.stayType],
      'budget': instance.budget,
      'preferences': instance.preferences,
    };

_AmendIntentRequest _$AmendIntentRequestFromJson(Map<String, dynamic> json) =>
    _AmendIntentRequest(
      answer: json['answer'] as String?,
      corrections: json['corrections'] == null
          ? null
          : IntentCorrections.fromJson(
              json['corrections'] as Map<String, dynamic>,
            ),
    );

Map<String, dynamic> _$AmendIntentRequestToJson(_AmendIntentRequest instance) =>
    <String, dynamic>{
      'answer': instance.answer,
      'corrections': instance.corrections,
    };
