import 'package:flutter/foundation.dart';

/// Runtime API configuration for the Batam MedHub core backend.
abstract final class ApiConfig {
  /// Base URL for the core API.
  ///
  /// Override at build/run time with:
  /// ```sh
  /// flutter run --dart-define=API_BASE_URL=http://host:port
  /// ```
  ///
  /// When no override is supplied the default is platform-aware:
  /// - Android: `http://10.0.2.2:8080` — the Android emulator's alias for the
  ///   host machine's `localhost`, so the app reaches a backend running on the
  ///   host. On a physical Android device `10.0.2.2` does not work; override
  ///   `API_BASE_URL` with your host's LAN IP instead.
  /// - All other platforms (web, desktop): `http://localhost:8080`.
  static String get baseUrl {
    const override = String.fromEnvironment('API_BASE_URL');
    if (override.isNotEmpty) return override;
    if (!kIsWeb && defaultTargetPlatform == TargetPlatform.android) {
      return 'http://10.0.2.2:8080';
    }
    return 'http://localhost:8080';
  }
}
