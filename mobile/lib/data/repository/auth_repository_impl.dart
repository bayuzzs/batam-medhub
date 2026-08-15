import 'package:mobile/models/auth_session.dart';

import '../service/auth_api.dart';
import '../service/token_store.dart';
import 'auth_repository.dart';

/// Real [AuthRepository]: calls the core API over [AuthApi] and persists the
/// session via [TokenStore].
class AuthRepositoryImpl implements AuthRepository {
  const AuthRepositoryImpl({required this._api, required this._tokenStore});

  final AuthApi _api;
  final TokenStore _tokenStore;

  @override
  Future<AuthSession> login({
    required String email,
    required String password,
  }) async {
    final session = await _api.login(email: email, password: password);
    await _tokenStore.writeSession(session);
    return session;
  }

  @override
  Future<AuthSession> register({
    required String fullName,
    required String email,
    required String password,
  }) async {
    final session = await _api.register(
      fullName: fullName,
      email: email,
      password: password,
    );
    await _tokenStore.writeSession(session);
    return session;
  }

  @override
  Future<AuthSession> refresh({required String refreshToken}) async {
    final session = await _api.refresh(refreshToken: refreshToken);
    await _tokenStore.writeSession(session);
    return session;
  }

  @override
  Future<void> logout({required String refreshToken}) async {
    try {
      await _api.logout(refreshToken: refreshToken);
    } finally {
      // Always clear the local session, even if the API call fails.
      await _tokenStore.clear();
    }
  }

  @override
  Future<AuthSession?> restore() => _tokenStore.readSession();
}
