import 'dart:convert';

import 'package:mobile/data/service/token_store.dart';
import 'package:mobile/models/auth_session.dart';
import 'package:mobile/models/patient_profile.dart';

import 'auth_repository.dart';

/// In-memory [AuthRepository] used while the core backend isn't implemented.
///
/// Accepts any well-formed credentials and returns a synthetic fixture
/// session whose access token is a JWT-shaped token carrying an `exp` claim
/// (so the JWT-expiry path in [AuthSession.accessExpiresAt] is exercised).
///
/// `refresh` rotates the session — the presented refresh token is consumed
/// (single-use) and a new token pair is issued — and the session is persisted
/// through [TokenStore] so `restore()` works within a process.
class FakeAuthRepository implements AuthRepository {
  FakeAuthRepository({this._tokenStore, this.accessTtlSeconds = 900});

  final TokenStore? _tokenStore;

  /// Access-token lifetime issued by the fake. Set short (e.g. 5) to observe
  /// auto-refresh in manual QA.
  final int accessTtlSeconds;

  /// Documented demo user for manual QA. Any well-formed credentials also
  /// authenticate — the fake performs no real password check.
  static const String demoEmail = 'rina.tan@example.test';
  static const String demoPassword = 'Demo-Only-Password-2026!';

  final Set<String> _issuedRefreshTokens = {};
  int _counter = 0;

  @override
  Future<AuthSession> login({
    required String email,
    required String password,
  }) async {
    final session = _buildSession(
      fullName: 'Rina Tan',
      email: email.trim().toLowerCase(),
    );
    await _persist(session);
    return session;
  }

  @override
  Future<AuthSession> register({
    required String fullName,
    required String email,
    required String password,
  }) async {
    final session = _buildSession(
      fullName: fullName.trim(),
      email: email.trim().toLowerCase(),
    );
    await _persist(session);
    return session;
  }

  @override
  Future<AuthSession> refresh({required String refreshToken}) async {
    if (!_issuedRefreshTokens.remove(refreshToken)) {
      throw const AuthException('Invalid or expired refresh token');
    }
    final session = _buildSession(
      fullName: 'Rina Tan',
      email: FakeAuthRepository.demoEmail,
    );
    await _persist(session);
    return session;
  }

  @override
  Future<void> logout({required String refreshToken}) async {
    _issuedRefreshTokens.remove(refreshToken);
    await _tokenStore?.clear();
  }

  @override
  Future<AuthSession?> restore() async {
    final store = _tokenStore;
    if (store == null) {
      return null;
    }
    final session = await store.readSession();
    // A fresh process (after a full restart) doesn't remember the refresh
    // tokens it issued in memory, but the persisted session's refresh token is
    // still the live one. Register it so `refresh()` succeeds after a restart
    // instead of treating the session as dead and logging the user out.
    if (session != null) {
      _issuedRefreshTokens.add(session.refreshToken);
    }
    return session;
  }

  AuthSession _buildSession({required String fullName, required String email}) {
    _counter += 1;
    final now = DateTime.now().toUtc();
    final iat = now.millisecondsSinceEpoch ~/ 1000;
    final exp = iat + accessTtlSeconds;
    final id = 'patient-${_counter.toString().padLeft(6, '0')}';

    final session = AuthSession(
      tokenType: 'Bearer',
      accessToken: _mintJwt(iat: iat, exp: exp, sub: id),
      refreshToken: _mintRefreshToken(),
      expiresInSeconds: accessTtlSeconds,
      refreshExpiresAt: now.add(const Duration(days: 30)),
      profile: PatientProfile(
        id: id,
        fullName: fullName,
        email: email,
        preferredCurrency: 'SGD',
        synthetic: true,
        createdAt: now,
        updatedAt: now,
      ),
    );
    _issuedRefreshTokens.add(session.refreshToken);
    return session;
  }

  /// Minimal HS256-shaped JWT: header + base64url payload with `exp` and a
  /// unique `jti` claim. The signature segment is a fixed placeholder (never
  /// verified client-side). The `jti` is a microsecond timestamp so rotated
  /// access tokens are observably different even within the same second or
  /// after a process restart (the in-memory counter resets).
  String _mintJwt({required int iat, required int exp, required String sub}) {
    final header = _b64('{"alg":"HS256","typ":"JWT"}');
    final payload = _b64(
      jsonEncode({
        'iss': 'batam-medhub',
        'aud': 'batam-medhub-mobile',
        'sub': sub,
        'iat': iat,
        'exp': exp,
        'jti': DateTime.now().microsecondsSinceEpoch,
      }),
    );
    return '$header.$payload.synthetic-signature';
  }

  String _mintRefreshToken() =>
      'refresh_fake_${_counter}_${DateTime.now().microsecondsSinceEpoch}';

  String _b64(String value) =>
      base64Url.encode(utf8.encode(value)).replaceAll('=', '');

  Future<void> _persist(AuthSession session) async {
    await _tokenStore?.writeSession(session);
  }
}
