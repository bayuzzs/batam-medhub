import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:lucide_icons_flutter/lucide_icons.dart';

import 'package:mobile/ui/core/theme/app_colors.dart';
import 'package:mobile/ui/core/theme/app_spacing.dart';
import 'package:mobile/ui/core/widgets/app_container.dart';
import 'package:mobile/ui/core/widgets/primary_radial_gradient.dart';

/// A single entry in the journey timeline: a time, the activity, and where it
/// happens.
class _JourneyStep {
  const _JourneyStep({
    required this.time,
    required this.activity,
    required this.location,
  });

  final String time;
  final String activity;
  final String location;
}

/// Itinerary Journey Detail screen.
///
/// Pushed from the chat screen's itinerary option card ("View Itinerary
/// Details"). Shows the hospital, a day-by-day journey timeline, an itinerary
/// summary, and a "Choose This Itinerary" action. Content is synthetic demo
/// data until the core API exposes itinerary details.
class ItineraryJourneyDetailPage extends StatelessWidget {
  const ItineraryJourneyDetailPage({
    super.key,
    required this.imageUrl,
    required this.providerName,
    required this.serviceName,
    required this.location,
    required this.appointment,
    required this.duration,
    required this.price,
  });

  /// Hospital / provider photo (local `assets/...` path or remote URL).
  final String imageUrl;

  /// Hospital name.
  final String providerName;

  /// The offered medical service.
  final String serviceName;

  /// Hospital location (e.g. "Batu Aji · 5 km").
  final String location;

  /// Appointment availability (e.g. "Tomorrow, 09:00").
  final String appointment;

  /// Total trip duration (e.g. "3 days").
  final String duration;

  /// Estimated price per person (e.g. "IDR 4.5jt").
  final String price;

  /// Header label for the timeline card.
  static const String _journeyDate = 'Day 1 · Tomorrow';

  /// What the estimated cost includes, shown as chips in the summary.
  static const List<String> _includedItems = [
    'Ferry',
    'Transport',
    'Medical',
    'Hotel',
  ];

  List<_JourneyStep> _steps() {
    return [
      const _JourneyStep(
        time: '06:30',
        activity: 'Ferry departure',
        location: 'HarbourFront Ferry Terminal, Singapore',
      ),
      const _JourneyStep(
        time: '08:00',
        activity: 'Arrive & transfer',
        location: 'Batam Centre Ferry Terminal',
      ),
      _JourneyStep(
        time: '09:00',
        activity: serviceName,
        location: providerName,
      ),
      const _JourneyStep(
        time: '12:30',
        activity: 'Lunch & rest',
        location: 'Batu Aji',
      ),
    ];
  }

  void _chooseItinerary(BuildContext context) {
    // Booking flow is out of scope for now; acknowledge the selection.
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Itinerary selected! Booking is coming soon.'),
      ),
    );
  }

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
                      Text(
                        'Itinerary Journey',
                        style: theme.textTheme.titleLarge,
                      ),
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
                        _HospitalCard(
                          imageUrl: imageUrl,
                          providerName: providerName,
                          serviceName: serviceName,
                          location: location,
                          appointment: appointment,
                        ),
                        const SizedBox(height: 16),
                        _JourneyTimelineCard(
                          date: _journeyDate,
                          steps: _steps(),
                        ),
                        const SizedBox(height: 16),
                        _ItinerarySummaryCard(
                          duration: duration,
                          price: price,
                          includedItems: _includedItems,
                        ),
                      ],
                    ),
                  ),
                ),
                // Pinned "choose itinerary" action.
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
                      onPressed: () => _chooseItinerary(context),
                      child: const Text('Choose This Itinerary'),
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

/// Top card: hospital photo, name, medical service, location, appointment.
class _HospitalCard extends StatelessWidget {
  const _HospitalCard({
    required this.imageUrl,
    required this.providerName,
    required this.serviceName,
    required this.location,
    required this.appointment,
  });

  final String imageUrl;
  final String providerName;
  final String serviceName;
  final String location;
  final String appointment;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final muted = AppColors.text.withValues(alpha: 0.75);

    return Container(
      decoration: _cardDecoration(),
      clipBehavior: Clip.antiAlias,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          _HospitalImage(imageUrl: imageUrl),
          Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  providerName,
                  style: theme.textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.w800,
                  ),
                ),
                const SizedBox(height: 6),
                Text(
                  serviceName,
                  style: theme.textTheme.bodyMedium?.copyWith(color: muted),
                ),
                const SizedBox(height: 12),
                _InfoRow(icon: LucideIcons.mapPin, text: location),
                const SizedBox(height: 6),
                _InfoRow(
                  icon: LucideIcons.clock,
                  text: 'Appointment · $appointment',
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

/// Full-width hospital photo with a themed placeholder when it can't load.
class _HospitalImage extends StatelessWidget {
  const _HospitalImage({required this.imageUrl});

  final String imageUrl;

  @override
  Widget build(BuildContext context) {
    final placeholder = Container(
      height: 150,
      width: double.infinity,
      color: AppColors.primary.withValues(alpha: 0.08),
      alignment: Alignment.center,
      child: const Icon(
        LucideIcons.imageOff,
        size: 40,
        color: AppColors.primary,
      ),
    );

    if (imageUrl.startsWith('assets/')) {
      return Image.asset(
        imageUrl,
        height: 150,
        width: double.infinity,
        fit: BoxFit.cover,
        errorBuilder: (_, _, _) => placeholder,
      );
    }
    return Image.network(
      imageUrl,
      height: 150,
      width: double.infinity,
      fit: BoxFit.cover,
      errorBuilder: (_, _, _) => placeholder,
    );
  }
}

/// An icon + text row (e.g. location, appointment).
class _InfoRow extends StatelessWidget {
  const _InfoRow({required this.icon, required this.text});

  final IconData icon;
  final String text;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(icon, size: 20, color: AppColors.text),
        const SizedBox(width: 6),
        Flexible(child: Text(text, style: theme.textTheme.bodySmall)),
      ],
    );
  }
}

