import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:mobile/data/config/api_config.dart';
import 'package:mobile/data/repository/auth_repository.dart';
import 'package:mobile/data/repository/auth_repository_impl.dart';
import 'package:mobile/data/repository/fake_auth_repository.dart';
import 'package:mobile/data/service/auth_api.dart';
import 'package:mobile/data/service/auth_interceptor.dart';
import 'package:mobile/data/service/dio_auth_api.dart';
import 'package:mobile/data/service/in_memory_token_store.dart';
import 'package:mobile/data/service/secure_token_store.dart';
import 'package:mobile/data/service/token_store.dart';

import 'auth_controller.dart';

/// DI switch: `true` (default) uses the in-memory [FakeAuthRepository] and an
/// in-memory token store — no backend or platform secure storage needed.
/// Set to `false` (or override [authRepositoryProvider] / [tokenStoreProvider])
/// to use the real Dio backend + `flutter_secure_storage`.
const bool kUseFakeBackend = true;

/// Token persistence. Fake mode keeps it in memory so tests and demo runs
/// never touch platform channels; real mode uses secure storage.
final tokenStoreProvider = Provider<TokenStore>((ref) {
  return kUseFakeBackend ? InMemoryTokenStore() : SecureTokenStore();
});

/// [Dio] configured with the API base URL and the [AuthInterceptor].
final dioProvider = Provider<Dio>((ref) {
  final dio = Dio(
    BaseOptions(
      baseUrl: ApiConfig.baseUrl,
      connectTimeout: const Duration(seconds: 15),
      receiveTimeout: const Duration(seconds: 15),
    ),
  );
  dio.interceptors.add(
    AuthInterceptor(
      dio: dio,
      tokenStore: ref.watch(tokenStoreProvider),
      // Lazy read: only invoked on a 401, so no provider cycle at build time.
      refresh: () => ref.read(authControllerProvider.notifier).refresh(),
    ),
  );
  return dio;
});

final authApiProvider = Provider<AuthApi>((ref) {
  return DioAuthApi(ref.watch(dioProvider));
});

/// The auth repository. Swap here (or override this provider in tests) to
/// change between fake and real backends.
final authRepositoryProvider = Provider<AuthRepository>((ref) {
  if (kUseFakeBackend) {
    return FakeAuthRepository(tokenStore: ref.watch(tokenStoreProvider));
  }
  return AuthRepositoryImpl(
    api: ref.watch(authApiProvider),
    tokenStore: ref.watch(tokenStoreProvider),
  );
});

final authControllerProvider = NotifierProvider<AuthController, AuthState>(
  AuthController.new,
);
