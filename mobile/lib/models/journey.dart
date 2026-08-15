import 'package:freezed_annotation/freezed_annotation.dart';

import 'money.dart';
import 'plan_option.dart';
import 'time_window.dart';

part 'journey.freezed.dart';
part 'journey.g.dart';

/// Lifecycle of a confirmed journey (`JourneyStatus` schema).
enum JourneyStatus {
  @JsonValue('ACTIVE')
  active,
  @JsonValue('MANUAL_REVIEW')
  manualReview,
}

/// Booking status of an itinerary item.
enum ItineraryItemStatus {
  @JsonValue('CONFIRMED')
  confirmed,
  @JsonValue('BUFFER')
  buffer,
  @JsonValue('SUPERSEDED')
  superseded,
}

/// Version status of an itinerary (`ItineraryVersionStatus` schema).
enum ItineraryVersionStatus {
  @JsonValue('ACTIVE')
  active,
  @JsonValue('SUPERSEDED')
  superseded,
  @JsonValue('ABANDONED')
  abandoned,
}

/// A confirmed journey owned by the patient (`Journey` schema).
@freezed
abstract class Journey with _$Journey {
  const factory Journey({
    required String id,
    @JsonKey(name: 'trip_request_id') required String tripRequestId,
    required JourneyStatus status,
    @JsonKey(name: 'active_itinerary_version')
    required int activeItineraryVersion,
    @JsonKey(name: 'created_at') required DateTime createdAt,
    @JsonKey(name: 'updated_at') required DateTime updatedAt,
  }) = _Journey;

  factory Journey.fromJson(Map<String, dynamic> json) =>
      _$JourneyFromJson(json);
}

/// A booked leg within an itinerary version (`ItineraryItem` schema).
@freezed
abstract class ItineraryItem with _$ItineraryItem {
  const factory ItineraryItem({
    required String id,
    @JsonKey(name: 'item_type') required ItemType itemType,
    @JsonKey(name: 'provider_id') String? providerId,
    @JsonKey(name: 'external_reservation_id') String? externalReservationId,
    required String title,
    required ItineraryItemStatus status,
    @JsonKey(name: 'time_window') required TimeWindow timeWindow,
    @JsonKey(name: 'origin_code') String? originCode,
    @JsonKey(name: 'destination_code') String? destinationCode,
    ConvertedMoney? price,
    @JsonKey(name: 'operational_notes') required List<String> operationalNotes,
    required bool synthetic,
    required String source,
  }) = _ItineraryItem;

  factory ItineraryItem.fromJson(Map<String, dynamic> json) =>
      _$ItineraryItemFromJson(json);
}

/// An immutable version of a confirmed itinerary (`ItineraryVersion` schema).
/// Activated versions are append-only; recovery creates a successor version.
@freezed
abstract class ItineraryVersion with _$ItineraryVersion {
  const factory ItineraryVersion({
    required String id,
    @JsonKey(name: 'journey_id') required String journeyId,
    required int version,
    required ItineraryVersionStatus status,
    @JsonKey(name: 'based_on_disruption_id') String? basedOnDisruptionId,
    @JsonKey(name: 'total_price') required PriceSummary totalPrice,
    required List<ItineraryItem> items,
    @JsonKey(name: 'created_at') required DateTime createdAt,
  }) = _ItineraryVersion;

  factory ItineraryVersion.fromJson(Map<String, dynamic> json) =>
      _$ItineraryVersionFromJson(json);
}

/// Lightweight reference to a non-active itinerary version
/// (`ItineraryVersionSummary` schema).
@freezed
abstract class ItineraryVersionSummary with _$ItineraryVersionSummary {
  const factory ItineraryVersionSummary({
    required String id,
    required int version,
    required ItineraryVersionStatus status,
    @JsonKey(name: 'based_on_disruption_id') String? basedOnDisruptionId,
    @JsonKey(name: 'created_at') required DateTime createdAt,
  }) = _ItineraryVersionSummary;

  factory ItineraryVersionSummary.fromJson(Map<String, dynamic> json) =>
      _$ItineraryVersionSummaryFromJson(json);
}

/// A journey with its active itinerary and version history
/// (`JourneyDetail` schema) — returned by confirm and itinerary endpoints.
@freezed
abstract class JourneyDetail with _$JourneyDetail {
  const factory JourneyDetail({
    required Journey journey,
    @JsonKey(name: 'active_itinerary')
    required ItineraryVersion activeItinerary,
    @JsonKey(name: 'itinerary_versions')
    required List<ItineraryVersionSummary> itineraryVersions,
  }) = _JourneyDetail;

  factory JourneyDetail.fromJson(Map<String, dynamic> json) =>
      _$JourneyDetailFromJson(json);
}
