import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lucide_icons_flutter/lucide_icons.dart';

import 'package:mobile/application/auth/providers.dart'
    show authControllerProvider;
import 'package:mobile/application/chat/chat_controller.dart' show ChatState;
import 'package:mobile/application/chat/providers.dart';
import 'package:mobile/models/plan_option.dart';
import 'package:mobile/ui/core/navigation/app_router.dart';
import 'package:mobile/ui/core/theme/app_assets.dart';
import 'package:mobile/ui/core/theme/app_colors.dart';
import 'package:mobile/ui/core/theme/app_spacing.dart';
import 'package:mobile/ui/core/widgets/app_container.dart';
import 'package:mobile/ui/core/widgets/primary_radial_gradient.dart';
import 'package:mobile/ui/chat/chat_item.dart';
import 'package:mobile/ui/chat/chat_message.dart';
import 'package:mobile/ui/chat/plan_option_card.dart';

/// AI chat screen — the app's primary authenticated screen.
///
/// Has a pinned top bar with **History** (top left) and **Profile** (top
/// right) actions; both push full-screen pages. Greets the user, shows a card
/// introducing the AI assistant, and pins a multiline message input at the
/// bottom. The conversation is driven by [chatControllerProvider], which walks
/// the trip-request state machine: the patient's typed text becomes a user
/// bubble, the assistant replies with a clarification question or plan cards,
/// and selecting a plan books the journey.
class ChatPage extends ConsumerStatefulWidget {
  const ChatPage({super.key});

  @override
  ConsumerState<ChatPage> createState() => _ChatPageState();
}

class _ChatPageState extends ConsumerState<ChatPage> {
  final _messageController = TextEditingController();
  final _scrollController = ScrollController();

  @override
  void dispose() {
    _messageController.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  void _sendMessage() {
    final message = _messageController.text.trim();
    if (message.isEmpty) return;
    _messageController.clear();
    ref.read(chatControllerProvider.notifier).send(message);
  }

  /// The patient's initial for the user avatar, from the authenticated
  /// profile (falls back to a neutral letter when unavailable).
  String get _userInitial {
    final session = ref.read(authControllerProvider).session;
    final name = session?.profile.fullName.trim() ?? '';
    return name.isEmpty ? 'H' : name[0].toUpperCase();
  }

  void _scrollToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!_scrollController.hasClients) return;
      _scrollController.animateTo(
        _scrollController.position.maxScrollExtent,
        duration: const Duration(milliseconds: 250),
        curve: Curves.easeOut,
      );
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final chat = ref.watch(chatControllerProvider);

    ref.listen<ChatState>(chatControllerProvider, (_, _) {
      _scrollToBottom();
    });

    return Scaffold(
      body: Stack(
        children: [
          const PrimaryRadialGradient(startOpacity: 0.20),
          SafeArea(
            child: Column(
              children: [
                // Pinned top bar (does not scroll): History top-left, Profile
                // top-right. Both push full-screen pages on top of the chat.
                Padding(
                  padding: const EdgeInsets.fromLTRB(
                    AppSpacing.screen,
                    12,
                    AppSpacing.screen,
                    0,
                  ),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      _TopBarButton(
                        key: const Key('chat_history_button'),
                        icon: LucideIcons.history,
                        tooltip: 'History',
                        onPressed: () => context.pushHistory(),
                      ),
                      _TopBarButton(
                        key: const Key('chat_profile_button'),
                        icon: LucideIcons.user,
                        tooltip: 'Profile',
                        onPressed: () => context.pushProfile(),
                      ),
                    ],
                  ),
                ),
                Expanded(
                  child: AppContainer(
                    vertical: 40,
                    child: ListView(
                      controller: _scrollController,
                      padding: const EdgeInsets.only(bottom: 16),
                      children: [
                        const SizedBox(height: 24),
                        Text(
                          'Hi, Name',
                          style: theme.textTheme.titleMedium?.copyWith(
                            color: AppColors.primary,
                          ),
                        ),
                        const SizedBox(height: 8),
                        // Decorative landmark behind the texts, anchored to
                        // the bottom-right of the header. Overflows its stack
                        // on purpose, so clipping is disabled.
                        Stack(
                          clipBehavior: Clip.none,
                          children: [
                            Positioned(
                              right: -40,
                              bottom: -40,
                              child: Image.asset(
                                AppAssets.barelang,
                                height: 200,
                              ),
                            ),
                            Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  'Let\'s orchestrate your medical journey.',
                                  style: theme.textTheme.headlineMedium,
                                ),
                                const SizedBox(height: 8),
                                Text(
                                  'One request, One connected journey across '
                                  'care, travel and stay.',
                                  style: theme.textTheme.bodyLarge,
                                ),
                              ],
                            ),
                          ],
                        ),
                        const SizedBox(height: 24),
                        const _AssistantCard(),
                        const SizedBox(height: 24),
                        // The conversation scrolls together with the header
                        // and assistant card as one scrollable column.
                        for (final message in chat.messages)
                          Padding(
                            padding: const EdgeInsets.only(bottom: 12),
                            child: _buildMessage(message),
                          ),
                      ],
                    ),
                  ),
                ),
                Padding(
                  padding: EdgeInsets.fromLTRB(
                    AppSpacing.screen,
                    12,
                    AppSpacing.screen,
                    AppSpacing.screen,
                  ),
                  child: _ChatInput(
                    controller: _messageController,
                    onSend: _sendMessage,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildMessage(ChatMessage message) {
    if (message.planOptions != null) {
      return _PlanOptionsList(
        options: message.planOptions!,
        onViewDetails: (option) => context.pushPlanDetail(option),
      );
    }
    if (message.isTyping) {
      return const ChatItem.assistant(child: _TypingIndicator());
    }
    final text = message.text ?? '';
    return message.role == ChatRole.user
        ? ChatItem.user(message: text, userInitial: _userInitial)
        : ChatItem.assistant(message: text);
  }
}

/// A column of [PlanOptionCard]s rendered full-width (outside the chat bubble
/// so each card can use the whole content width). Cards only open the plan
/// detail screen; booking happens there.
class _PlanOptionsList extends StatelessWidget {
  const _PlanOptionsList({required this.options, required this.onViewDetails});

  final List<PlanOption> options;
  final ValueChanged<PlanOption> onViewDetails;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        for (final option in options) ...[
          PlanOptionCard(
            option: option,
            onViewDetails: () => onViewDetails(option),
          ),
          const SizedBox(height: 12),
        ],
      ],
    );
  }
}

