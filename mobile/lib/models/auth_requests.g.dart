// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'auth_requests.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

LoginRequest _$LoginRequestFromJson(Map<String, dynamic> json) => LoginRequest(
  email: json['email'] as String,
  password: json['password'] as String,
);

Map<String, dynamic> _$LoginRequestToJson(LoginRequest instance) =>
    <String, dynamic>{'email': instance.email, 'password': instance.password};

RegisterPatientRequest _$RegisterPatientRequestFromJson(
  Map<String, dynamic> json,
) => RegisterPatientRequest(
  fullName: json['full_name'] as String,
  email: json['email'] as String,
  password: json['password'] as String,
  preferredCurrency: json['preferred_currency'] as String?,
);

Map<String, dynamic> _$RegisterPatientRequestToJson(
  RegisterPatientRequest instance,
) => <String, dynamic>{
  'full_name': instance.fullName,
  'email': instance.email,
  'password': instance.password,
  'preferred_currency': instance.preferredCurrency,
};

RefreshTokenRequest _$RefreshTokenRequestFromJson(Map<String, dynamic> json) =>
    RefreshTokenRequest(refreshToken: json['refresh_token'] as String);

Map<String, dynamic> _$RefreshTokenRequestToJson(
  RefreshTokenRequest instance,
) => <String, dynamic>{'refresh_token': instance.refreshToken};
