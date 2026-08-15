import 'package:flutter/material.dart';

import 'package:mobile/ui/core/theme/app_colors.dart';

/// Text input with its label rendered above the box, not inside the border.
///
/// Wraps a [TextFormField] in a column with [label] on top, so the label never
/// overlaps the input border. Prefer this over raw `TextFormField` with
/// `InputDecoration(labelText:)` for form fields.
class AppTextField extends StatelessWidget {
  const AppTextField({
    super.key,
    required this.label,
    this.controller,
    this.validator,
    this.keyboardType,
    this.textInputAction,
    this.obscureText = false,
    this.autocorrect = true,
    this.textCapitalization = TextCapitalization.none,
    this.prefixIcon,
    this.suffixIcon,
  });

  /// The label shown above the input box.
  final String label;

  /// Controller for the underlying [TextFormField].
  final TextEditingController? controller;

  /// Form validator for the underlying [TextFormField].
  final String? Function(String?)? validator;

  /// Keyboard type for the underlying [TextFormField].
  final TextInputType? keyboardType;

  /// Keyboard action for the underlying [TextFormField].
  final TextInputAction? textInputAction;

  /// Whether the input is obscured (e.g. passwords).
  final bool obscureText;

  /// Whether autocorrect is enabled.
  final bool autocorrect;

  /// Text capitalization behavior (e.g. words for names).
  final TextCapitalization textCapitalization;

  /// Icon shown at the start of the input.
  final Widget? prefixIcon;

  /// Widget shown at the end of the input (e.g. visibility toggle).
  final Widget? suffixIcon;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: theme.textTheme.bodyMedium?.copyWith(
            color: theme.colorScheme.onSurfaceVariant,
            fontWeight: FontWeight.w500,
          ),
        ),
        const SizedBox(height: 8),
        TextFormField(
          controller: controller,
          validator: validator,
          keyboardType: keyboardType,
          textInputAction: textInputAction,
          obscureText: obscureText,
          autocorrect: autocorrect,
          textCapitalization: textCapitalization,
          decoration: InputDecoration(
            prefixIcon: prefixIcon,
            suffixIcon: suffixIcon,
            // Icons inside the field use the regular text color, not the
            // white icon-button default.
            prefixIconColor: AppColors.text,
            suffixIconColor: AppColors.text,
          ),
        ),
      ],
    );
  }
}
