// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'medical_service.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_MedicalService _$MedicalServiceFromJson(Map<String, dynamic> json) =>
    _MedicalService(
      code: json['code'] as String,
      name: json['name'] as String,
      category: json['category'] as String,
      description: json['description'] as String?,
      active: json['active'] as bool,
      synthetic: json['synthetic'] as bool,
      source: json['source'] as String,
    );

Map<String, dynamic> _$MedicalServiceToJson(_MedicalService instance) =>
    <String, dynamic>{
      'code': instance.code,
      'name': instance.name,
      'category': instance.category,
      'description': instance.description,
      'active': instance.active,
      'synthetic': instance.synthetic,
      'source': instance.source,
    };

_MedicalServiceListResponse _$MedicalServiceListResponseFromJson(
  Map<String, dynamic> json,
) => _MedicalServiceListResponse(
  services: (json['services'] as List<dynamic>)
      .map((e) => MedicalService.fromJson(e as Map<String, dynamic>))
      .toList(),
);

Map<String, dynamic> _$MedicalServiceListResponseToJson(
  _MedicalServiceListResponse instance,
) => <String, dynamic>{'services': instance.services};
