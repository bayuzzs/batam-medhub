import 'package:mobile/models/auth_session.dart';

import 'token_store.dart';

/// [TokenStore] that keeps the session in memory only.
///
/// Used by the fake-auth mode (no platform secure storage needed) and by
/// widget/unit tests. Nothing is persisted across app restarts.
class InMemoryTokenStore implements TokenStore {
  AuthSession? _session;

  @override
  Future<AuthSession?> readSession() async => _session;

  @override
  Future<String?> readAccessToken() async => _session?.accessToken;

  @override
  Future<void> writeSession(AuthSession session) async {
    _session = session;
  }

  @override
  Future<void> clear() async {
    _session = null;
  }
}
