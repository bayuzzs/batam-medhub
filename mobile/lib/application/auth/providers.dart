import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:mobile/data/config/api_config.dart';
import 'package:mobile/data/repository/auth_repository.dart';
import 'package:mobile/data/repository/auth_repository_impl.dart';
import 'package:mobile/data/repository/fake_auth_repository.dart';
import 'package:mobile/data/service/auth_api.dart';
import 'package:mobile/data/service/auth_interceptor.dart';
import 'package:mobile/data/service/dio_auth_api.dart';
import 'package:mobile/data/service/secure_token_store.dart';
import 'package:mobile/data/service/shared_prefs_token_store.dart';
import 'package:mobile/data/service/token_store.dart';

import 'auth_controller.dart';

/// DI switch. `false` (default) uses the real Dio backend
/// ([AuthRepositoryImpl] + [DioAuthApi]) with a secure token store, so the app
/// talks to the live core API. Set to `true` (or override
/// [authRepositoryProvider] / [tokenStoreProvider]) to use the fake
/// [FakeAuthRepository] with a `shared_preferences`-backed token store.
///
/// Overridable at build/run time without editing code:
/// ```sh
/// flutter run --dart-define=USE_FAKE_BACKEND=true
/// ```
const bool kUseFakeBackend = bool.fromEnvironment(
  'USE_FAKE_BACKEND',
  defaultValue: false,
);

/// Token persistence. Fake mode persists via `shared_preferences` so a full
/// reload restores the demo session; real mode uses secure storage. Tests
/// override this provider with [InMemoryTokenStore].
final tokenStoreProvider = Provider<TokenStore>((ref) {
  return kUseFakeBackend ? SharedPreferencesTokenStore() : SecureTokenStore();
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
