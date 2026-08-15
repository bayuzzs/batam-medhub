import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:mobile/application/auth/auth_controller.dart';
import 'package:mobile/application/auth/providers.dart';
import 'package:mobile/data/repository/auth_repository.dart';
import 'package:mobile/data/repository/fake_auth_repository.dart';
import 'package:mobile/data/service/in_memory_token_store.dart';
import 'package:mobile/data/service/token_store.dart';
import 'package:mobile/models/auth_session.dart';

/// Builds a [ProviderContainer] with the fake backend wired up.
///
/// The same [TokenStore] instance is shared between the store override and
/// the fake repository so restore/login persist across both.
ProviderContainer makeContainer({
  TokenStore? tokenStore,
  AuthRepository? repository,
}) {
  final store = tokenStore ?? InMemoryTokenStore();
  final repo = repository ?? FakeAuthRepository(tokenStore: store);
  return ProviderContainer(
    overrides: [
      tokenStoreProvider.overrideWithValue(store),
      authRepositoryProvider.overrideWithValue(repo),
    ],
  );
}

void main() {
  test('login authenticates, persists, and clears submitting', () async {
    final container = makeContainer();
    addTearDown(container.dispose);
    final controller = container.read(authControllerProvider.notifier);
    await pumpEventQueue();

    final ok = await controller.login(
      email: 'rina.tan@example.test',
      password: 'Demo-Only-Password-2026!',
    );
    expect(ok, isTrue);

    final state = container.read(authControllerProvider);
    expect(state.status, AuthStatus.authenticated);
    expect(state.session, isNotNull);
    expect(state.session!.profile.email, 'rina.tan@example.test');
    expect(state.isSubmitting, isFalse);

    final stored = await container.read(tokenStoreProvider).readSession();
    expect(stored, isNotNull);
  });

  test('register authenticates with the provided name', () async {
    final container = makeContainer();
    addTearDown(container.dispose);
    final controller = container.read(authControllerProvider.notifier);
    await pumpEventQueue();

    final ok = await controller.register(
      fullName: 'Budi Santoso',
      email: 'budi@example.test',
      password: 'secret123',
    );
    expect(ok, isTrue);

    final state = container.read(authControllerProvider);
    expect(state.status, AuthStatus.authenticated);
    expect(state.session!.profile.fullName, 'Budi Santoso');
  });

  test('restore restores a persisted session on a fresh container', () async {
    final store = InMemoryTokenStore();
    final c1 = makeContainer(tokenStore: store);
    await c1
        .read(authControllerProvider.notifier)
        .login(email: 'rina.tan@example.test', password: 'x');
    c1.dispose();

    final c2 = makeContainer(tokenStore: store);
    addTearDown(c2.dispose);

    // Reading the controller triggers build(), which fires the async restore.
    c2.read(authControllerProvider.notifier);
    await pumpEventQueue();
    expect(c2.read(authControllerProvider).status, AuthStatus.authenticated);
  });

  test('restore with no stored session stays unauthenticated', () async {
    final container = makeContainer();
    addTearDown(container.dispose);
    container.read(authControllerProvider.notifier);
    await pumpEventQueue();

    expect(
      container.read(authControllerProvider).status,
      AuthStatus.unauthenticated,
    );
  });

  test('refresh rotates the session while staying authenticated', () async {
    final container = makeContainer();
    addTearDown(container.dispose);
    final controller = container.read(authControllerProvider.notifier);
    await pumpEventQueue();
    await controller.login(email: 'rina.tan@example.test', password: 'x');
    final before = container.read(authControllerProvider).session!;

    await controller.refresh();

    final after = container.read(authControllerProvider).session!;
    expect(after.accessToken, isNot(before.accessToken));
    expect(after.refreshToken, isNot(before.refreshToken));
    expect(
      container.read(authControllerProvider).status,
      AuthStatus.authenticated,
    );
  });

  test('concurrent refresh calls share one rotation', () async {
    final container = makeContainer();
    addTearDown(container.dispose);
    final controller = container.read(authControllerProvider.notifier);
    await pumpEventQueue();
    await controller.login(email: 'rina.tan@example.test', password: 'x');

    // Fire two refreshes without awaiting between them. Thanks to the
    // in-flight dedup, the single-use refresh token is only rotated once.
    final f1 = controller.refresh();
    final f2 = controller.refresh();
    await Future.wait([f1, f2]);

    expect(
      container.read(authControllerProvider).status,
      AuthStatus.authenticated,
    );
  });

  test('refresh failure signs out', () async {
    final store = InMemoryTokenStore();
    final repo = FakeAuthRepository(tokenStore: store);
    final container = makeContainer(tokenStore: store, repository: repo);
    addTearDown(container.dispose);
    final controller = container.read(authControllerProvider.notifier);
    await pumpEventQueue();
    await controller.login(email: 'rina.tan@example.test', password: 'x');

    // Consume the single-use refresh token behind the controller's back.
    final token = container.read(authControllerProvider).session!.refreshToken;
    await repo.logout(refreshToken: token);

    await controller.refresh();

    final state = container.read(authControllerProvider);
    expect(state.status, AuthStatus.unauthenticated);
    expect(state.session, isNull);
  });

  test('logout signs out and clears the persisted session', () async {
    final container = makeContainer();
    addTearDown(container.dispose);
    final controller = container.read(authControllerProvider.notifier);
    await pumpEventQueue();
    await controller.login(email: 'rina.tan@example.test', password: 'x');

    await controller.logout();

    expect(
      container.read(authControllerProvider).status,
      AuthStatus.unauthenticated,
    );
    expect(await container.read(tokenStoreProvider).readSession(), isNull);
  });

  test('login failure surfaces an error message and resets submitting',
      () async {
    // A repository that always rejects login.
    final failing = _FailingRepository();
    final container = makeContainer(repository: failing);
    addTearDown(container.dispose);
    final controller = container.read(authControllerProvider.notifier);
    await pumpEventQueue();

    final ok = await controller.login(email: 'a@b.c', password: 'x');
    expect(ok, isFalse);

    final state = container.read(authControllerProvider);
    expect(state.isSubmitting, isFalse);
    expect(state.errorMessage, isNotNull);
    expect(state.status, isNot(AuthStatus.authenticated));
  });
}

class _FailingRepository implements AuthRepository {
  @override
  Future<AuthSession> login({
    required String email,
    required String password,
  }) {
    throw const AuthException('Invalid credentials');
  }

  @override
  Future<AuthSession> register({
    required String fullName,
    required String email,
    required String password,
  }) {
    throw const AuthException('Invalid credentials');
  }

  @override
  Future<void> logout({required String refreshToken}) async {}

  @override
  Future<AuthSession> refresh({required String refreshToken}) {
    throw const AuthException('Invalid or expired refresh token');
  }

  @override
  Future<AuthSession?> restore() async => null;
}
