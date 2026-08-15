import 'package:freezed_annotation/freezed_annotation.dart';

import 'jwt_claims.dart';
import 'patient_profile.dart';

part 'auth_session.freezed.dart';
part 'auth_session.g.dart';

/// An authenticated session issued by the core API (`AuthSession` schema).
///
/// Holds the short-lived access JWT, the single-use refresh token, and the
/// patient profile. [accessExpiresAt] derives the access-token expiry from
/// the JWT `exp` claim, falling back to `receivedAt + expiresInSeconds`.
@freezed
abstract class AuthSession with _$AuthSession {
  const factory AuthSession({
    @JsonKey(name: 'token_type') required String tokenType,
    @JsonKey(name: 'access_token') required String accessToken,
    @JsonKey(name: 'refresh_token') required String refreshToken,
    @JsonKey(name: 'expires_in_seconds') required int expiresInSeconds,
    @JsonKey(name: 'refresh_expires_at') required DateTime refreshExpiresAt,
    required PatientProfile profile,
  }) = _AuthSession;

  factory AuthSession.fromJson(Map<String, dynamic> json) =>
      _$AuthSessionFromJson(json);
}

/// When the access token expires. Prefers the JWT `exp` claim (backend
/// authoritative); falls back to [receivedAt] + [expiresInSeconds] when the
/// token isn't a decodable JWT.
extension AuthSessionExpiry on AuthSession {
  DateTime accessExpiresAt({required DateTime receivedAt}) {
    final exp = tryDecodeJwtExp(accessToken);
    if (exp != null) {
      return DateTime.fromMillisecondsSinceEpoch(exp * 1000, isUtc: true);
    }
    return receivedAt.add(Duration(seconds: expiresInSeconds));
  }
}