/// Transient "assistant is typing" indicator shown inside a chat bubble.
class _TypingIndicator extends StatelessWidget {
  const _TypingIndicator();

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Text('typing', style: theme.textTheme.bodyMedium),
        const SizedBox(width: 6),
        const _TypingDots(),
      ],
    );
  }
}

/// Three animated dots for the typing indicator.
class _TypingDots extends StatefulWidget {
  const _TypingDots();

  @override
  State<_TypingDots> createState() => _TypingDotsState();
}

class _TypingDotsState extends State<_TypingDots>
    with SingleTickerProviderStateMixin {
  late final AnimationController _controller = AnimationController(
    vsync: this,
    duration: const Duration(milliseconds: 900),
  )..repeat();

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: _controller,
      builder: (context, _) {
        return Row(
          mainAxisSize: MainAxisSize.min,
          children: List.generate(3, (index) {
            final phase = (_controller.value - index * 0.15) % 1.0;
            final t = (1 - (phase * 2 - 1).abs()).clamp(0.0, 1.0);
            final opacity = (0.3 + 0.7 * t).toDouble();
            return Padding(
              padding: const EdgeInsets.symmetric(horizontal: 2),
              child: Opacity(
                opacity: opacity,
                child: const Icon(LucideIcons.circle, size: 6),
              ),
            );
          }),
        );
      },
    );
  }
}

/// Circular top-bar icon button used in the chat screen's pinned header.
class _TopBarButton extends StatelessWidget {
  const _TopBarButton({
    super.key,
    required this.icon,
    required this.tooltip,
    required this.onPressed,
  });

  final IconData icon;
  final String tooltip;
  final VoidCallback onPressed;

  @override
  Widget build(BuildContext context) {
    return Tooltip(
      message: tooltip,
      child: IconButton(
        onPressed: onPressed,
        style: IconButton.styleFrom(
          backgroundColor: AppColors.background,
          foregroundColor: AppColors.heading,
          side: BorderSide(color: AppColors.inputBorder),
          minimumSize: const Size(44, 44),
        ),
        icon: Icon(icon, size: 22),
      ),
    );
  }
}

/// Card introducing the AI assistant: icon, description, and a robot image.
class _AssistantCard extends StatelessWidget {
  const _AssistantCard();

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(24),
        border: Border.all(color: AppColors.inputBorder),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            width: 44,
            height: 44,
            decoration: BoxDecoration(
              color: AppColors.primary,
              borderRadius: BorderRadius.circular(100),
            ),
            child: const Icon(LucideIcons.sparkles, color: Colors.white),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Talk with your AI Assistant',
                  style: theme.textTheme.titleSmall?.copyWith(
                    fontWeight: FontWeight.w400,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  'Tell me what you need and I\'ll orchestrate the best '
                  'medical journey for you.',
                  style: theme.textTheme.bodySmall,
                ),
              ],
            ),
          ),
          const SizedBox(width: 6),
          Image.asset(AppAssets.robot, height: 42),
        ],
      ),
    );
  }
}

/// Message input pinned at the bottom of the chat screen, with a send button
/// on the right side.
class _ChatInput extends StatelessWidget {
  const _ChatInput({required this.controller, required this.onSend});

  final TextEditingController controller;
  final VoidCallback onSend;

  @override
  Widget build(BuildContext context) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.end,
      children: [
        // Multi-line input that grows as the user types (1–4 rows). The send
        // button sits beside it so it stays bottom-aligned as the box grows.
        Expanded(
          child: TextField(
            controller: controller,
            minLines: 1,
            maxLines: 4,
            keyboardType: TextInputType.multiline,
            onSubmitted: (_) => onSend(),
            style: Theme.of(context).textTheme.bodyLarge,
            decoration: const InputDecoration(
              hintText: 'Type your message...',
              contentPadding: EdgeInsets.symmetric(
                horizontal: 18,
                vertical: 16,
              ),
            ),
          ),
        ),
        const SizedBox(width: 8),
        IconButton.filled(
          key: const Key('chat_send_button'),
          onPressed: onSend,
          icon: const Icon(LucideIcons.send, size: 22),
        ),
      ],
    );
  }
}
