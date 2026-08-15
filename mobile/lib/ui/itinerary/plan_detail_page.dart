import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:lucide_icons_flutter/lucide_icons.dart';

import 'package:mobile/application/chat/providers.dart';
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

/// Plan / itinerary detail screen.
///
/// Pushed from a plan option card in the chat. Shows the hospital service,
/// the itemized journey legs (ferry, transfer, appointment, return) as a
/// timeline with per-leg times and prices, the planner's explanation, and a
/// pinned "Book this journey" action. Booking confirms the plan and opens the
/// [ActiveItineraryPage] — the patient lands on their confirmed journey
/// instead of returning to the chat.
class PlanDetailPage extends ConsumerStatefulWidget {
  const PlanDetailPage({super.key, required this.option});

  final PlanOption option;

  @override
  ConsumerState<PlanDetailPage> createState() => _PlanDetailPageState();
}

class _PlanDetailPageState extends ConsumerState<PlanDetailPage> {
  /// Whether the booking request is in flight (disables the button).
  bool _booking = false;

  /// Confirms the selected plan, then opens the active itinerary screen.
  Future<void> _book() async {
    if (_booking) return;
    setState(() => _booking = true);
    final notifier = ref.read(chatControllerProvider.notifier);
    await notifier.selectPlanOption(widget.option);
    final journey = ref.read(chatControllerProvider).journey;
    if (!mounted) return;
    if (journey != null) {
      context.goActiveItinerary(journey);
    } else {
      setState(() => _booking = false);
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text("Sorry, that journey couldn't be booked. Try again."),
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final option = widget.option;
    final hospital = _hospitalItem(option);
    final serviceTitle = hospital?.title ?? option.items.first.title;
    final window = hospital?.timeWindow ?? option.items.first.timeWindow;
    final total = option.totalPrice.displayTotal;

    return Scaffold(
      body: Stack(
        children: [
          const PrimaryRadialGradient(startOpacity: 0.20),
          SafeArea(
            child: Column(
              children: [
                // Pinned header: back button + screen title.
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
                      Text('Journey Plan', style: theme.textTheme.titleLarge),
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
                        _SummaryCard(
                          serviceTitle: serviceTitle,
                          providerId: hospital?.providerId,
                          window: window,
                          total: MoneyFormatter.format(total),
                          recommended: option.rank == 1,
                        ),
                        const SizedBox(height: 16),
                        _TimelineCard(items: option.items),
                        if (option.explanation.isNotEmpty) ...[
                          const SizedBox(height: 16),
                          _ExplanationCard(text: option.explanation.join(' ')),
                        ],
                      ],
                    ),
                  ),
                ),
                // Pinned book action: confirms the option, then opens the
                // active itinerary screen (not the chat).
                Padding(
                  padding: const EdgeInsets.fromLTRB(
                    AppSpacing.screen,
                    8,
                    AppSpacing.screen,
                    AppSpacing.screen,
                  ),
                  child: SizedBox(
                    width: double.infinity,
                    child: FilledButton(
                      onPressed: _booking ? null : _book,
                      child: _booking
                          ? const SizedBox(
                              width: 20,
                              height: 20,
                              child: CircularProgressIndicator(
                                strokeWidth: 2,
                                color: AppColors.heading,
                              ),
                            )
                          : const Text('Book this journey'),
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

  /// The hospital-appointment leg, or the first item when none is bookable.
  static PlanItem? _hospitalItem(PlanOption option) {
    for (final item in option.items) {
      if (item.itemType == ItemType.hospitalAppointment) return item;
    }
    return option.items.isEmpty ? null : option.items.first;
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

/// Top card: recommended banner (when rank 1), service, provider, time window
/// and total price.
class _SummaryCard extends StatelessWidget {
  const _SummaryCard({
    required this.serviceTitle,
    required this.providerId,
    required this.window,
    required this.total,
    required this.recommended,
  });

  final String serviceTitle;
  final String? providerId;
  final TimeWindow window;
  final String total;
  final bool recommended;

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
              if (recommended) ...[
                const Icon(
                  LucideIcons.star,
                  size: 20,
                  color: AppColors.heading,
                ),
                const SizedBox(width: 6),
                Text(
                  'Recommended',
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: AppColors.heading,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ],
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

/// Itemized journey legs as a vertical timeline: time, title, provider,
/// price, and any operational notes per leg.
class _TimelineCard extends StatelessWidget {
  const _TimelineCard({required this.items});

  final List<PlanItem> items;

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
            'Your journey, step by step',
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

/// A single timeline entry: a leading dot + connector, then the leg details.
class _TimelineStep extends StatelessWidget {
  const _TimelineStep({
    required this.item,
    required this.isLast,
    required this.time,
    required this.price,
    required this.provider,
  });

  final PlanItem item;
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
                  if (price != null) ...[
                    const SizedBox(height: 2),
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

/// The planner's rationale for the option.
class _ExplanationCard extends StatelessWidget {
  const _ExplanationCard({required this.text});

  final String text;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final muted = AppColors.text.withValues(alpha: 0.75);

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: _cardDecoration(),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Why this plan',
            style: theme.textTheme.titleMedium?.copyWith(
              color: AppColors.heading,
              fontWeight: FontWeight.w700,
            ),
          ),
          const SizedBox(height: 8),
          Text(text, style: theme.textTheme.bodyMedium?.copyWith(color: muted)),
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
