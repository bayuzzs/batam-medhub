import 'package:flutter_test/flutter_test.dart';

import 'package:mobile/data/repository/auth_repository.dart';
import 'package:mobile/data/repository/fake_auth_repository.dart';
import 'package:mobile/data/service/in_memory_token_store.dart';

void main() {
  group('FakeAuthRepository', () {
    test('login returns a session and persists it', () async {
      final store = InMemoryTokenStore();
      final repo = FakeAuthRepository(tokenStore: store);

      final session = await repo.login(
        email: 'rina.tan@example.test',
        password: 'Demo-Only-Password-2026!',
      );

      expect(session.profile.fullName, 'Rina Tan');
      expect(session.profile.email, 'rina.tan@example.test');
      expect(session.profile.synthetic, isTrue);
      expect(session.refreshToken, isNotEmpty);

      final stored = await store.readSession();
      expect(stored, isNotNull);
      expect(stored!.refreshToken, session.refreshToken);
    });

    test('refresh rotates the token pair and honors single use', () async {
      final repo = FakeAuthRepository();
      final s1 = await repo.login(
        email: 'rina.tan@example.test',
        password: 'x',
      );

      final s2 = await repo.refresh(refreshToken: s1.refreshToken);
      expect(s2.accessToken, isNot(s1.accessToken));
      expect(s2.refreshToken, isNot(s1.refreshToken));

      // The old refresh token was consumed — reusing it must fail.
      await expectLater(
        repo.refresh(refreshToken: s1.refreshToken),
        throwsA(isA<AuthException>()),
      );
    });

    test('restore reads back the persisted session', () async {
      final store = InMemoryTokenStore();
      final repo = FakeAuthRepository(tokenStore: store);
      await repo.login(email: 'rina.tan@example.test', password: 'x');

      final restored = await repo.restore();
      expect(restored, isNotNull);
    });

    test('logout clears the persisted session', () async {
      final store = InMemoryTokenStore();
      final repo = FakeAuthRepository(tokenStore: store);
      final session = await repo.login(
        email: 'rina.tan@example.test',
        password: 'x',
      );

      await repo.logout(refreshToken: session.refreshToken);

      expect(await store.readSession(), isNull);
      expect(await repo.restore(), isNull);
    });
  });
}
