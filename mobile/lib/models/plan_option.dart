import 'package:freezed_annotation/freezed_annotation.dart';

import 'money.dart';
import 'time_window.dart';
import 'trip_request.dart';

part 'plan_option.freezed.dart';
part 'plan_option.g.dart';

/// Category of a leg in a planned journey (`ItemType` schema).
enum ItemType {
  @JsonValue('FERRY_OUTBOUND')
  ferryOutbound,
  @JsonValue('ARRIVAL_BUFFER')
  arrivalBuffer,
  @JsonValue('TRANSPORT_PICKUP')
  transportPickup,
  @JsonValue('HOSPITAL_APPOINTMENT')
  hospitalAppointment,
  @JsonValue('ADDITIONAL_CARE')
  additionalCare,
  @JsonValue('HOTEL_STAY')
  hotelStay,
  @JsonValue('TRANSPORT_DROPOFF')
  transportDropoff,
  @JsonValue('DEPARTURE_BUFFER')
  departureBuffer,
  @JsonValue('FERRY_RETURN')
  ferryReturn,
}

/// Lifecycle of a plan option (`PlanOptionStatus` schema).
enum PlanOptionStatus {
  @JsonValue('PROPOSED')
  proposed,
  @JsonValue('SELECTED')
  selected,
  @JsonValue('EXPIRED')
  expired,
  @JsonValue('CONFIRMED')
  confirmed,
}

/// One bookable or operational leg within a [PlanOption] (`PlanItem` schema).
/// Bookable items carry an `external_offer_id`; operational buffers don't.
@freezed
abstract class PlanItem with _$PlanItem {
  const factory PlanItem({
    required String id,
    @JsonKey(name: 'item_type') required ItemType itemType,
    @JsonKey(name: 'provider_id') String? providerId,
    @JsonKey(name: 'external_offer_id') String? externalOfferId,
    required String title,
    @JsonKey(name: 'time_window') required TimeWindow timeWindow,
    @JsonKey(name: 'origin_code') String? originCode,
    @JsonKey(name: 'destination_code') String? destinationCode,
    ConvertedMoney? price,
    @JsonKey(name: 'offer_expires_at') DateTime? offerExpiresAt,
    @JsonKey(name: 'operational_notes') required List<String> operationalNotes,
    required bool synthetic,
    required String source,
  }) = _PlanItem;

  factory PlanItem.fromJson(Map<String, dynamic> json) =>
      _$PlanItemFromJson(json);
}

/// A proposed cross-provider journey (`PlanOption` schema). At most two
/// options exist per planning revision.
@freezed
abstract class PlanOption with _$PlanOption {
  const factory PlanOption({
    required String id,
    @JsonKey(name: 'trip_request_id') required String tripRequestId,
    @JsonKey(name: 'planning_revision') required int planningRevision,
    required int rank,
    required PlanOptionStatus status,
    @JsonKey(name: 'expires_at') required DateTime expiresAt,
    required List<String> explanation,
    required List<PlanItem> items,
    @JsonKey(name: 'total_price') required PriceSummary totalPrice,
  }) = _PlanOption;

  factory PlanOption.fromJson(Map<String, dynamic> json) =>
      _$PlanOptionFromJson(json);
}

/// Result of `POST /v1/trip-requests/{id}/plans` (`PlanningResult` schema).
@freezed
abstract class PlanningResult with _$PlanningResult {
  const factory PlanningResult({
    @JsonKey(name: 'trip_request') required TripRequest tripRequest,
    required List<PlanOption> options,
    @JsonKey(name: 'no_match_reasons') required List<String> noMatchReasons,
    @JsonKey(name: 'provider_warnings') required List<String> providerWarnings,
  }) = _PlanningResult;

  factory PlanningResult.fromJson(Map<String, dynamic> json) =>
      _$PlanningResultFromJson(json);
}

/// Payload returned by create/amend trip-request endpoints
/// (`TripRequestDetail` schema): the trip request plus its current plan
/// options.
@freezed
abstract class TripRequestDetail with _$TripRequestDetail {
  const factory TripRequestDetail({
    @JsonKey(name: 'trip_request') required TripRequest tripRequest,
    @JsonKey(name: 'plan_options') required List<PlanOption> planOptions,
  }) = _TripRequestDetail;

  factory TripRequestDetail.fromJson(Map<String, dynamic> json) =>
      _$TripRequestDetailFromJson(json);
}
