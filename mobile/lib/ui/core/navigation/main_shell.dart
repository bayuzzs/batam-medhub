import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'package:mobile/ui/core/navigation/app_bottom_nav.dart';

/// Root shell for the bottom-nav destinations.
///
/// Owns the [Scaffold] and [AppBottomNav]; the active tab's page is shown via
/// a [StatefulNavigationShell]. Each branch keeps its own navigation stack and
/// state across tab switches (chat history, form input, etc. are preserved).
class MainShell extends StatelessWidget {
  const MainShell({super.key, required this.navigationShell});

  final StatefulNavigationShell navigationShell;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: navigationShell,
      bottomNavigationBar: AppBottomNav(
        currentIndex: navigationShell.currentIndex,
        onSelect: (index) => navigationShell.goBranch(
          index,
          // Re-tapping the active tab resets that branch to its start.
          initialLocation: index == navigationShell.currentIndex,
        ),
      ),
    );
  }
}
