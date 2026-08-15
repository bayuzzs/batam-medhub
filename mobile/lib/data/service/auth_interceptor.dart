import 'package:dio/dio.dart';

import 'token_store.dart';

/// Attaches the Bearer access token to requests and transparently recovers
/// from a single `401` by refreshing the session and retrying once.
///
/// [refresh] must refresh AND persist the session (this is
/// `AuthController.refresh()`); the interceptor then re-reads the fresh token
/// from [tokenStore] and retries the original request.
///
/// The `_refreshing` guard prevents retry loops — e.g. a `401` from the
/// refresh call itself is surfaced as-is rather than triggering another
/// refresh. Concurrent dedup of the refresh call itself lives in the
/// controller (single in-flight refresh future).
class AuthInterceptor extends Interceptor {
  AuthInterceptor({
    required this._dio,
    required this._tokenStore,
    required this._refresh,
  });

  final Dio _dio;
  final TokenStore _tokenStore;
  final Future<void> Function() _refresh;

  bool _refreshing = false;

  static bool _isAuthPath(String path) => path.contains('/v1/auth/');

  @override
  Future<void> onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    // Auth endpoints are unauthenticated; don't attach a (possibly stale)
    // access token to them.
    if (!_isAuthPath(options.path)) {
      final token = await _tokenStore.readAccessToken();
      if (token != null) {
        options.headers['Authorization'] = 'Bearer $token';
      }
    }
    handler.next(options);
  }

  @override
  Future<void> onError(
    DioException err,
    ErrorInterceptorHandler handler,
  ) async {
    final status = err.response?.statusCode;
    final request = err.requestOptions;
    if (status == 401 && !_refreshing && !_isAuthPath(request.path)) {
      _refreshing = true;
      try {
        await _refresh();
        final token = await _tokenStore.readAccessToken();
        if (token == null) {
          // Refresh failed and the session was signed out; surface the
          // original 401 instead of retrying without a token.
          return handler.reject(err);
        }
        request.headers['Authorization'] = 'Bearer $token';
        final response = await _dio.fetch<dynamic>(request);
        return handler.resolve(response);
      } on DioException catch (e) {
        return handler.reject(e);
      } catch (_) {
        return handler.reject(err);
      } finally {
        _refreshing = false;
      }
    }
    handler.next(err);
  }
}
