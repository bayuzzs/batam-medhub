// End-to-end integration test verifying the full journey flow between the Flutter
// mobile client and the live Go core orchestrator backend (along with 4 mock providers).

import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mobile/application/auth/auth_controller.dart';
import 'package:mobile/application/auth/providers.dart';
import 'package:mobile/application/chat/chat_controller.dart';
import 'package:mobile/application/chat/providers.dart';
import 'package:mobile/application/journey/providers.dart';
import 'package:mobile/data/repository/auth_repository_impl.dart';
import 'package:mobile/data/repository/journey_repository_impl.dart';
import 'package:mobile/data/service/auth_interceptor.dart';
import 'package:mobile/data/service/dio_auth_api.dart';
import 'package:mobile/data/service/dio_journey_api.dart';
import 'package:mobile/data/service/in_memory_token_store.dart';
import 'package:mobile/models/journey.dart';
import 'package:mobile/models/plan_option.dart';
import 'package:mobile/models/structured_intent.dart';

void main() {
  const baseUrl = 'http://localhost:8080';
  late InMemoryTokenStore tokenStore;
  late Dio unauthenticatedDio;
  late Dio authenticatedDio;
  late DioAuthApi authApi;
  late DioJourneyApi journeyApi;
  late AuthRepositoryImpl authRepo;
  late JourneyRepositoryImpl journeyRepo;
  late ProviderContainer container;

  setUpAll(() async {
    // 1. Reset demo backend state to guarantee deterministic initial DB state
    final resetDio = Dio(BaseOptions(baseUrl: baseUrl));
    final resetResp = await resetDio.post<Map<String, dynamic>>(
      '/v1/demo/reset',
      options: Options(headers: {'X-Demo-Key': 'demo_dev_secret'}),
      data: {'scenario': 'DEFAULT', 'confirm': true},
    );
    expect(resetResp.statusCode, 200);
    expect(resetResp.data!['status'], 'RESET');
  });

  setUp(() {
    tokenStore = InMemoryTokenStore();
    unauthenticatedDio = Dio(BaseOptions(baseUrl: baseUrl));
    authenticatedDio = Dio(BaseOptions(baseUrl: baseUrl));

    authApi = DioAuthApi(unauthenticatedDio);
    authRepo = AuthRepositoryImpl(api: authApi, tokenStore: tokenStore);

    authenticatedDio.interceptors.add(
      AuthInterceptor(
        dio: authenticatedDio,
        tokenStore: tokenStore,
        refresh: () async {
          final current = await tokenStore.readSession();
          if (current == null) return;
          try {
            await authRepo.refresh(refreshToken: current.refreshToken);
          } catch (_) {}
        },
      ),
    );

    journeyApi = DioJourneyApi(authenticatedDio);
    journeyRepo = JourneyRepositoryImpl(journeyApi);

    container = ProviderContainer(
      overrides: [
        tokenStoreProvider.overrideWithValue(tokenStore),
        dioProvider.overrideWithValue(authenticatedDio),
        authRepositoryProvider.overrideWithValue(authRepo),
        journeyRepositoryProvider.overrideWithValue(journeyRepo),
      ],
    );
    addTearDown(container.dispose);
  });

  group('Mobile & Backend Full-Flow Integration Tests', () {
    test('1. Patient Registration, Profile, and Token Persistence', () async {
      final session = await authRepo.register(
        fullName: 'Eleanor Vance',
        email: 'eleanor.vance.e2e@example.com',
        password: 'Password123!',
      );

      expect(session.tokenType, 'Bearer');
      expect(session.accessToken, isNotEmpty);
      expect(session.refreshToken, isNotEmpty);
      expect(session.profile.fullName, 'Eleanor Vance');
      expect(session.profile.email, 'eleanor.vance.e2e@example.com');
      expect(session.profile.preferredCurrency, 'SGD');

      // Verify token store persistence
      final storedSession = await tokenStore.readSession();
      expect(storedSession, isNotNull);
      expect(storedSession!.accessToken, session.accessToken);
    });

    test('2. Patient Login, Session Restore, and Token Rotation', () async {
      // Login with registered user
      final loginSession = await authRepo.login(
        email: 'eleanor.vance.e2e@example.com',
        password: 'Password123!',
      );

      expect(loginSession.profile.email, 'eleanor.vance.e2e@example.com');

      // Restore session
      final restored = await authRepo.restore();
      expect(restored?.accessToken, loginSession.accessToken);

      // Refresh / Rotate tokens
      final refreshed = await authRepo.refresh(
        refreshToken: loginSession.refreshToken,
      );
      expect(refreshed.accessToken, isNotEmpty);
      expect(refreshed.refreshToken, isNot(loginSession.refreshToken));

      final updatedSession = await tokenStore.readSession();
      expect(updatedSession?.accessToken, refreshed.accessToken);
    });

    test('3. Medical Catalog Retrieval', () async {
      // Log in to set token
      await authRepo.login(
        email: 'eleanor.vance.e2e@example.com',
        password: 'Password123!',
      );

      final catalog = await journeyApi.listMedicalServices();
      expect(catalog.services, isNotEmpty);
      final codes = catalog.services.map((item) => item.code).toList();
      expect(codes, contains('MCU_COMPREHENSIVE'));
      expect(codes, contains('MCU_BASIC'));
    });

    test('4. Full Journey Flow: Intent -> Plans -> Saga Booking -> Itinerary', () async {
      // Step A: Ensure logged in
      final session = await authRepo.login(
        email: 'eleanor.vance.e2e@example.com',
        password: 'Password123!',
      );
      expect(session.accessToken, isNotEmpty);

      // Step B: Submit natural language prompt
      final tripDetail = await journeyRepo.createTripRequest(
        prompt: 'I need a comprehensive health screening in Batam on 22 August 2026 for 1 person.',
        locale: 'en',
      );

      expect(tripDetail.tripRequest.id, isNotEmpty);
      expect(
        tripDetail.tripRequest.intent.resolution,
        IntentResolution.matched,
      );
      expect(
        tripDetail.tripRequest.intent.serviceCode,
        'MCU_COMPREHENSIVE',
      );

      // Step C: Generate Ranked Plan Options from 4 Mock Providers
      final planningResult = await journeyRepo.generatePlans(
        tripRequestId: tripDetail.tripRequest.id,
      );

      expect(planningResult.options, isNotEmpty);
      final rank1Option = planningResult.options.firstWhere((o) => o.rank == 1);
      expect(rank1Option.items, isNotEmpty);
      expect(rank1Option.totalPrice.displayTotal.currency, 'SGD');

      // Verify itemized legs
      final itemTypes = rank1Option.items.map((i) => i.itemType).toList();
      expect(itemTypes, contains(ItemType.ferryOutbound));
      expect(itemTypes, contains(ItemType.hospitalAppointment));
      expect(itemTypes, contains(ItemType.transportPickup));

      // Step D: Execute Distributed Booking Saga
      final journeyDetail = await journeyRepo.confirmPlanOption(
        planOptionId: rank1Option.id,
      );

      expect(journeyDetail.journey.status, JourneyStatus.active);
      expect(journeyDetail.journey.activeItineraryVersion, 1);
      expect(journeyDetail.activeItinerary.items.length, 6);

      // All 6 legs should be confirmed
      for (final item in journeyDetail.activeItinerary.items) {
        expect(item.status, ItineraryItemStatus.confirmed);
      }

      // Step E: Fetch Active Itinerary by Journey ID
      final fetchedJourney = await journeyRepo.getJourneyItinerary(
        journeyId: journeyDetail.journey.id,
      );
      expect(fetchedJourney.journey.id, journeyDetail.journey.id);
      expect(fetchedJourney.activeItinerary.version, 1);
    });

    test('5. Disruption Ingestion & Versioned Itinerary Recovery', () async {
      // Step A: Login & create a journey to disrupt
      await authRepo.login(
        email: 'eleanor.vance.e2e@example.com',
        password: 'Password123!',
      );

      final tripDetail = await journeyRepo.createTripRequest(
        prompt: 'I need a comprehensive health screening in Batam on 22 August 2026 for 1 person.',
        locale: 'en',
      );

      final planningResult = await journeyRepo.generatePlans(
        tripRequestId: tripDetail.tripRequest.id,
      );
      final planOption = planningResult.options.first;

      final journeyDetail = await journeyRepo.confirmPlanOption(
        planOptionId: planOption.id,
      );
      final journeyId = journeyDetail.journey.id;
      final hospitalItem = journeyDetail.activeItinerary.items.firstWhere(
        (i) => i.itemType == ItemType.hospitalAppointment,
      );

      // Step B: Simulate Hospital Provider Disruption Event via HTTP Webhook
      final disruptionPayload = {
        'external_event_id': 'evt-hosp-mobile-e2e-${DateTime.now().millisecondsSinceEpoch}',
        'journey_id': journeyId,
        'event_type': 'HOSPITAL_ADDITIONAL_CARE_REQUESTED',
        'occurred_at': '2026-08-22T04:30:00Z',
        'target': {
          'itinerary_item_id': hospitalItem.id,
        },
        'actor': {
          'actor_id': 'dr-lee-tan',
          'name': 'Dr Lee Tan',
          'role': 'Cardiologist',
        },
        'details': {
          'reason': 'Patient requires additional cardiac observation following exam.',
          'instruction_reference': 'hospital-instruction://followup-observation/FO-20260822-0001',
          'replacement_time_window': {
            'starts_at': '2026-08-22T05:00:00Z',
            'ends_at': '2026-08-22T06:30:00Z',
            'start_time_zone': 'Asia/Jakarta',
            'end_time_zone': 'Asia/Jakarta',
          },
          'additional_service_code': 'FOLLOWUP_OBSERVATION',
          'additional_duration_minutes': 90,
          'priority': 'MEDIUM',
          'travel_clearance_status': 'CLEARED',
        },
      };

      final disruptionResp = await unauthenticatedDio.post<Map<String, dynamic>>(
        '/v1/provider/disruptions',
        options: Options(headers: {'X-Provider-Key': 'hospital_dev_secret'}),
        data: disruptionPayload,
      );
      expect(disruptionResp.statusCode, 202);
      final disruptionId = disruptionResp.data!['disruption_id'] as String;
      expect(disruptionId, isNotEmpty);

      // Step C: Patient retrieves disruption details & recovery options
      final disruptionDetailResp = await authenticatedDio.get<Map<String, dynamic>>(
        '/v1/disruptions/$disruptionId',
      );
      expect(disruptionDetailResp.statusCode, 200);
      final recoveryOptions = disruptionDetailResp.data!['recovery_options'] as List<dynamic>;
      expect(recoveryOptions, isNotEmpty);
      final recoveryOptionId = recoveryOptions.first['id'] as String;

      // Step D: Patient approves recovery option -> Activates Itinerary Version 2
      final approveResp = await authenticatedDio.post<Map<String, dynamic>>(
        '/v1/recovery-options/$recoveryOptionId/approve',
        data: {'approved': true},
        options: Options(headers: {'Idempotency-Key': 'idem-approve-mobile-${DateTime.now().millisecondsSinceEpoch}'}),
      );
      expect(approveResp.statusCode, 200);
      final updatedJourney = JourneyDetail.fromJson(approveResp.data!);

      expect(updatedJourney.journey.activeItineraryVersion, 2);
      expect(updatedJourney.activeItinerary.version, 2);
      expect(updatedJourney.activeItinerary.items.length, 8);
      expect(updatedJourney.itineraryVersions.length, 1);
      expect(updatedJourney.itineraryVersions.first.version, 1);
      expect(updatedJourney.itineraryVersions.first.status, ItineraryVersionStatus.superseded);
    });

    test('6. ChatController End-to-End Orchestration with Real Backend', () async {
      // Step A: Login through AuthController
      final authCtrl = container.read(authControllerProvider.notifier);
      await authCtrl.login(
        email: 'eleanor.vance.e2e@example.com',
        password: 'Password123!',
      );

      // Step B: Drive chat conversation with ChatController
      final chatCtrl = container.read(chatControllerProvider.notifier);
      await chatCtrl.send('I need a comprehensive health screening in Batam on 22 August 2026 for 1 person.');

      final state = container.read(chatControllerProvider);
      expect(state.tripRequest, isNotNull);
      expect(state.tripRequest!.intent.resolution, IntentResolution.matched);

      final planOptionMsg = state.messages.lastWhere((m) => m.planOptions != null);
      expect(planOptionMsg.planOptions, isNotEmpty);

      // Step C: Confirm plan option via ChatController
      final chosenPlan = planOptionMsg.planOptions!.first;
      await chatCtrl.selectPlanOption(chosenPlan);

      final finalState = container.read(chatControllerProvider);
      expect(finalState.journey, isNotNull);
      expect(finalState.journey!.journey.status, JourneyStatus.active);
      expect(finalState.messages.last.text, contains('journey is confirmed'));
    });

    test('7. Patient Logout and Session Cleanup', () async {
      final authCtrl = container.read(authControllerProvider.notifier);
      await authCtrl.login(
        email: 'eleanor.vance.e2e@example.com',
        password: 'Password123!',
      );

      expect(await tokenStore.readSession(), isNotNull);

      await authCtrl.logout();
      expect(await tokenStore.readSession(), isNull);
    });
  });
}
