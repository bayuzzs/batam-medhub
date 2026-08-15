import 'package:mobile/models/auth_session.dart';

/// Contract for the core API's authentication operations
/// (`/v1/auth/register`, `/v1/auth/login`, `/v1/auth/refresh`,
/// `/v1/auth/logout`).
abstract class AuthApi {
  Future<AuthSession> login({required String email, required String password});

  Future<AuthSession> register({
    required String fullName,
    required String email,
    required String password,
  });

  Future<AuthSession> refresh({required String refreshToken});

  Future<void> logout({required String refreshToken});
}
