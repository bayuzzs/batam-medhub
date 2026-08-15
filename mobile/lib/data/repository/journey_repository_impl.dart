import 'package:mobile/models/journey.dart';
import 'package:mobile/models/medical_service.dart';
import 'package:mobile/models/plan_option.dart';
import 'package:mobile/models/structured_intent.dart';

import '../service/journey_api.dart';
import 'journey_repository.dart';

/// Real [JourneyRepository]: calls the core API over [JourneyApi].
class JourneyRepositoryImpl implements JourneyRepository {
  const JourneyRepositoryImpl(this._api);

  final JourneyApi _api;

  @override
  Future<TripRequestDetail> createTripRequest({
    required String prompt,
    required String locale,
  }) {
    return _api.createTripRequest(prompt: prompt, locale: locale);
  }

  @override
  Future<TripRequestDetail> amendIntent({
    required String tripRequestId,
    String? answer,
    IntentCorrections? corrections,
  }) {
    return _api.amendIntent(
      tripRequestId: tripRequestId,
      answer: answer,
      corrections: corrections,
    );
  }

  @override
  Future<PlanningResult> generatePlans({required String tripRequestId}) {
    return _api.generatePlans(tripRequestId: tripRequestId);
  }

  @override
  Future<JourneyDetail> confirmPlanOption({required String planOptionId}) {
    return _api.confirmPlanOption(planOptionId: planOptionId);
  }

  @override
  Future<JourneyDetail> getJourneyItinerary({required String journeyId}) {
    return _api.getJourneyItinerary(journeyId: journeyId);
  }

  @override
  Future<MedicalServiceListResponse> listMedicalServices() {
    return _api.listMedicalServices();
  }
}
