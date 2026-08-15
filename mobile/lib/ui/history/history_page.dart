import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:lucide_icons_flutter/lucide_icons.dart';

import 'package:mobile/ui/core/theme/app_colors.dart';
import 'package:mobile/ui/core/theme/app_spacing.dart';
import 'package:mobile/ui/core/widgets/app_container.dart';
import 'package:mobile/ui/core/widgets/primary_radial_gradient.dart';

/// History screen, pushed full-screen from the chat top bar.
///
/// Placeholder until past journeys are fetched from the data layer. Back
/// returns to the chat screen.
class HistoryPage extends StatelessWidget {
  const HistoryPage({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      body: Stack(
        children: [
          const PrimaryRadialGradient(startOpacity: 0.20),
          SafeArea(
            child: Column(
              children: [
                // Header with back button and screen title.
                Padding(
                  padding: const EdgeInsets.fromLTRB(
                    AppSpacing.screen,
                    8,
                    AppSpacing.screen,
                    8,
                  ),
                  child: Row(
                    children: [
                      IconButton(
                        key: const Key('page_back_button'),
                        onPressed: () => context.pop(),
                        style: IconButton.styleFrom(
                          backgroundColor: AppColors.primary.withValues(
                            alpha: 0.12,
                          ),
                          foregroundColor: AppColors.heading,
                        ),
                        icon: const Icon(LucideIcons.chevronLeft),
                      ),
                      const SizedBox(width: 12),
                      Text('History', style: theme.textTheme.titleLarge),
                    ],
                  ),
                ),
                Expanded(
                  child: AppContainer(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const SizedBox(height: 24),
                        Text(
                          'Your past medical journeys will appear here.',
                          style: theme.textTheme.bodyLarge,
                        ),
                      ],
                    ),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
