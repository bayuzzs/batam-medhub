import 'package:flutter/material.dart';

import 'package:mobile/ui/core/app_colors.dart';

/// A selectable option card (white card with a border, no radio indicator).
///
/// Tapping the card calls [onTap]; [selected] controls the highlight
/// (primary border). Selection state is managed externally — see
/// [ChatOptions], which wires several of these together.
class ChatOption extends StatelessWidget {
  const ChatOption({
    super.key,
    required this.label,
    this.subtitle,
    required this.selected,
    required this.onTap,
  });

  /// The main option text.
  final String label;

  /// Optional secondary text shown below [label].
  final String? subtitle;

  /// Whether this option is currently selected.
  final bool selected;

  /// Called when the option is tapped.
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Container(
          width: double.infinity,
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
          decoration: BoxDecoration(
            color: AppColors.background,
            borderRadius: BorderRadius.circular(12),
            border: Border.all(
              color: selected ? AppColors.primary : AppColors.inputBorder,
              width: selected ? 2 : 1,
            ),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                label,
                style: theme.textTheme.bodyMedium?.copyWith(
                  fontWeight: FontWeight.w600,
                ),
              ),
              if (subtitle != null) ...[
                const SizedBox(height: 2),
                Text(
                  subtitle!,
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: theme.colorScheme.onSurfaceVariant,
                  ),
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }
}

/// An option for [ChatOptions].
typedef ChatOptionData = ({String label, String? subtitle});

/// A column of selectable [ChatOption]s that manages its own selection.
///
/// Only one option can be selected at a time; tapping an option selects it
/// and calls [onChanged] with its index.
class ChatOptions extends StatefulWidget {
  const ChatOptions({super.key, required this.options, this.onChanged});

  /// The option labels (and optional subtitles).
  final List<ChatOptionData> options;

  /// Called with the index of the selected option.
  final ValueChanged<int>? onChanged;

  @override
  State<ChatOptions> createState() => _ChatOptionsState();
}

class _ChatOptionsState extends State<ChatOptions> {
  int? _selected;

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        for (final (index, option) in widget.options.indexed) ...[
          if (index > 0) const SizedBox(height: 8),
          ChatOption(
            label: option.label,
            subtitle: option.subtitle,
            selected: _selected == index,
            onTap: () {
              setState(() => _selected = index);
              widget.onChanged?.call(index);
            },
          ),
        ],
      ],
    );
  }
}
