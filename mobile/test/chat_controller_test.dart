// Unit tests for the chat controller's trip-request conversation, backed by
// the instant fake journey repository.

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'package:mobile/application/chat/chat_controller.dart';
import 'package:mobile/application/chat/providers.dart';
import 'package:mobile/application/journey/providers.dart';
import 'package:mobile/data/repository/fake_journey_repository.dart';
import 'package:mobile/data/repository/journey_repository.dart';
import 'package:mobile/models/journey.dart';
import 'package:mobile/models/medical_service.dart';
import 'package:mobile/models/plan_option.dart';
import 'package:mobile/models/structured_intent.dart';
import 'package:mobile/ui/chat/chat_item.dart' show ChatRole;

/// A repository that fails every journey-orchestration call, used to exercise
/// the controller's error bubble path.
class _ThrowingJourneyRepository implements JourneyRepository {
  const _ThrowingJourneyRepository();

  Never _fail() =>
      throw const JourneyException('backend is down', code: 'UPSTREAM_DOWN');

  @override
  Future<TripRequestDetail> createTripRequest({
    required String prompt,
    required String locale,
  }) async => _fail();

  @override
  Future<TripRequestDetail> amendIntent({
    required String tripRequestId,
    String? answer,
    IntentCorrections? corrections,
  }) async => _fail();

  @override
  Future<PlanningResult> generatePlans({required String tripRequestId}) async =>
      _fail();

  @override
  Future<JourneyDetail> confirmPlanOption({
    required String planOptionId,
  }) async => _fail();

  @override
  Future<JourneyDetail> getJourneyItinerary({
    required String journeyId,
  }) async => _fail();

  @override
  Future<MedicalServiceListResponse> listMedicalServices() async => _fail();
}

void main() {
  late ProviderContainer container;

  setUp(() {
    SharedPreferences.setMockInitialValues({});
    container = ProviderContainer(
      overrides: [
        journeyRepositoryProvider.overrideWithValue(
          FakeJourneyRepository(delay: Duration.zero),
        ),
      ],
    );
    addTearDown(container.dispose);
  });

  ChatController controller() =>
      container.read(chatControllerProvider.notifier);

  ChatState state() => container.read(chatControllerProvider);

  group('ChatController', () {
    test('build seeds an assistant greeting', () {
      expect(state().messages, isNotEmpty);
      expect(state().messages.first.role, ChatRole.assistant);
      expect(state().messages.first.isTyping, isFalse);
    });

    test('send appends a user bubble and surfaces the clarification', () async {
      await controller().send('I need a check-up in Batam');

      final messages = state().messages;
      expect(
        messages.any(
          (m) =>
              m.role == ChatRole.user && m.text == 'I need a check-up in Batam',
        ),
        isTrue,
      );
      expect(messages.last.text, contains('basic or comprehensive check-up'));
      // No plan cards yet.
      expect(messages.where((m) => m.planOptions != null), isEmpty);
      // The transient typing bubble is cleaned up.
      expect(messages.any((m) => m.isTyping), isFalse);
    });

    test('answering the clarification produces plan cards', () async {
      await controller().send('I need a check-up in Batam');
      await controller().send('Basic, next Friday please');

      final messages = state().messages;
      final planMessage = messages.lastWhere((m) => m.planOptions != null);
      expect(planMessage.planOptions!.length, 2);
      expect(planMessage.planOptions!.first.rank, 1);
      expect(state().tripRequest, isNotNull);
    });

    test('selecting a plan confirms the journey', () async {
      await controller().send('I need a check-up in Batam');
      await controller().send('Basic, next Friday please');
      final plan = state().messages
          .lastWhere((m) => m.planOptions != null)
          .planOptions!
          .first;

      await controller().selectPlanOption(plan);

      expect(state().journey, isNotNull);
      expect(state().journey!.journey.id, 'journey-000001');
      expect(state().messages.last.text, contains('journey is confirmed'));
    });

    test('a journey exception surfaces as an assistant error bubble', () async {
      final errorContainer = ProviderContainer(
        overrides: [
          journeyRepositoryProvider.overrideWithValue(
            const _ThrowingJourneyRepository(),
          ),
        ],
      );
      addTearDown(errorContainer.dispose);

      final notifier = errorContainer.read(chatControllerProvider.notifier);
      await notifier.send('I need a check-up in Batam');

      final messages = errorContainer.read(chatControllerProvider).messages;
      expect(
        messages.any(
          (m) =>
              m.text != null && m.text!.contains('Sorry, I couldn\'t do that'),
        ),
        isTrue,
      );
      expect(messages.any((m) => m.isTyping), isFalse);
    });
  });
}
