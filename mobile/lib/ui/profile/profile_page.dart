import 'package:flutter/material.dart';

import 'package:mobile/ui/core/widgets/app_container.dart';

/// Profile screen (bottom nav destination 2).
///
/// Shown inside [MainShell]; the shell owns the [Scaffold] and bottom nav.
/// Placeholder until user data is available from the data layer.
class ProfilePage extends StatelessWidget {
  const ProfilePage({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return SafeArea(
      child: AppContainer(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const SizedBox(height: 24),
            Text('Profile', style: theme.textTheme.headlineMedium),
            const SizedBox(height: 8),
            Text(
              'Your profile details will appear here.',
              style: theme.textTheme.bodyLarge,
            ),
          ],
        ),
      ),
    );
  }
}
