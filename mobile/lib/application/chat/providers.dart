import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'chat_controller.dart';

/// The chat conversation state. Drives the chat UI through the trip-request
/// state machine backed by [JourneyRepository] (see `journey/providers.dart`).
final chatControllerProvider = NotifierProvider<ChatController, ChatState>(
  ChatController.new,
);
