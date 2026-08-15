import 'package:flutter_test/flutter_test.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'package:mobile/data/service/shared_prefs_token_store.dart';
import 'package:mobile/models/auth_session.dart';
import 'package:mobile/models/patient_profile.dart';

void main() {
  group('SharedPreferencesTokenStore', () {
    setUp(() {
      SharedPreferences.setMockInitialValues({});
    });

    test('round-trips a session across store instances', () async {
      final store = SharedPreferencesTokenStore();
      await store.writeSession(_sampleSession());

      // A new instance simulates a fresh app process reading the same storage.
      final restored = await SharedPreferencesTokenStore().readSession();
      expect(restored, isNotNull);
      expect(restored!.accessToken, 'access-token');
      expect(restored.refreshToken, 'refresh-token');
      expect(restored.profile.fullName, 'Rina Tan');
    });

    test('readAccessToken returns the persisted token', () async {
      final store = SharedPreferencesTokenStore();
      await store.writeSession(_sampleSession());

      expect(await store.readAccessToken(), 'access-token');
    });

    test('clear removes the persisted session', () async {
      final store = SharedPreferencesTokenStore();
      await store.writeSession(_sampleSession());

      await store.clear();

      expect(await store.readSession(), isNull);
    });
  });
}

AuthSession _sampleSession() {
  final now = DateTime.now().toUtc();
  return AuthSession(
    tokenType: 'Bearer',
    accessToken: 'access-token',
    refreshToken: 'refresh-token',
    expiresInSeconds: 900,
    refreshExpiresAt: now.add(const Duration(days: 7)),
    profile: PatientProfile(
      id: 'pt_1',
      fullName: 'Rina Tan',
      email: 'rina.tan@example.test',
      preferredCurrency: 'IDR',
      synthetic: true,
      createdAt: now,
      updatedAt: now,
    ),
  );
}
