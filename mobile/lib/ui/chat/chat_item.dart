import 'package:flutter/material.dart';
import 'package:lucide_icons_flutter/lucide_icons.dart';

import 'package:mobile/ui/core/app_colors.dart';

/// Role of a chat message, controlling the bubble and avatar styling.
enum ChatRole { user, assistant }

/// A single chat message with its sender avatar on the outer edge.
///
/// - [ChatRole.user] — bubble with a [AppColors.primary] background and
///   [AppColors.heading]-colored text, aligned to the right with a circle
///   avatar (user initial) on its right.
/// - [ChatRole.assistant] — white bubble with the standard [AppColors.inputBorder]
///   border, aligned to the left with a sparkles icon on its left.
///
/// Prefer `ChatItem.user(...)` / `ChatItem.assistant(...)` over constructing
/// with [ChatRole] directly.
class ChatItem extends StatelessWidget {
  /// A message sent by the user.
  const ChatItem.user({
    super.key,
    this.message,
    this.child,
    required this.userInitial,
    this.indent = 80,
  }) : role = ChatRole.user;

  /// A message sent by the AI assistant.
  const ChatItem.assistant({
    super.key,
    this.message,
    this.child,
    this.indent = 80,
  }) : role = ChatRole.assistant,
       userInitial = null;

  final ChatRole role;

  /// The message text. Ignored when [child] is provided.
  final String? message;

  /// Custom content rendered inside the bubble, e.g. text with selectable
  /// options below it. Prefer this over [message] for rich content.
  final Widget? child;

  /// The user's initial, shown in the avatar circle. Only used for
  /// [ChatRole.user].
  final String? userInitial;

  /// Horizontal inset: on user messages it pushes the bubble and avatar toward
  /// the right (padding on the left); on assistant messages it pushes the
  /// bubble and icon toward the left (padding on the right of the content).
  final double indent;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isUser = role == ChatRole.user;

    // User: avatar on the right of its text (message hugs the right edge).
    // Assistant: icon on the left of its text (message hugs the left edge).
    final children = isUser
        ? <Widget>[
            Expanded(child: _bubble(theme, isUser)),
            const SizedBox(width: 8),
            _avatar(theme, isUser),
          ]
        : <Widget>[
            _avatar(theme, isUser),
            const SizedBox(width: 8),
            Expanded(child: _bubble(theme, isUser)),
          ];

    return Padding(
      padding: EdgeInsets.only(
        left: isUser ? indent : 0,
        right: isUser ? 0 : indent,
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.end,
        children: children,
      ),
    );
  }

  Widget _bubble(ThemeData theme, bool isUser) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      decoration: BoxDecoration(
        color: isUser ? AppColors.primary : AppColors.background,
        borderRadius: BorderRadius.circular(12),
        border: isUser ? null : Border.all(color: AppColors.inputBorder),
      ),
      child:
          child ??
          Text(
            message ?? '',
            style: theme.textTheme.bodyMedium?.copyWith(
              color: isUser ? AppColors.heading : null,
            ),
          ),
    );
  }

  Widget _avatar(ThemeData theme, bool isUser) {
    if (isUser) {
      return CircleAvatar(
        radius: 15,
        backgroundColor: AppColors.primary,
        child: Text(
          userInitial ?? '?',
          style: theme.textTheme.bodyMedium?.copyWith(
            color: Colors.white,
            fontWeight: FontWeight.w600,
          ),
        ),
      );
    }

    return Container(
      width: 30,
      height: 30,
      decoration: const BoxDecoration(
        color: AppColors.primary,
        shape: BoxShape.circle,
      ),
      child: const Icon(LucideIcons.sparkles, color: Colors.white, size: 18),
    );
  }
}
