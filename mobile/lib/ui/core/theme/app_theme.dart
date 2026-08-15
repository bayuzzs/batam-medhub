import 'package:flutter/material.dart';
import 'package:mobile/ui/core/theme/app_colors.dart';

/// Font families used across the app.
///
/// Registered in `pubspec.yaml` under `flutter: fonts:`.
abstract final class AppFonts {
  /// Heading / title font.
  static const String heading = 'Momo Trust Display';

  /// Regular / body font.
  static const String body = 'Poppins';
}

/// Application themes.
///
/// Exposes both a light and a dark `ThemeData` built on a Material 3
/// [ColorScheme] seeded from [AppColors.primary]. Use these in
/// [MaterialApp.theme] or with
/// [MaterialApp.router].
abstract final class AppTheme {
  static ThemeData get light => _build(Brightness.light);

  static ThemeData _build(Brightness brightness) {
    final scheme = ColorScheme.fromSeed(
      seedColor: AppColors.primary,
      brightness: brightness,
      primary: AppColors.primary,
      onPrimary: AppColors.heading,
    );

    final base = ThemeData(
      useMaterial3: true,
      brightness: brightness,
      colorScheme: scheme,
      scaffoldBackgroundColor: AppColors.background,
      fontFamily: AppFonts.body,
    );

    return base.copyWith(
      textTheme: _appTextTheme(base.textTheme),
      filledButtonTheme: _filledButtonTheme(),
      iconButtonTheme: _iconButtonTheme(),
      inputDecorationTheme: _inputDecorationTheme(scheme),
      chipTheme: _chipTheme(),
    );
  }

  /// Global chip style: fully rounded (stadium) with no border, compact
  /// label/icon sizing, and a tighter gap between the icon and label.
  static ChipThemeData _chipTheme() {
    return ChipThemeData(
      side: BorderSide.none,
      shape: const StadiumBorder(),
      labelStyle: const TextStyle(fontSize: 11, fontWeight: FontWeight.w600),
      iconTheme: const IconThemeData(size: 14),
      labelPadding: const EdgeInsets.symmetric(horizontal: 4),
    );
  }

  /// Global primary button style: a bigger touch target and a moderate
  /// corner radius (deliberately NOT the fully-rounded M3 default).
  static FilledButtonThemeData _filledButtonTheme() {
    return FilledButtonThemeData(
      style: FilledButton.styleFrom(
        minimumSize: const Size(64, 32),
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        textStyle: const TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
      ),
    );
  }

  /// Global icon button style: icons render in white by default.
  static IconButtonThemeData _iconButtonTheme() {
    return IconButtonThemeData(
      style: IconButton.styleFrom(foregroundColor: Colors.white),
    );
  }

  /// Global input style: fields use a soft outlined border in
  /// [AppColors.inputBorder] (primary on focus). Labels are rendered above the
  /// box by [AppTextField] (see `app_text_field.dart`), so no floating label
  /// behavior is set here.
  static InputDecorationTheme _inputDecorationTheme(ColorScheme scheme) {
    const border = OutlineInputBorder(
      borderRadius: BorderRadius.all(Radius.circular(12)),
    );

    return InputDecorationTheme(
      filled: false,
      enabledBorder: border.copyWith(
        borderSide: const BorderSide(color: AppColors.inputBorder),
      ),
      focusedBorder: border.copyWith(
        borderSide: const BorderSide(color: AppColors.primary, width: 2),
      ),
      errorBorder: border.copyWith(borderSide: BorderSide(color: scheme.error)),
      focusedErrorBorder: border.copyWith(
        borderSide: BorderSide(color: scheme.error, width: 2),
      ),
      border: border,
    );
  }

  /// Applies the app font families and the heading/body colors to the text
  /// theme: headings use [AppColors.heading], regular text uses
  /// [AppColors.text].
  static TextTheme _appTextTheme(TextTheme base) {
    const heading = TextStyle(
      fontFamily: AppFonts.heading,
      color: AppColors.heading,
    );
    const body = TextStyle(color: AppColors.text);

    return base.copyWith(
      displayLarge: base.displayLarge?.merge(heading),
      displayMedium: base.displayMedium?.merge(heading),
      displaySmall: base.displaySmall?.merge(heading),
      headlineLarge: base.headlineLarge?.merge(heading),
      headlineMedium: base.headlineMedium?.merge(heading),
      headlineSmall: base.headlineSmall?.merge(heading),
      titleLarge: base.titleLarge?.merge(heading),
      titleMedium: base.titleMedium?.merge(heading),
      titleSmall: base.titleSmall?.merge(heading),
      bodyLarge: base.bodyLarge?.merge(body),
      bodyMedium: base.bodyMedium?.merge(body),
      bodySmall: base.bodySmall?.merge(body),
    );
  }
}
