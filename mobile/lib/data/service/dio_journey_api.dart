import 'dart:math';

import 'package:dio/dio.dart';

import 'package:mobile/models/journey.dart';
import 'package:mobile/models/medical_service.dart';
import 'package:mobile/models/plan_option.dart';
import 'package:mobile/models/structured_intent.dart';
import 'package:mobile/models/trip_request.dart';

import 'journey_api.dart';

/// Generates an idempotency key accepted by the core API
/// (pattern `^[A-Za-z0-9][A-Za-z0-9._:-]{7,127}$`). A fresh key per logical
/// operation lets retries replay instead of double-booking.
String _newIdempotencyKey() {
  final rng = Random.secure();
  final hex = List<int>.generate(16, (_) => rng.nextInt(256))
      .map((b) => b.toRadixString(16).padLeft(2, '0'))
      .join();
  return 'idem-$hex';
}

/// [JourneyApi] implemented over [Dio].
///
/// All endpoints are `PatientBearer`-secured, so the [AuthInterceptor] already
/// attached by the shared `dioProvider` adds the access token and transparently
/// refreshes + retries on a 401. The state-changing calls also send an
/// `Idempotency-Key` header, which the core API requires for trip-request
/// mutations and uses to make booking confirmations safe to retry.
class DioJourneyApi implements JourneyApi {
  const DioJourneyApi(this._dio);

  final Dio _dio;

  @override
  Future<TripRequestDetail> createTripRequest({
    required String prompt,
    required String locale,
  }) async {
    final response = await _dio.post<Map<String, dynamic>>(
      '/v1/trip-requests',
      data: CreateTripRequest(prompt: prompt, locale: locale).toJson(),
      options: Options(headers: {'Idempotency-Key': _newIdempotencyKey()}),
    );
    return TripRequestDetail.fromJson(response.data!);
  }

  @override
  Future<TripRequestDetail> amendIntent({
    required String tripRequestId,
    String? answer,
    IntentCorrections? corrections,
  }) async {
    final response = await _dio.patch<Map<String, dynamic>>(
      '/v1/trip-requests/$tripRequestId/intent',
      data: AmendIntentRequest(
        answer: answer,
        corrections: corrections,
      ).toJson(),
      options: Options(headers: {'Idempotency-Key': _newIdempotencyKey()}),
    );
    return TripRequestDetail.fromJson(response.data!);
  }

  @override
  Future<PlanningResult> generatePlans({required String tripRequestId}) async {
    final response = await _dio.post<Map<String, dynamic>>(
      '/v1/trip-requests/$tripRequestId/plans',
      options: Options(headers: {'Idempotency-Key': _newIdempotencyKey()}),
    );
    return PlanningResult.fromJson(response.data!);
  }

  @override
  Future<JourneyDetail> confirmPlanOption({
    required String planOptionId,
  }) async {
    final response = await _dio.post<Map<String, dynamic>>(
      '/v1/plan-options/$planOptionId/confirm',
      data: const ApprovalRequest(approved: true).toJson(),
      options: Options(headers: {'Idempotency-Key': _newIdempotencyKey()}),
    );
    return JourneyDetail.fromJson(response.data!);
  }

  @override
  Future<JourneyDetail> getJourneyItinerary({required String journeyId}) async {
    final response = await _dio.get<Map<String, dynamic>>(
      '/v1/journeys/$journeyId/itinerary',
    );
    return JourneyDetail.fromJson(response.data!);
  }

  @override
  Future<MedicalServiceListResponse> listMedicalServices() async {
    final response = await _dio.get<Map<String, dynamic>>(
      '/v1/medical-services',
    );
    return MedicalServiceListResponse.fromJson(response.data!);
  }
}

/// Body for `POST /v1/plan-options/{id}/confirm`: explicit patient approval.
/// Mirrors the `ApprovalRequest` schema (`approved` must be `true`).
class ApprovalRequest {
  const ApprovalRequest({required this.approved});

  final bool approved;

  Map<String, dynamic> toJson() => {'approved': approved};
}
