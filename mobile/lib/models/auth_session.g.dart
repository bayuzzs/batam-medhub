// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'auth_session.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_AuthSession _$AuthSessionFromJson(Map<String, dynamic> json) => _AuthSession(
  tokenType: json['token_type'] as String,
  accessToken: json['access_token'] as String,
  refreshToken: json['refresh_token'] as String,
  expiresInSeconds: (json['expires_in_seconds'] as num).toInt(),
  refreshExpiresAt: DateTime.parse(json['refresh_expires_at'] as String),
  profile: PatientProfile.fromJson(json['profile'] as Map<String, dynamic>),
);

Map<String, dynamic> _$AuthSessionToJson(_AuthSession instance) =>
    <String, dynamic>{
      'token_type': instance.tokenType,
      'access_token': instance.accessToken,
      'refresh_token': instance.refreshToken,
      'expires_in_seconds': instance.expiresInSeconds,
      'refresh_expires_at': instance.refreshExpiresAt.toIso8601String(),
      'profile': instance.profile,
    };
