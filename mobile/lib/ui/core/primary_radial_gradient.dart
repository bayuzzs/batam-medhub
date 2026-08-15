import 'package:flutter/material.dart';

import 'package:mobile/ui/core/app_colors.dart';

/// Decorative radial gradient, anchored to the top of the screen by default.
///
/// Draws a soft radial glow of [AppColors.primary] that fades to transparent.
/// Place it as the first (bottom-most) child of a [Stack] so it sits behind
/// screen content — e.g. `OnboardingPage`.
class PrimaryRadialGradient extends StatelessWidget {
  const PrimaryRadialGradient({
    super.key,
    this.alignment = Alignment.topCenter,
    this.radius = 1.0,
    this.startOpacity = 0.60,
  });

  /// Where the gradient originates. Defaults to the top-center.
  final Alignment alignment;

  /// Gradient radius relative to the shortest side of the box.
  final double radius;

  /// Opacity of the primary color at the gradient's origin.
  final double startOpacity;

  @override
  Widget build(BuildContext context) {
    return SizedBox.expand(
      child: DecoratedBox(
        decoration: BoxDecoration(
          gradient: RadialGradient(
            center: alignment,
            radius: radius,
            colors: [
              AppColors.primary.withValues(alpha: startOpacity),
              AppColors.primary.withValues(alpha: 0),
            ],
          ),
        ),
      ),
    );
  }
}
