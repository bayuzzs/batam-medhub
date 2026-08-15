import 'package:freezed_annotation/freezed_annotation.dart';

part 'patient_profile.freezed.dart';
part 'patient_profile.g.dart';

/// The authenticated patient's profile, as returned by the core API.
///
/// Field names follow the OpenAPI contract (`full_name`, `preferred_currency`,
/// `created_at`, `updated_at`) via [JsonKey].
@freezed
abstract class PatientProfile with _$PatientProfile {
  const factory PatientProfile({
    required String id,
    @JsonKey(name: 'full_name') required String fullName,
    required String email,
    @JsonKey(name: 'preferred_currency') required String preferredCurrency,
    required bool synthetic,
    @JsonKey(name: 'created_at') required DateTime createdAt,
    @JsonKey(name: 'updated_at') required DateTime updatedAt,
  }) = _PatientProfile;

  factory PatientProfile.fromJson(Map<String, dynamic> json) =>
      _$PatientProfileFromJson(json);
}
