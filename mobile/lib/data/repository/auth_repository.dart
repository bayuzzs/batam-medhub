import 'package:mobile/models/auth_session.dart';

/// Thrown when an authentication operation fails (e.g. invalid credentials,
/// unknown/expired refresh token, network failure).
class AuthException implements Exception {
  const AuthException(this.message);

  final String message;

  @override
  String toString() => 'AuthException: $message';
}

/// Contract for authentication, agnostic of the transport.
///
/// The real implementation ([AuthRepositoryImpl]) talks to the core API over
/// [AuthApi]; [FakeAuthRepository] is the in-memory stand-in used while the
/// backend isn't implemented. Both implement this interface so they can be
/// swapped via dependency injection (see `application/auth/providers.dart`).
abstract class AuthRepository {
  Future<AuthSession> login({required String email, required String password});

  Future<AuthSession> register({
    required String fullName,
    required String email,
    required String password,
  });

  Future<AuthSession> refresh({required String refreshToken});

  Future<void> logout({required String refreshToken});

  /// Returns the persisted session, or `null` when none is stored.
  Future<AuthSession?> restore();
}
