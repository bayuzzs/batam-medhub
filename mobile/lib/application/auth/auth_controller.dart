import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:mobile/data/repository/auth_repository.dart';
import 'package:mobile/models/auth_session.dart';

import 'providers.dart';

/// Lifecycle of the auth session.
enum AuthStatus {
  /// A persisted session is being restored on startup.
  restoring,

  /// A valid session is active.
  authenticated,

  /// No session (signed out, or restore failed).
  unauthenticated,
}

/// Immutable auth state exposed to the UI.
class AuthState {
  const AuthState({
    this.status = AuthStatus.restoring,
    this.session,
    this.errorMessage,
    this.isSubmitting = false,
  });

  final AuthStatus status;

  /// The active session, when [status] is [AuthStatus.authenticated].
  final AuthSession? session;

  /// Message from the last failed login/register (for the UI to surface).
  final String? errorMessage;

  /// Whether a login/register request is in flight.
  final bool isSubmitting;

  bool get isAuthenticated => status == AuthStatus.authenticated;

  AuthState copyWith({
    AuthStatus? status,
    AuthSession? session,
    bool clearSession = false,
    String? errorMessage,
    bool clearError = false,
    bool? isSubmitting,
  }) {
    return AuthState(
      status: status ?? this.status,
      session: clearSession ? null : (session ?? this.session),
      errorMessage: clearError ? null : (errorMessage ?? this.errorMessage),
      isSubmitting: isSubmitting ?? this.isSubmitting,
    );
  }
}

/// Owns the auth session: restore on startup, login/register, scheduled
/// refresh, and logout.
///
/// The scheduled refresh fires just before the access token expires
/// (derived from the JWT `exp` claim); [refresh] is also called by the Dio
/// `AuthInterceptor` on a `401`. Concurrent callers share one in-flight
/// refresh future, so the single-use refresh token is only ever rotated once
/// at a time.
class AuthController extends Notifier<AuthState> {
  /// How long before expiry the scheduled refresh fires.
  static const Duration refreshBuffer = Duration(seconds: 30);

  Timer? _refreshTimer;
  Future<void>? _inFlightRefresh;

  AuthRepository get _repository => ref.read(authRepositoryProvider);

  @override
  AuthState build() {
    ref.onDispose(() => _refreshTimer?.cancel());
    // Fire-and-forget restore; the router waits on [AuthStatus.restoring].
    _restore();
    return const AuthState();
  }

  Future<void> _restore() async {
    try {
      final session = await _repository.restore();
      // A login/logout may have transitioned state while restore was in
      // flight; don't clobber a newer state.
      if (state.status != AuthStatus.restoring) {
        return;
      }
      if (session == null) {
        state = const AuthState(status: AuthStatus.unauthenticated);
        return;
      }
      final receivedAt = DateTime.now().toUtc();
      if (session.accessExpiresAt(receivedAt: receivedAt).isAfter(receivedAt)) {
        _applySession(session);
      } else {
        // Access token already expired — rotate via the refresh token.
        await refresh();
      }
    } catch (_) {
      state = const AuthState(status: AuthStatus.unauthenticated);
    }
  }

  /// Logs in and persists the session. Returns `true` on success.
  Future<bool> login({required String email, required String password}) async {
    state = state.copyWith(isSubmitting: true, clearError: true);
    try {
      final session = await _repository.login(email: email, password: password);
      _applySession(session);
      return true;
    } on AuthException catch (e) {
      state = state.copyWith(isSubmitting: false, errorMessage: e.message);
      return false;
    } catch (_) {
      state = state.copyWith(
        isSubmitting: false,
        errorMessage: 'Unable to log in. Please try again.',
      );
      return false;
    }
  }

  /// Registers and persists the session. Returns `true` on success.
  Future<bool> register({
    required String fullName,
    required String email,
    required String password,
  }) async {
    state = state.copyWith(isSubmitting: true, clearError: true);
    try {
      final session = await _repository.register(
        fullName: fullName,
        email: email,
        password: password,
      );
      _applySession(session);
      return true;
    } on AuthException catch (e) {
      state = state.copyWith(isSubmitting: false, errorMessage: e.message);
      return false;
    } catch (_) {
      state = state.copyWith(
        isSubmitting: false,
        errorMessage: 'Unable to register. Please try again.',
      );
      return false;
    }
  }

  /// Refreshes the session, deduping concurrent callers. On failure the
  /// session is treated as dead and the user is signed out.
  Future<void> refresh() {
    return _inFlightRefresh ??= _doRefresh().whenComplete(() {
      _inFlightRefresh = null;
    });
  }

  Future<void> _doRefresh() async {
    final current = state.session;
    if (current == null) {
      return;
    }
    try {
      final session = await _repository.refresh(
        refreshToken: current.refreshToken,
      );
      _applySession(session);
    } catch (_) {
      await _signOut(clearSession: true);
    }
  }

  /// Logs out (revokes on the backend when possible) and clears the session.
  Future<void> logout() async {
    _refreshTimer?.cancel();
    _refreshTimer = null;
    final refreshToken = state.session?.refreshToken;
    try {
      await _repository.logout(refreshToken: refreshToken ?? '');
    } catch (_) {
      // Local sign-out still proceeds.
    } finally {
      state = const AuthState(status: AuthStatus.unauthenticated);
    }
  }

  Future<void> _signOut({bool clearSession = false}) async {
    _refreshTimer?.cancel();
    _refreshTimer = null;
    if (clearSession) {
      await _repository.logout(refreshToken: state.session?.refreshToken ?? '');
    }
    state = const AuthState(status: AuthStatus.unauthenticated);
  }

  void _applySession(AuthSession session) {
    state = AuthState(status: AuthStatus.authenticated, session: session);
    _scheduleRefresh(session);
  }

  void _scheduleRefresh(AuthSession session) {
    _refreshTimer?.cancel();
    final receivedAt = DateTime.now().toUtc();
    final expiresAt = session.accessExpiresAt(receivedAt: receivedAt);
    final remaining = expiresAt.difference(receivedAt) - refreshBuffer;
    if (remaining <= Duration.zero) {
      refresh();
      return;
    }
    _refreshTimer = Timer(remaining, refresh);
  }
}
