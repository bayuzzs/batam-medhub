/// Runtime API configuration for the Batam MedHub core backend.
abstract final class ApiConfig {
  /// Base URL for the core API.
  ///
  /// Override at build/run time with:
  /// ```sh
  /// flutter run --dart-define=API_BASE_URL=http://host:port
  /// ```
  static const String baseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://localhost:8080',
  );
}
