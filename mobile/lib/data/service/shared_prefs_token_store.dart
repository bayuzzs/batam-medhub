import 'dart:convert';

import 'package:shared_preferences/shared_preferences.dart';

import 'package:mobile/models/auth_session.dart';

import 'token_store.dart';

/// [TokenStore] backed by `shared_preferences`.
///
/// Used by the fake-auth mode so demo sessions survive a full app restart
/// without needing platform secure storage (works on Linux, Android, and web
/// out of the box). Storing the fake session is fine: it only ever holds
/// visibly synthetic hackathon credentials.
class SharedPreferencesTokenStore implements TokenStore {
  static const String _sessionKey = 'auth_session';

  @override
  Future<AuthSession?> readSession() async {
    final prefs = await SharedPreferences.getInstance();
    final raw = prefs.getString(_sessionKey);
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
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_sessionKey, jsonEncode(session.toJson()));
  }

  @override
  Future<void> clear() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_sessionKey);
  }
}
