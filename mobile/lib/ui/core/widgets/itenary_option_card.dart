import 'package:flutter/material.dart';
import 'package:lucide_icons_flutter/lucide_icons.dart';

import 'package:mobile/ui/core/theme/app_colors.dart';

/// A card describing an itinerary option (a provider's medical travel
/// package).
///
/// Shows an optional [Recommended] banner, the provider's image + name +
/// service, key info (location, appointment), rating / language / duration
/// chips, an estimated total price, and a "View Details Itinerary" action.
///
/// Styling follows the app theme: a white card with a primary-tinted border
/// and rounded corners (matching `chat_options.dart`).
class ItenaryOptionCard extends StatelessWidget {
  const ItenaryOptionCard({
    super.key,
    required this.imageUrl,
    required this.providerName,
    required this.serviceName,
    required this.location,
    required this.appointment,
    required this.rating,
    required this.reviewCount,
    required this.duration,
    required this.price,
    this.recommended = true,
    this.onViewDetails,
  });

  /// Image shown for the provider. Supports local `assets/...` paths and
  /// remote URLs; a placeholder is rendered when it can't be loaded.
  final String imageUrl;

  /// Provider / hospital name.
  final String providerName;

  /// The offered service (e.g. "Cardiac Screening Package").
  final String serviceName;

  /// Provider location (e.g. "Batu Aji · 5 km").
  final String location;

  /// Appointment availability text (e.g. "Tomorrow, 09:00").
  final String appointment;

  /// Provider rating out of 5.
  final double rating;

  /// Number of reviews behind [rating].
  final int reviewCount;

  /// Trip duration, shown as a chip (e.g. "3 days").
  final String duration;

  /// Estimated price per person (e.g. "IDR 4.500.000").
  final String price;

  /// Whether to show the "Recommended" banner.
  final bool recommended;

  /// Called when "View Details Itinerary" is tapped.
  final VoidCallback? onViewDetails;

  static const double _imageWidth = 100;
  static const double _imageHeight = 50;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final muted = AppColors.text.withValues(alpha: 0.75);

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
          if (recommended) ...[
            _RecommendedBanner(),
            const SizedBox(height: 24),
          ],
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              _ProviderImage(imageUrl: imageUrl),
              const SizedBox(width: 16),
              Expanded(
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
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 4),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: [
              _InfoItem(icon: LucideIcons.mapPin, text: location),
              _InfoItem(
                icon: LucideIcons.clock,
                text: 'Appointment · $appointment',
              ),
            ],
          ),
          const SizedBox(height: 12),
          Wrap(
            spacing: 4,
            runSpacing: 4,
            children: [
              _OptionChip(
                icon: LucideIcons.award,
                label: '$rating ($reviewCount reviews)',
                backgroundColor: const Color(0xFFFFE6A8),
                foregroundColor: const Color(0xFF70451D),
              ),
              const _OptionChip(
                icon: LucideIcons.globe,
                label: 'English Speaking',
                backgroundColor: Color(0xFFC7E9FF),
                foregroundColor: Color(0xFF1F3E6D),
              ),
              _OptionChip(
                icon: LucideIcons.calendarClock,
                label: duration,
                backgroundColor: const Color(0xFFC7F4DC),
                foregroundColor: const Color(0xFF1D5C3F),
              ),
            ],
          ),
          const SizedBox(height: 12),
          _PriceSection(price: price),
          const SizedBox(height: 12),
          FilledButton(
            onPressed: onViewDetails,
            child: Text('View Itinerary Details'),
          ),
        ],
      ),
    );
  }
}

/// "Recommended" banner shown at the top of the card.
class _RecommendedBanner extends StatelessWidget {
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

/// Provider image, with a themed placeholder when it can't be loaded.
class _ProviderImage extends StatelessWidget {
  const _ProviderImage({required this.imageUrl});

  final String imageUrl;

  @override
  Widget build(BuildContext context) {
    final placeholder = Container(
      width: ItenaryOptionCard._imageWidth,
      height: ItenaryOptionCard._imageHeight,
      color: AppColors.primary.withValues(alpha: 0.08),
      alignment: Alignment.center,
      child: const Icon(
        LucideIcons.imageOff,
        size: 32,
        color: AppColors.primary,
      ),
    );

    return ClipRRect(
      borderRadius: BorderRadius.circular(16),
      child: imageUrl.startsWith('assets/')
          ? Image.asset(
              imageUrl,
              width: ItenaryOptionCard._imageWidth,
              height: ItenaryOptionCard._imageHeight,
              fit: BoxFit.cover,
              errorBuilder: (_, _, _) => placeholder,
            )
          : Image.network(
              imageUrl,
              width: ItenaryOptionCard._imageWidth,
              height: ItenaryOptionCard._imageHeight,
              fit: BoxFit.cover,
              errorBuilder: (_, _, _) => placeholder,
            ),
    );
  }
}

/// An icon + label row (e.g. location, appointment).
class _InfoItem extends StatelessWidget {
  const _InfoItem({required this.icon, required this.text});

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
        // Flexible so the text wraps instead of overflowing when the card is
        // rendered in a narrow context (e.g. inside a chat bubble).
        Flexible(child: Text(text, style: theme.textTheme.bodySmall)),
      ],
    );
  }
}

/// A colored metadata chip (rating, language, duration).
class _OptionChip extends StatelessWidget {
  const _OptionChip({
    required this.icon,
    required this.label,
    required this.backgroundColor,
    required this.foregroundColor,
  });

  final IconData icon;
  final String label;
  final Color backgroundColor;
  final Color foregroundColor;

  @override
  Widget build(BuildContext context) {
    return Chip(
      avatar: Icon(icon, color: foregroundColor),
      label: Text(label),
      backgroundColor: backgroundColor,
      labelStyle: TextStyle(color: foregroundColor),
    );
  }
}

/// Bordered price block: "Estimated Total", what's included, price per person.
class _PriceSection extends StatelessWidget {
  const _PriceSection({required this.price});

  final String price;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Includes Ferry · Transport · Medical · Hotel',
          style: theme.textTheme.bodyMedium,
        ),
        const SizedBox(height: 6),
        Row(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            Expanded(
              child: Text('Estimated Total', style: theme.textTheme.bodyMedium),
            ),
            Text(
              '$price/person',
              style: theme.textTheme.headlineSmall?.copyWith(
                fontWeight: FontWeight.w800,
                color: AppColors.primary,
              ),
            ),
          ],
        ),
      ],
    );
  }
}
