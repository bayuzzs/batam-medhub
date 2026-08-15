import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:mobile/data/repository/journey_repository.dart';
import 'package:mobile/models/journey.dart';
import 'package:mobile/models/plan_option.dart';
import 'package:mobile/models/structured_intent.dart';
import 'package:mobile/models/trip_request.dart';
import 'package:mobile/ui/chat/chat_message.dart';

import '../journey/providers.dart';

/// Whether the chat is waiting on a journey-orchestration call.
enum ChatStatus { idle, sending }

/// Immutable chat state exposed to the UI.
class ChatState {
  const ChatState({
    this.messages = const [],
    this.status = ChatStatus.idle,
    this.tripRequest,
    this.journey,
  });

  /// The in-memory conversation transcript (client-side only).
  final List<ChatMessage> messages;

  /// Whether a journey-orchestration call is in flight.
  final ChatStatus status;

  /// The current trip request driving the conversation, once created.
  final TripRequest? tripRequest;

  /// The confirmed journey, once a plan option is selected.
  final JourneyDetail? journey;

  ChatState copyWith({
    List<ChatMessage>? messages,
    ChatStatus? status,
    TripRequest? tripRequest,
    bool clearTripRequest = false,
    JourneyDetail? journey,
    bool clearJourney = false,
  }) {
    return ChatState(
      messages: messages ?? this.messages,
      status: status ?? this.status,
      tripRequest: clearTripRequest ? null : (tripRequest ?? this.tripRequest),
      journey: clearJourney ? null : (journey ?? this.journey),
    );
  }
}

/// Owns the chat conversation and drives it through the trip-request state
/// machine backed by [JourneyRepository]:
///
/// - `send(text)` — append the user bubble, `POST /v1/trip-requests`, then
///   either surface the clarification question or generate plans.
/// - `answerClarification(text)` — `PATCH .../intent`, then plans.
/// - `selectPlanOption(option)` — `POST /plan-options/{id}/confirm` → journey.
///
/// Transcript is in-memory only (the spec has no chat resource). While a
/// request is in flight a transient "typing…" bubble is shown.
class ChatController extends Notifier<ChatState> {
  JourneyRepository get _repository => ref.read(journeyRepositoryProvider);

  @override
  ChatState build() {
    return const ChatState(
      messages: [
        ChatMessage.assistant(
          'Hi! I\'m your medical-journey assistant. Tell me what you need — '
          'for example a same-day check-up in Batam — and I\'ll orchestrate it.',
        ),
      ],
    );
  }

  /// Start a new trip request from the patient's [text].
  ///
  /// When the current trip request is awaiting clarification (the assistant
  /// asked a question and no plan is pending), the text is treated as the
  /// patient's answer and routed to [answerClarification] instead.
  Future<void> send(String text) async {
    final trimmed = text.trim();
    if (trimmed.isEmpty || state.status == ChatStatus.sending) return;

    final current = state.tripRequest;
    final awaitingClarification =
        current != null &&
        current.intent.resolution == IntentResolution.needsClarification;
    if (awaitingClarification) {
      await answerClarification(trimmed);
      return;
    }

    _append(ChatMessage.user(trimmed));
    await _run(() async {
      final detail = await _repository.createTripRequest(
        prompt: trimmed,
        locale: 'en',
      );
      state = state.copyWith(tripRequest: detail.tripRequest);
      await _handleTripRequestDetail(detail);
    });
  }

  /// Answer the assistant's clarification question (or supply corrections).
  Future<void> answerClarification(String text) async {
    if (state.status == ChatStatus.sending) return;
    final trimmed = text.trim();
    final tripRequestId = state.tripRequest?.id;
    if (trimmed.isEmpty || tripRequestId == null) return;

    _append(ChatMessage.user(trimmed));
    await _run(() async {
      final detail = await _repository.amendIntent(
        tripRequestId: tripRequestId,
        answer: trimmed,
      );
      state = state.copyWith(tripRequest: detail.tripRequest);
      await _handleTripRequestDetail(detail);
    });
  }

  /// Confirm a plan option, creating the journey.
  Future<void> selectPlanOption(PlanOption option) async {
    if (state.status == ChatStatus.sending) return;
    _append(ChatMessage.assistant('Booking your journey — one moment…'));
    await _run(() async {
      final journey = await _repository.confirmPlanOption(
        planOptionId: option.id,
      );
      state = state.copyWith(journey: journey);
      _append(
        ChatMessage.assistant(
          'Your journey is confirmed! I\'ve arranged your ferry, transport, '
          'and hospital appointment.',
        ),
      );
    });
  }

  Future<void> _handleTripRequestDetail(TripRequestDetail detail) async {
    final intent = detail.tripRequest.intent;
    switch (intent.resolution) {
      case IntentResolution.needsClarification:
        _append(
          ChatMessage.assistant(
            intent.clarificationQuestion ??
                'Could you give me a few more '
                    'details so I can plan your trip?',
          ),
        );
      case IntentResolution.matched:
        final result = await _repository.generatePlans(
          tripRequestId: detail.tripRequest.id,
        );
        state = state.copyWith(tripRequest: result.tripRequest);
        if (result.options.isEmpty) {
          _append(
            ChatMessage.assistant(
              result.noMatchReasons.isEmpty
                  ? 'I couldn\'t find any feasible options for that trip.'
                  : result.noMatchReasons.join(' '),
            ),
          );
        } else {
          _append(ChatMessage.planOptions(result.options));
        }
      case IntentResolution.unsupportedService:
        _append(
          ChatMessage.assistant(
            intent.unsupportedReason ?? 'That service isn\'t supported yet.',
          ),
        );
      case IntentResolution.outOfScope:
        _append(
          ChatMessage.assistant(
            intent.outOfScopeReason ??
                'Sorry, that\'s outside what I can help '
                    'with on this platform.',
          ),
        );
    }
  }

  /// Shows a typing bubble while [action] runs, appends the assistant's
  /// outcome, and resets the sending status in `finally`.
  Future<void> _run(Future<void> Function() action) async {
    _setSending(true);
    _setTyping(true);
    try {
      await action();
    } on JourneyException catch (error) {
      _append(
        ChatMessage.assistant(
          'Sorry, I couldn\'t do that: '
          '${error.message}',
        ),
      );
    } catch (_) {
      _append(ChatMessage.assistant('Something went wrong. Please try again.'));
    } finally {
      _setTyping(false);
      _setSending(false);
    }
  }

  void _append(ChatMessage message) {
    state = state.copyWith(messages: [...state.messages, message]);
  }

  void _setSending(bool sending) {
    state = state.copyWith(
      status: sending ? ChatStatus.sending : ChatStatus.idle,
    );
  }

  void _setTyping(bool typing) {
    final messages = state.messages;
    if (typing) {
      if (messages.isNotEmpty && messages.last.isTyping) return;
      state = state.copyWith(
        messages: [...messages, const ChatMessage.typing()],
      );
    } else {
      // The typing bubble may no longer be the last message (the assistant's
      // reply is appended before `finally` runs), so remove it wherever it is.
      final withoutTyping = messages.where((m) => !m.isTyping).toList();
      if (withoutTyping.length != messages.length) {
        state = state.copyWith(messages: withoutTyping);
      }
    }
  }
}
