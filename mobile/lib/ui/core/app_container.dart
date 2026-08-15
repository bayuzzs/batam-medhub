import 'package:flutter/material.dart';

import 'package:mobile/ui/core/app_spacing.dart';

/// Standard screen container.
///
/// Wraps content in a [SafeArea] with the app's standard horizontal padding
/// ([AppSpacing.screen]) so screens don't repeat the padding themselves.
///
/// - Use [AppContainer] with no extra args for a padded, non-scrolling screen
///   (e.g. the onboarding screen).
/// - Set [scrollable] to `true` for content that may overflow the viewport
///   (e.g. forms), and [vertical] to add vertical breathing room.
class AppContainer extends StatelessWidget {
  const AppContainer({
    super.key,
    this.horizontal = AppSpacing.screen,
    this.vertical = 0,
    this.scrollable = false,
    this.child,
  });

  /// Horizontal padding for the content edges.
  final double horizontal;

  /// Vertical padding for the content edges.
  final double vertical;

  /// Whether the content is wrapped in a [SingleChildScrollView].
  final bool scrollable;

  /// The content to display.
  final Widget? child;

  @override
  Widget build(BuildContext context) {
    final content = Padding(
      padding: EdgeInsets.symmetric(horizontal: horizontal, vertical: vertical),
      child: child,
    );

    return SafeArea(
      child: scrollable ? SingleChildScrollView(child: content) : content,
    );
  }
}
