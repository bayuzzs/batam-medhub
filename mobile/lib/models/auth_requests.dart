import 'package:json_annotation/json_annotation.dart';

part 'auth_requests.g.dart';

/// Request body for `POST /v1/auth/login`.
@JsonSerializable()
class LoginRequest {
  const LoginRequest({required this.email, required this.password});

  final String email;
  final String password;

  factory LoginRequest.fromJson(Map<String, dynamic> json) =>
      _$LoginRequestFromJson(json);

  Map<String, dynamic> toJson() => _$LoginRequestToJson(this);
}

/// Request body for `POST /v1/auth/register`.
///
/// `confirm_password` is a mobile-only validation field and must never be sent
/// to the backend. `preferred_currency` is optional; the backend defaults to
/// `SGD` when omitted.
@JsonSerializable()
class RegisterPatientRequest {
  const RegisterPatientRequest({
    required this.fullName,
    required this.email,
    required this.password,
    this.preferredCurrency,
  });

  @JsonKey(name: 'full_name')
  final String fullName;

  final String email;
  final String password;

  @JsonKey(name: 'preferred_currency')
  final String? preferredCurrency;

  factory RegisterPatientRequest.fromJson(Map<String, dynamic> json) =>
      _$RegisterPatientRequestFromJson(json);

  Map<String, dynamic> toJson() => _$RegisterPatientRequestToJson(this);
}

/// Request body for `POST /v1/auth/refresh` and `POST /v1/auth/logout`.
@JsonSerializable()
class RefreshTokenRequest {
  const RefreshTokenRequest({required this.refreshToken});

  @JsonKey(name: 'refresh_token')
  final String refreshToken;

  factory RefreshTokenRequest.fromJson(Map<String, dynamic> json) =>
      _$RefreshTokenRequestFromJson(json);

  Map<String, dynamic> toJson() => _$RefreshTokenRequestToJson(this);
}
