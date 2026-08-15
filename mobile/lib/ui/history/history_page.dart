import 'package:flutter/material.dart';

import 'package:mobile/ui/core/app_container.dart';

/// History screen (bottom nav destination 0).
///
/// Shown inside [MainShell]; the shell owns the [Scaffold] and bottom nav.
/// Placeholder until past journeys are fetched from the data layer.
class HistoryPage extends StatelessWidget {
  const HistoryPage({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return SafeArea(
      child: AppContainer(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const SizedBox(height: 24),
            Text('History', style: theme.textTheme.headlineMedium),
            const SizedBox(height: 8),
            Text(
              'Your past medical journeys will appear here.',
              style: theme.textTheme.bodyLarge,
            ),
          ],
        ),
      ),
    );
  }
}
