/// Spacing / padding tokens used across the app.
///
/// Centralizes the app's padding values so screens don't repeat raw numbers.
/// Use [AppContainer] (`app_container.dart`) to apply the standard screen
/// padding without hand-rolling `Padding` widgets.
abstract final class AppSpacing {
  /// Standard horizontal padding applied to the edges of screen content.
  static const double screen = 24;
}
