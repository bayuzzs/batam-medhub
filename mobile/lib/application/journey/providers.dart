import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:mobile/application/auth/providers.dart'
    show kUseFakeBackend, dioProvider;
import 'package:mobile/data/repository/fake_journey_repository.dart';
import 'package:mobile/data/repository/journey_repository.dart';
import 'package:mobile/data/repository/journey_repository_impl.dart';
import 'package:mobile/data/service/dio_journey_api.dart';
import 'package:mobile/data/service/journey_api.dart';

/// [JourneyApi] over the shared Dio client (rides the same [AuthInterceptor]
/// that attaches the Bearer token and refreshes on 401).
final journeyApiProvider = Provider<JourneyApi>((ref) {
  final dio = ref.watch(dioProvider);
  return DioJourneyApi(dio);
});

/// The journey repository. Swap here (or override this provider in tests) to
/// change between fake and real backends, matching `kUseFakeBackend` used by
/// the auth layer.
final journeyRepositoryProvider = Provider<JourneyRepository>((ref) {
  if (kUseFakeBackend) {
    return FakeJourneyRepository();
  }
  return JourneyRepositoryImpl(ref.watch(journeyApiProvider));
});
