import 'dart:convert';

import 'package:flutter_secure_storage/flutter_secure_storage.dart';

import 'package:mobile/models/auth_session.dart';

import 'token_store.dart';

/// [TokenStore] backed by `flutter_secure_storage` (Keychain / Keystore).
///
/// Stores the serialized [AuthSession] under a single key.
class SecureTokenStore implements TokenStore {
  SecureTokenStore({FlutterSecureStorage? storage})
      : _storage = storage ?? const FlutterSecureStorage();

  final FlutterSecureStorage _storage;

  static const String _sessionKey = 'auth_session';

  @override
  Future<AuthSession?> readSession() async {
    final raw = await _storage.read(key: _sessionKey);
    if (raw == null) {
      return null;
    }
    try {
      final json = jsonDecode(raw);
      if (json is! Map<String, dynamic>) {
        return null;
      }
      return AuthSession.fromJson(json);
    } on FormatException {
      return null;
    }
  }

  @override
  Future<String?> readAccessToken() async {
    final session = await readSession();
    return session?.accessToken;
  }

  @override
  Future<void> writeSession(AuthSession session) async {
    await _storage.write(
      key: _sessionKey,
      value: jsonEncode(session.toJson()),
    );
  }

  @override
  Future<void> clear() async {
    await _storage.delete(key: _sessionKey);
  }
}
