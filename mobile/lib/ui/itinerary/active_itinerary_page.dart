import 'package:flutter/material.dart';
import 'package:lucide_icons_flutter/lucide_icons.dart';

import 'package:mobile/models/journey.dart';
import 'package:mobile/models/plan_option.dart';
import 'package:mobile/models/time_window.dart';
import 'package:mobile/ui/core/navigation/app_router.dart';
import 'package:mobile/ui/core/theme/app_colors.dart';
import 'package:mobile/ui/core/theme/app_spacing.dart';
import 'package:mobile/ui/core/utils/item_type_presentation.dart';
import 'package:mobile/ui/core/utils/money_formatter.dart';
import 'package:mobile/ui/core/utils/provider_label.dart';
import 'package:mobile/ui/core/utils/time_window_format.dart';
import 'package:mobile/ui/core/widgets/app_container.dart';
import 'package:mobile/ui/core/widgets/primary_radial_gradient.dart';

/// Active Itinerary screen — the journey the patient is following.
///
/// Shown immediately after a plan is booked (replacing the plan detail and the
/// chat underneath), so the patient lands on their confirmed journey instead
/// of returning to the chat conversation. Renders the confirmed journey from
/// [JourneyDetail.activeItinerary]: a status banner, the hospital summary,
/// each booked leg as a timeline with its status, and the total price.
class ActiveItineraryPage extends StatelessWidget {
  const ActiveItineraryPage({super.key, required this.detail});

  /// The confirmed journey (with its active itinerary version).
  final JourneyDetail detail;

  /// The hospital-appointment leg, or the first item when none exists.
  static ItineraryItem _hospitalItem(JourneyDetail detail) {
    final items = detail.activeItinerary.items;
    for (final item in items) {
      if (item.itemType == ItemType.hospitalAppointment) return item;
    }
    return items.isEmpty ? _stubItem() : items.first;
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final hospital = _hospitalItem(detail);
    final total = detail.activeItinerary.totalPrice.displayTotal;

    return Scaffold(
      body: Stack(
        children: [
          const PrimaryRadialGradient(startOpacity: 0.20),
          SafeArea(
            child: Column(
              children: [
                // Pinned header: back to the home (chat) + screen title.
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
                        onPressed: () => context.goChat(),
                        style: IconButton.styleFrom(
                          backgroundColor: AppColors.primary.withValues(
                            alpha: 0.12,
                          ),
                          foregroundColor: AppColors.heading,
                        ),
                        icon: const Icon(LucideIcons.chevronLeft),
                      ),
                      const SizedBox(width: 12),
                      Text('My Itinerary', style: theme.textTheme.titleLarge),
                    ],
                  ),
                ),
                Expanded(
                  child: AppContainer(
                    scrollable: true,
                    vertical: 12,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      children: [
                        _StatusBanner(detail: detail),
                        const SizedBox(height: 16),
                        _SummaryCard(
                          serviceTitle: hospital.title,
                          providerId: hospital.providerId,
                          window: hospital.timeWindow,
                          total: MoneyFormatter.format(total),
                        ),
                        const SizedBox(height: 16),
                        _TimelineCard(items: detail.activeItinerary.items),
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

/// Teal banner confirming the journey is active, with the booking reference.
class _StatusBanner extends StatelessWidget {
  const _StatusBanner({required this.detail});

  final JourneyDetail detail;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final active = detail.journey.status == JourneyStatus.active;

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.primary.withValues(alpha: 0.10),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(
          color: AppColors.primary.withValues(alpha: 0.35),
          width: 1.5,
        ),
      ),
      child: Row(
        children: [
          Container(
            width: 44,
            height: 44,
            decoration: BoxDecoration(
              color: AppColors.primary.withValues(alpha: 0.15),
              shape: BoxShape.circle,
            ),
            child: const Icon(
              LucideIcons.check,
              color: AppColors.primary,
              size: 22,
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  active ? 'Journey active' : 'Journey in review',
                  style: theme.textTheme.titleMedium?.copyWith(
                    color: AppColors.heading,
                    fontWeight: FontWeight.w700,
                  ),
                ),
                const SizedBox(height: 2),
                Text(
                  'Booking ref · ${detail.journey.id}',
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: AppColors.text.withValues(alpha: 0.75),
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

/// Shared card decoration: white card with a primary-tinted border and
/// rounded corners (matching `itenary_option_card.dart`).
BoxDecoration _cardDecoration() {
  return BoxDecoration(
    color: AppColors.background,
    borderRadius: BorderRadius.circular(24),
    border: Border.all(
      color: AppColors.primary.withValues(alpha: 0.35),
      width: 1.5,
    ),
  );
}

/// Top card: service, provider, appointment window and total price.
class _SummaryCard extends StatelessWidget {
  const _SummaryCard({
    required this.serviceTitle,
    required this.providerId,
    required this.window,
    required this.total,
  });

  final String serviceTitle;
  final String? providerId;
  final TimeWindow window;
  final String total;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final muted = AppColors.text.withValues(alpha: 0.75);

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: _cardDecoration(),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
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
                      formatProviderLabel(providerId),
                      style: theme.textTheme.bodyMedium?.copyWith(color: muted),
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 14),
          _InfoRow(icon: LucideIcons.clock, text: formatWindowWithZone(window)),
          const SizedBox(height: 10),
          Text(
            total,
            style: theme.textTheme.titleLarge?.copyWith(
              color: AppColors.heading,
              fontWeight: FontWeight.w700,
            ),
          ),
        ],
      ),
    );
  }
}

/// A labeled row (icon + text) used inside cards.
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
            ).textTheme.bodySmall?.copyWith(color: muted),
          ),
        ),
      ],
    );
  }
}

