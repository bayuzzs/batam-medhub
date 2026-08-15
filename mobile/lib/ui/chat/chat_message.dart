import 'package:mobile/models/plan_option.dart';

import 'chat_item.dart' show ChatRole;

/// A single message in the in-memory chat transcript.
///
/// This is a **client-side** model (not an OpenAPI schema) that renders the
/// trip-request conversation: user text bubbles, assistant text, plan-option
/// cards, and a transient "typing…" indicator.
class ChatMessage {
  const ChatMessage._({
    required this.role,
    this.text,
    this.planOptions,
    this.isTyping = false,
  });

  /// A message the patient typed.
  const ChatMessage.user(String text) : this._(role: ChatRole.user, text: text);

  /// A plain text reply from the assistant.
  const ChatMessage.assistant(String text)
    : this._(role: ChatRole.assistant, text: text);

  /// A bubble rendering one or more [PlanOption] cards.
  const ChatMessage.planOptions(List<PlanOption> options)
    : this._(role: ChatRole.assistant, planOptions: options);

  /// Transient "assistant is typing" indicator bubble.
  const ChatMessage.typing() : this._(role: ChatRole.assistant, isTyping: true);

  final ChatRole role;

  /// Message text; null when the bubble renders [planOptions] or is typing.
  final String? text;

  /// Plan options rendered as selectable cards, when present.
  final List<PlanOption>? planOptions;

  /// Whether this is the transient typing indicator.
  final bool isTyping;
}
