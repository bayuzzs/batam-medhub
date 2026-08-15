import 'package:flutter/material.dart';
import 'package:lucide_icons_flutter/lucide_icons.dart';

import 'package:mobile/ui/core/theme/app_colors.dart';

/// The app's main bottom navigation bar.
///
/// Three destinations: **History**, **New Itinerary** (the chat screen), and
/// **Profile**. The active destination is controlled by [currentIndex]; taps
/// are forwarded to [onSelect] with the destination's index.
///
/// [AppRoutes.history] = 0, [AppRoutes.chat] = 1, [AppRoutes.profile] = 2.
///
/// A custom floating pill (iOS-style): the container shrink-wraps its circular
/// items and is centered above the bottom inset, so it's never full-width.
/// Only the selected item gets a primary circle background (an
/// [AnimatedPositioned] overlay that slides between items on tab change);
/// unselected items are label-less and have no background.
class AppBottomNav extends StatelessWidget {
  const AppBottomNav({
    super.key,
    required this.currentIndex,
    required this.onSelect,
  });

  /// The index of the currently selected destination.
  final int currentIndex;

  /// Called with the tapped destination index.
  final ValueChanged<int> onSelect;

  static const List<({IconData icon, String label})> _destinations = [
    (icon: LucideIcons.history, label: 'History'),
    (icon: LucideIcons.messageCircle, label: 'New Itinerary'),
    (icon: LucideIcons.user, label: 'Profile'),
  ];

  /// Hit-area extent of each item (square).
  static const double _itemExtent = 52;

  /// Diameter of the selected-circle indicator.
  static const double _circleSize = 44;

  /// Horizontal inset of the circle inside an item cell so both stay centered.
  static const double _circleInset = (_itemExtent - _circleSize) / 2;

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      top: false,
      child: Padding(
        padding: const EdgeInsets.fromLTRB(16, 0, 16, 12),
        // The Scaffold's bottomNavigationBar slot passes a *full-screen-height*
        // constraint (maxHeight = screen height). A plain Center would expand
        // to fill it and cover the whole screen, so heightFactor: 1.0 makes the
        // bar shrink-wrap to its content height while still centering it
        // horizontally.
        child: Align(
          alignment: Alignment.bottomCenter,
          heightFactor: 1.0,
          child: Container(
            key: const Key('app_bottom_nav_pill'),
            padding: const EdgeInsets.all(6),
            decoration: BoxDecoration(
              color: AppColors.background,
              borderRadius: BorderRadius.circular(36),
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withValues(alpha: 0.08),
                  blurRadius: 24,
                  offset: const Offset(0, 8),
                ),
              ],
            ),
            // Transparent Material so each item's ripple draws on top of the
            // pill instead of the Scaffold behind it.
            child: Material(
              color: Colors.transparent,
              child: SizedBox(
                height: _itemExtent,
                child: Stack(
                  children: [
                    // The selected circle, sliding between items whenever the
                    // active tab changes.
                    AnimatedPositioned(
                      key: const Key('app_bottom_nav_indicator'),
                      duration: const Duration(milliseconds: 250),
                      curve: Curves.easeOutCubic,
                      left: currentIndex * _itemExtent + _circleInset,
                      top: (_itemExtent - _circleSize) / 2,
                      child: Container(
                        width: _circleSize,
                        height: _circleSize,
                        decoration: const BoxDecoration(
                          shape: BoxShape.circle,
                          color: AppColors.primary,
                        ),
                      ),
                    ),
                    Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        for (var i = 0; i < _destinations.length; i++)
                          Semantics(
                            label: _destinations[i].label,
                            button: true,
                            child: _NavItem(
                              key: Key('app_bottom_nav_item_$i'),
                              icon: _destinations[i].icon,
                              selected: i == currentIndex,
                              onTap: () => onSelect(i),
                            ),
                          ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}

/// A single icon tab inside [AppBottomNav].
///
/// No background of its own: the sliding circle in the [Stack] marks the
/// selected item, so unselected items just show a muted icon.
class _NavItem extends StatelessWidget {
  const _NavItem({
    super.key,
    required this.icon,
    required this.selected,
    required this.onTap,
  });

  final IconData icon;
  final bool selected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final muted = AppColors.text.withValues(alpha: 0.45);

    return InkWell(
      onTap: onTap,
      customBorder: const CircleBorder(),
      child: SizedBox(
        width: AppBottomNav._itemExtent,
        height: AppBottomNav._itemExtent,
        child: Icon(icon, size: 22, color: selected ? Colors.white : muted),
      ),
    );
  }
}
