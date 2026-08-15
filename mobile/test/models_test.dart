import 'dart:convert';

import 'package:flutter_test/flutter_test.dart';

import 'package:mobile/models/auth_session.dart';
import 'package:mobile/models/jwt_claims.dart';

/// Builds a JWT-shaped token with the given `exp` claim.
String jwtWithExp(int exp) {
  final header = base64Url.encode(utf8.encode('{"alg":"HS256","typ":"JWT"}'));
  final payload = base64Url.encode(
    utf8.encode(jsonEncode({'exp': exp, 'sub': 'patient-000001'})),
  );
  return '$header.$payload.signature';
}

void main() {
  group('tryDecodeJwtExp', () {
    test('decodes the exp claim from a JWT payload', () {
      expect(tryDecodeJwtExp(jwtWithExp(1700000000)), 1700000000);
    });

    test('returns null for a non-JWT (opaque) token', () {
      expect(tryDecodeJwtExp('opaque-access-token'), isNull);
    });

    test('returns null for a malformed payload', () {
      final bad = 'eyJhbGciOiJIUzI1NiJ9.not-base64url.sig';
      expect(tryDecodeJwtExp(bad), isNull);
    });

    test('returns null when exp is missing or not an integer', () {
      final noExp =
          'eyJhbGciOiJIUzI1NiJ9.'
          '${base64Url.encode(utf8.encode('{"sub":"p1"}'))}.sig';
      expect(tryDecodeJwtExp(noExp), isNull);
    });
  });

  group('AuthSession', () {
    final profileJson = {
      'id': 'patient-000001',
      'full_name': 'Rina Tan',
      'email': 'rina.tan@example.test',
      'preferred_currency': 'SGD',
      'synthetic': true,
      'created_at': '2026-01-01T00:00:00.000Z',
      'updated_at': '2026-01-01T00:00:00.000Z',
    };

    test('fromJson maps the OpenAPI contract fields', () {
      final session = AuthSession.fromJson({
        'token_type': 'Bearer',
        'access_token': 'access-token',
        'refresh_token': 'refresh-token',
        'expires_in_seconds': 900,
        'refresh_expires_at': '2026-02-01T00:00:00.000Z',
        'profile': profileJson,
      });

      expect(session.tokenType, 'Bearer');
      expect(session.accessToken, 'access-token');
      expect(session.refreshToken, 'refresh-token');
      expect(session.expiresInSeconds, 900);
      expect(session.profile.fullName, 'Rina Tan');
      expect(session.profile.preferredCurrency, 'SGD');
      expect(session.profile.synthetic, isTrue);
    });

    test('accessExpiresAt prefers the JWT exp claim', () {
      final session = AuthSession.fromJson({
        'token_type': 'Bearer',
        'access_token': jwtWithExp(1700000000),
        'refresh_token': 'refresh-token',
        'expires_in_seconds': 900,
        'refresh_expires_at': '2026-02-01T00:00:00.000Z',
        'profile': profileJson,
      });

      final receivedAt = DateTime.utc(2023, 11, 1);
      expect(
        session.accessExpiresAt(receivedAt: receivedAt),
        DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true),
      );
    });

    test('accessExpiresAt falls back to receivedAt + expiresInSeconds', () {
      final session = AuthSession.fromJson({
        'token_type': 'Bearer',
        'access_token': 'opaque-access-token', // not a JWT
        'refresh_token': 'refresh-token',
        'expires_in_seconds': 60,
        'refresh_expires_at': '2026-02-01T00:00:00.000Z',
        'profile': profileJson,
      });

      final receivedAt = DateTime.utc(2026, 1, 1);
      expect(
        session.accessExpiresAt(receivedAt: receivedAt),
        receivedAt.add(const Duration(seconds: 60)),
      );
    });
  });
}
