import 'package:flutter/material.dart';
import 'package:lucide_icons_flutter/lucide_icons.dart';

import 'package:mobile/ui/core/theme/app_assets.dart';
import 'package:mobile/ui/core/theme/app_colors.dart';
import 'package:mobile/ui/core/navigation/app_router.dart';
import 'package:mobile/ui/core/widgets/app_container.dart';
import 'package:mobile/ui/core/widgets/itenary_option_card.dart';
import 'package:mobile/ui/core/theme/app_spacing.dart';
import 'package:mobile/ui/core/widgets/primary_radial_gradient.dart';
import 'package:mobile/ui/chat/chat_item.dart';
import 'package:mobile/ui/chat/chat_options.dart';

/// AI chat screen — the app's primary authenticated screen.
///
/// Has a pinned top bar with **History** (top left) and **Profile** (top
/// right) actions; both push full-screen pages. Greets the user, shows a card
/// introducing the AI assistant, and pins a message input at the bottom of
/// the screen. Pure UI — chat/AI integration is not wired up yet.
class ChatPage extends StatefulWidget {
  const ChatPage({super.key});

  @override
  State<ChatPage> createState() => _ChatPageState();
}

class _ChatPageState extends State<ChatPage> {
  final _messageController = TextEditingController();

  @override
  void dispose() {
    _messageController.dispose();
    super.dispose();
  }

  void _sendMessage() {
    final message = _messageController.text.trim();
    if (message.isEmpty) return;
    // TODO(chat): send the message to the AI assistant service once integrated.
    _messageController.clear();
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
                    scrollable: true,
                    vertical: 40,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const SizedBox(height: 24),
                        Text(
                          'Hi, Name',
                          style: theme.textTheme.titleMedium?.copyWith(
                            color: AppColors.primary,
                          ),
                        ),
                        const SizedBox(height: 8),
                        Stack(
                          clipBehavior: Clip.none,
                          children: [
                            // Decorative landmark behind the texts, anchored to
                            // the bottom-right of the header. Overflows the stack
                            // bounds on purpose, so clipping is disabled.
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
                        const ChatItem.user(
                          message:
                              'I need help planning my medical trip to Batam',
                          userInitial: 'H',
                        ),
                        const SizedBox(height: 12),
                        const ChatItem.assistant(
                          message:
                              'Of course! I\'ll orchestrate the best '
                              'medical journey for you.',
                        ),
                        const SizedBox(height: 12),
                        const ChatItem.assistant(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text('Which hospital would you prefer?'),
                              SizedBox(height: 12),
                              ChatOptions(
                                options: [
                                  (
                                    label: 'RS Awal Bros Batam',
                                    subtitle: 'Batu Aji · 5 km',
                                  ),
                                  (
                                    label: 'RS Hermina Batam',
                                    subtitle: 'Batam Center · 7 km',
                                  ),
                                  (
                                    label: 'RSUD Embung Fatimah',
                                    subtitle: 'Sekupang · 9 km',
                                  ),
                                ],
                              ),
                            ],
                          ),
                        ),
                        const SizedBox(height: 12),
                        const ChatItem.assistant(
                          message: 'Here\'s a great itinerary for you:',
                        ),
                        const SizedBox(height: 12),
                        // Full-width itinerary card, outside the chat bubble so
                        // it can use the whole content width.
                        ItenaryOptionCard(
                          imageUrl: AppAssets.hospitalAwalBros,
                          providerName: 'RS Awal Bros Batam',
                          serviceName: 'Cardiac Screening Package',
                          location: 'Batu Aji · 5 km',
                          appointment: 'Tomorrow, 09:00',
                          rating: 4.8,
                          reviewCount: 212,
                          duration: '3 days',
                          price: 'IDR 4.5jt',
                          recommended: true,
                          onViewDetails: () => context.pushItinerary(),
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
