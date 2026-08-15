import 'package:flutter/material.dart';
import 'package:lucide_icons_flutter/lucide_icons.dart';

import 'package:mobile/models/plan_option.dart';
import 'package:mobile/ui/core/theme/app_colors.dart';
import 'package:mobile/ui/core/utils/money_formatter.dart';
import 'package:mobile/ui/core/utils/provider_label.dart';
import 'package:mobile/ui/core/utils/time_window_format.dart';

/// A card describing a [PlanOption] returned by the planner.
///
/// Renders real API data (no fabricated photo/rating/distance — those aren't
/// in the spec): the hospital service, provider, appointment time, total
/// price, and a short explanation, plus a "View Details" action. Booking
/// happens only on the plan detail screen (matching the previous "View
/// Itinerary Details" flow). Styling matches `ItenaryOptionCard` (white card,
/// primary-tinted border).
class PlanOptionCard extends StatelessWidget {
  const PlanOptionCard({
    super.key,
    required this.option,
    required this.onViewDetails,
  });

  final PlanOption option;

  /// Called when "View Details" is tapped (opens the itemized plan, where the
  /// journey is booked).
  final VoidCallback onViewDetails;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final muted = AppColors.text.withValues(alpha: 0.75);

    final hospital = _hospitalItem(option);
    final serviceTitle = hospital?.title ?? option.items.first.title;
    final window = hospital?.timeWindow ?? option.items.first.timeWindow;
    final total = option.totalPrice.displayTotal;

    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: AppColors.background,
        borderRadius: BorderRadius.circular(24),
        border: Border.all(
          color: AppColors.primary.withValues(alpha: 0.35),
          width: 1.5,
        ),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          if (option.rank == 1) ...[
            const _RecommendedBanner(),
            const SizedBox(height: 24),
          ],
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Container(
                width: 44,
                height: 44,
                decoration: BoxDecoration(
                  color: AppColors.primary.withValues(alpha: 0.10),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: const Icon(
                  LucideIcons.hospital,
                  color: AppColors.primary,
                  size: 22,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(serviceTitle, style: theme.textTheme.titleMedium),
                    const SizedBox(height: 4),
                    Text(
                      formatProviderLabel(hospital?.providerId),
                      style: theme.textTheme.bodyMedium?.copyWith(color: muted),
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          _InfoRow(icon: LucideIcons.clock, text: formatWindowWithZone(window)),
          const SizedBox(height: 12),
          Text(
            MoneyFormatter.format(total),
            style: theme.textTheme.titleMedium?.copyWith(
              color: AppColors.heading,
              fontWeight: FontWeight.w600,
            ),
          ),
          if (option.explanation.isNotEmpty) ...[
            const SizedBox(height: 8),
            Text(
              option.explanation.join(' '),
              maxLines: 3,
              overflow: TextOverflow.ellipsis,
              style: theme.textTheme.bodySmall?.copyWith(color: muted),
            ),
          ],
          const SizedBox(height: 12),
          FilledButton(
            onPressed: onViewDetails,
            child: const Text('View Details'),
          ),
        ],
      ),
    );
  }

  /// The hospital-appointment leg, or the first item when none is bookable.
  static PlanItem? _hospitalItem(PlanOption option) {
    for (final item in option.items) {
      if (item.itemType == ItemType.hospitalAppointment) return item;
    }
    return option.items.isEmpty ? null : option.items.first;
  }
}

/// "Recommended" banner for the top-ranked option (matches
/// `ItenaryOptionCard`'s banner).
class _RecommendedBanner extends StatelessWidget {
  const _RecommendedBanner();

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Container(
      height: 32,
      decoration: BoxDecoration(
        color: AppColors.primary.withValues(alpha: 0.12),
        borderRadius: BorderRadius.circular(24),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Icon(LucideIcons.star, size: 20, color: AppColors.heading),
          const SizedBox(width: 8),
          Text(
            'Recommended',
            style: theme.textTheme.bodyMedium?.copyWith(
              fontWeight: FontWeight.w600,
              color: AppColors.heading,
            ),
          ),
        ],
      ),
    );
  }
}

/// A single icon + text information row.
class _InfoRow extends StatelessWidget {
  const _InfoRow({required this.icon, required this.text});

  final IconData icon;
  final String text;

  @override
  Widget build(BuildContext context) {
    final muted = AppColors.text.withValues(alpha: 0.75);

    return Row(
      children: [
        Icon(icon, size: 16, color: AppColors.primary),
        const SizedBox(width: 8),
        Expanded(
          child: Text(
            text,
            style: Theme.of(
              context,
            ).textTheme.bodyMedium?.copyWith(color: muted),
          ),
        ),
      ],
    );
  }
}
