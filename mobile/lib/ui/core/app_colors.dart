import 'dart:ui';

/// Central palette for the app.
///
/// Batam MedHub uses a medical teal as the primary brand color with a
/// complementary accent. Keep color choices here so they stay consistent
/// across light and dark themes.
abstract final class AppColors {
  // Brand / seed colors.
  static const Color primary = Color(0xFF28D1D0);
  static const Color heading = Color(0xFF022857);
  static const Color text = Color(0xFF1E1E1E);
  static const Color background = Color(0xFFFFFFFF);

  /// Input field border
  static const Color inputBorder = Color(0x8854F3F2);
}
