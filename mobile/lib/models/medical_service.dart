import 'package:freezed_annotation/freezed_annotation.dart';

part 'medical_service.freezed.dart';
part 'medical_service.g.dart';

/// A supported, provider-authored medical service (`MedicalService` schema).
/// Only `HOSPITAL` providers may own a medical-service capability.
@freezed
abstract class MedicalService with _$MedicalService {
  const factory MedicalService({
    required String code,
    required String name,
    required String category,
    String? description,
    required bool active,
    required bool synthetic,
    required String source,
  }) = _MedicalService;

  factory MedicalService.fromJson(Map<String, dynamic> json) =>
      _$MedicalServiceFromJson(json);
}

/// Response body of `GET /v1/medical-services`.
@freezed
abstract class MedicalServiceListResponse with _$MedicalServiceListResponse {
  const factory MedicalServiceListResponse({
    required List<MedicalService> services,
  }) = _MedicalServiceListResponse;

  factory MedicalServiceListResponse.fromJson(Map<String, dynamic> json) =>
      _$MedicalServiceListResponseFromJson(json);
}
