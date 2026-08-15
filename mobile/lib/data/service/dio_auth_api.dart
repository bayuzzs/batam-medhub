import 'package:dio/dio.dart';

import 'package:mobile/models/auth_requests.dart';
import 'package:mobile/models/auth_session.dart';

import 'auth_api.dart';

/// [AuthApi] implemented over [Dio].
///
/// Endpoints are unauthenticated (`security: []` in the OpenAPI spec), so no
/// access token is required; the [AuthInterceptor] skips them.
class DioAuthApi implements AuthApi {
  const DioAuthApi(this._dio);

  final Dio _dio;

  @override
  Future<AuthSession> login({
    required String email,
    required String password,
  }) async {
    final response = await _dio.post<Map<String, dynamic>>(
      '/v1/auth/login',
      data: LoginRequest(email: email, password: password).toJson(),
    );
    return AuthSession.fromJson(response.data!);
  }

  @override
  Future<AuthSession> register({
    required String fullName,
    required String email,
    required String password,
  }) async {
    final response = await _dio.post<Map<String, dynamic>>(
      '/v1/auth/register',
      data: RegisterPatientRequest(
        fullName: fullName,
        email: email,
        password: password,
      ).toJson(),
    );
    return AuthSession.fromJson(response.data!);
  }

  @override
  Future<AuthSession> refresh({required String refreshToken}) async {
    final response = await _dio.post<Map<String, dynamic>>(
      '/v1/auth/refresh',
      data: RefreshTokenRequest(refreshToken: refreshToken).toJson(),
    );
    return AuthSession.fromJson(response.data!);
  }

  @override
  Future<void> logout({required String refreshToken}) async {
    await _dio.post<void>(
      '/v1/auth/logout',
      data: RefreshTokenRequest(refreshToken: refreshToken).toJson(),
    );
  }
}
