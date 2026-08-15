import 'package:mobile/models/auth_session.dart';

/// Persists the authenticated [AuthSession] (tokens + profile).
///
/// Abstract so the real implementation can use platform secure storage while
/// fake mode and widget tests use an in-memory implementation — no platform
/// channels needed outside real-auth flows.
abstract class TokenStore {
  /// Returns the persisted session, or `null` when none is stored.
  Future<AuthSession?> readSession();

  /// Returns just the current access token (used by the auth interceptor),
  /// or `null` when no session is stored.
  Future<String?> readAccessToken();

  /// Persists [session], replacing any previous one.
  Future<void> writeSession(AuthSession session);

  /// Removes any persisted session.
  Future<void> clear();
}