/// Booked journey legs as a vertical timeline, each with its status.
class _TimelineCard extends StatelessWidget {
  const _TimelineCard({required this.items});

  final List<ItineraryItem> items;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final muted = AppColors.text.withValues(alpha: 0.75);

    return Container(
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
      decoration: _cardDecoration(),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Text(
            'Confirmed journey, step by step',
            style: theme.textTheme.titleMedium?.copyWith(
              color: AppColors.heading,
              fontWeight: FontWeight.w700,
            ),
          ),
          const SizedBox(height: 8),
          for (var i = 0; i < items.length; i++) ...[
            _TimelineStep(
              item: items[i],
              isLast: i == items.length - 1,
              time: formatWindow(items[i].timeWindow),
              price: items[i].price == null
                  ? null
                  : MoneyFormatter.format(items[i].price!.display),
              provider: formatProviderLabel(items[i].providerId),
            ),
          ],
          const SizedBox(height: 8),
          Divider(color: muted.withValues(alpha: 0.25)),
          Padding(
            padding: const EdgeInsets.only(bottom: 8),
            child: Text(
              "All times are shown in your device's local time.",
              style: theme.textTheme.bodySmall?.copyWith(color: muted),
            ),
          ),
        ],
      ),
    );
  }
}

/// A single timeline entry: a leading dot + connector, then the leg details
/// and a status chip.
class _TimelineStep extends StatelessWidget {
  const _TimelineStep({
    required this.item,
    required this.isLast,
    required this.time,
    required this.price,
    required this.provider,
  });

  final ItineraryItem item;
  final bool isLast;
  final String time;
  final String? price;
  final String provider;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final muted = AppColors.text.withValues(alpha: 0.75);

    return IntrinsicHeight(
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          // Leading dot + connector line.
          SizedBox(
            width: 24,
            child: Column(
              children: [
                Container(
                  width: 12,
                  height: 12,
                  margin: const EdgeInsets.only(top: 6),
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    color: AppColors.primary.withValues(alpha: 0.35),
                    border: Border.all(color: AppColors.primary, width: 2),
                  ),
                ),
                if (!isLast)
                  Expanded(
                    child: Container(
                      width: 2,
                      color: AppColors.primary.withValues(alpha: 0.20),
                    ),
                  ),
              ],
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Padding(
              padding: const EdgeInsets.only(bottom: 16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Icon(
                        itemTypeIcon(item.itemType),
                        size: 16,
                        color: AppColors.primary,
                      ),
                      const SizedBox(width: 6),
                      Expanded(
                        child: Text(
                          item.title,
                          style: theme.textTheme.bodyMedium?.copyWith(
                            fontWeight: FontWeight.w600,
                            color: AppColors.heading,
                          ),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 2),
                  Text(
                    '${itemTypeLabel(item.itemType)} · $time · $provider',
                    style: theme.textTheme.bodySmall?.copyWith(color: muted),
                  ),
                  const SizedBox(height: 6),
                  _StatusChip(status: item.status),
                  if (price != null) ...[
                    const SizedBox(height: 4),
                    Text(
                      price!,
                      style: theme.textTheme.bodyMedium?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ],
                  if (item.operationalNotes.isNotEmpty) ...[
                    const SizedBox(height: 4),
                    Text(
                      item.operationalNotes.join(' '),
                      style: theme.textTheme.bodySmall?.copyWith(color: muted),
                    ),
                  ],
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}

/// Small pill showing an itinerary item's booking status.
class _StatusChip extends StatelessWidget {
  const _StatusChip({required this.status});

  final ItineraryItemStatus status;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final muted = AppColors.text.withValues(alpha: 0.75);
    final confirmed = status == ItineraryItemStatus.confirmed;

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
      decoration: BoxDecoration(
        color: confirmed
            ? AppColors.primary.withValues(alpha: 0.10)
            : AppColors.text.withValues(alpha: 0.08),
        borderRadius: BorderRadius.circular(999),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            confirmed ? LucideIcons.check : LucideIcons.hourglass,
            size: 12,
            color: confirmed ? AppColors.primary : muted,
          ),
          const SizedBox(width: 4),
          Text(
            confirmed ? 'Confirmed' : 'Pending',
            style: theme.textTheme.labelSmall?.copyWith(
              color: confirmed ? AppColors.primary : muted,
              fontWeight: FontWeight.w600,
            ),
          ),
        ],
      ),
    );
  }
}

/// Fallback when a journey has no itinerary items (shouldn't happen for an
/// active journey, but keeps the screen safe to render).
ItineraryItem _stubItem() {
  return ItineraryItem(
    id: 'item-unknown',
    itemType: ItemType.hospitalAppointment,
    providerId: null,
    title: 'Medical appointment',
    status: ItineraryItemStatus.confirmed,
    timeWindow: TimeWindow(
      startsAt: DateTime.now(),
      endsAt: DateTime.now(),
      startTimeZone: 'UTC',
      endTimeZone: 'UTC',
    ),
    operationalNotes: const [],
    synthetic: true,
    source: 'MOCK',
  );
}