/// Card with a date header and a vertical timeline of journey steps.
class _JourneyTimelineCard extends StatelessWidget {
  const _JourneyTimelineCard({required this.date, required this.steps});

  final String date;
  final List<_JourneyStep> steps;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: _cardDecoration(),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Icon(
                LucideIcons.calendarDays,
                size: 20,
                color: AppColors.primary,
              ),
              const SizedBox(width: 8),
              Text(
                date,
                style: theme.textTheme.titleMedium?.copyWith(
                  fontWeight: FontWeight.w800,
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          for (var i = 0; i < steps.length; i++)
            _TimelineItem(step: steps[i], isLast: i == steps.length - 1),
        ],
      ),
    );
  }
}

/// One timeline row: time, activity, and location, with a connector to the
/// next step.
class _TimelineItem extends StatelessWidget {
  const _TimelineItem({required this.step, required this.isLast});

  final _JourneyStep step;
  final bool isLast;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final muted = AppColors.text.withValues(alpha: 0.75);

    return IntrinsicHeight(
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          // Timeline rail: dot + connector to the next step.
          Column(
            children: [
              Container(
                width: 12,
                height: 12,
                margin: const EdgeInsets.only(top: 3),
                decoration: const BoxDecoration(
                  color: AppColors.primary,
                  shape: BoxShape.circle,
                ),
              ),
              if (!isLast)
                Expanded(
                  child: Container(
                    width: 2,
                    color: AppColors.primary.withValues(alpha: 0.35),
                  ),
                ),
            ],
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Padding(
              padding: EdgeInsets.only(bottom: isLast ? 0 : 24),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    step.time,
                    style: theme.textTheme.labelLarge?.copyWith(
                      color: AppColors.primary,
                      fontWeight: FontWeight.w700,
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    step.activity,
                    style: theme.textTheme.titleSmall?.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Icon(
                        LucideIcons.mapPin,
                        size: 14,
                        color: AppColors.text,
                      ),
                      const SizedBox(width: 4),
                      Expanded(
                        child: Text(
                          step.location,
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: muted,
                          ),
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}

/// Card with the itinerary summary: total duration, estimated cost, and what
/// is included.
class _ItinerarySummaryCard extends StatelessWidget {
  const _ItinerarySummaryCard({
    required this.duration,
    required this.price,
    required this.includedItems,
  });

  final String duration;
  final String price;
  final List<String> includedItems;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: _cardDecoration(),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Itinerary Summary',
            style: theme.textTheme.titleMedium?.copyWith(
              fontWeight: FontWeight.w800,
            ),
          ),
          const SizedBox(height: 16),
          _SummaryRow(
            icon: LucideIcons.calendarClock,
            label: 'Total Duration',
            value: duration,
          ),
          const SizedBox(height: 12),
          _SummaryRow(
            icon: LucideIcons.wallet,
            label: 'Cost Estimate',
            value: '$price/person',
          ),
          const SizedBox(height: 16),
          Text(
            'Includes',
            style: theme.textTheme.titleSmall?.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 8),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: [
              for (final item in includedItems) _IncludesChip(label: item),
            ],
          ),
        ],
      ),
    );
  }
}

/// A label/value row in the summary card (e.g. "Total Duration → 3 days").
class _SummaryRow extends StatelessWidget {
  const _SummaryRow({
    required this.icon,
    required this.label,
    required this.value,
  });

  final IconData icon;
  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final muted = AppColors.text.withValues(alpha: 0.75);

    return Row(
      children: [
        Icon(icon, size: 20, color: AppColors.primary),
        const SizedBox(width: 8),
        Text(label, style: theme.textTheme.bodyMedium?.copyWith(color: muted)),
        const Spacer(),
        Text(
          value,
          style: theme.textTheme.bodyMedium?.copyWith(
            fontWeight: FontWeight.w700,
          ),
        ),
      ],
    );
  }
}

/// A small pill showing an item the itinerary includes (e.g. "Ferry").
class _IncludesChip extends StatelessWidget {
  const _IncludesChip({required this.label});

  final String label;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
      decoration: BoxDecoration(
        color: AppColors.primary.withValues(alpha: 0.10),
        borderRadius: BorderRadius.circular(100),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          const Icon(LucideIcons.check, size: 14, color: AppColors.primary),
          const SizedBox(width: 4),
          Text(
            label,
            style: theme.textTheme.bodySmall?.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
        ],
      ),
    );
  }
}
