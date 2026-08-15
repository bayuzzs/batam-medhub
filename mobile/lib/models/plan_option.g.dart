// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'plan_option.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_PlanItem _$PlanItemFromJson(Map<String, dynamic> json) => _PlanItem(
  id: json['id'] as String,
  itemType: $enumDecode(_$ItemTypeEnumMap, json['item_type']),
  providerId: json['provider_id'] as String?,
  externalOfferId: json['external_offer_id'] as String?,
  title: json['title'] as String,
  timeWindow: TimeWindow.fromJson(json['time_window'] as Map<String, dynamic>),
  originCode: json['origin_code'] as String?,
  destinationCode: json['destination_code'] as String?,
  price: json['price'] == null
      ? null
      : ConvertedMoney.fromJson(json['price'] as Map<String, dynamic>),
  offerExpiresAt: json['offer_expires_at'] == null
      ? null
      : DateTime.parse(json['offer_expires_at'] as String),
  operationalNotes: (json['operational_notes'] as List<dynamic>)
      .map((e) => e as String)
      .toList(),
  synthetic: json['synthetic'] as bool,
  source: json['source'] as String,
);

Map<String, dynamic> _$PlanItemToJson(_PlanItem instance) => <String, dynamic>{
  'id': instance.id,
  'item_type': _$ItemTypeEnumMap[instance.itemType]!,
  'provider_id': instance.providerId,
  'external_offer_id': instance.externalOfferId,
  'title': instance.title,
  'time_window': instance.timeWindow,
  'origin_code': instance.originCode,
  'destination_code': instance.destinationCode,
  'price': instance.price,
  'offer_expires_at': instance.offerExpiresAt?.toIso8601String(),
  'operational_notes': instance.operationalNotes,
  'synthetic': instance.synthetic,
  'source': instance.source,
};

const _$ItemTypeEnumMap = {
  ItemType.ferryOutbound: 'FERRY_OUTBOUND',
  ItemType.arrivalBuffer: 'ARRIVAL_BUFFER',
  ItemType.transportPickup: 'TRANSPORT_PICKUP',
  ItemType.hospitalAppointment: 'HOSPITAL_APPOINTMENT',
  ItemType.additionalCare: 'ADDITIONAL_CARE',
  ItemType.hotelStay: 'HOTEL_STAY',
  ItemType.transportDropoff: 'TRANSPORT_DROPOFF',
  ItemType.departureBuffer: 'DEPARTURE_BUFFER',
  ItemType.ferryReturn: 'FERRY_RETURN',
};

_PlanOption _$PlanOptionFromJson(Map<String, dynamic> json) => _PlanOption(
  id: json['id'] as String,
  tripRequestId: json['trip_request_id'] as String,
  planningRevision: (json['planning_revision'] as num).toInt(),
  rank: (json['rank'] as num).toInt(),
  status: $enumDecode(_$PlanOptionStatusEnumMap, json['status']),
  expiresAt: DateTime.parse(json['expires_at'] as String),
  explanation: (json['explanation'] as List<dynamic>)
      .map((e) => e as String)
      .toList(),
  items: (json['items'] as List<dynamic>)
      .map((e) => PlanItem.fromJson(e as Map<String, dynamic>))
      .toList(),
  totalPrice: PriceSummary.fromJson(
    json['total_price'] as Map<String, dynamic>,
  ),
);

Map<String, dynamic> _$PlanOptionToJson(_PlanOption instance) =>
    <String, dynamic>{
      'id': instance.id,
      'trip_request_id': instance.tripRequestId,
      'planning_revision': instance.planningRevision,
      'rank': instance.rank,
      'status': _$PlanOptionStatusEnumMap[instance.status]!,
      'expires_at': instance.expiresAt.toIso8601String(),
      'explanation': instance.explanation,
      'items': instance.items,
      'total_price': instance.totalPrice,
    };

const _$PlanOptionStatusEnumMap = {
  PlanOptionStatus.proposed: 'PROPOSED',
  PlanOptionStatus.selected: 'SELECTED',
  PlanOptionStatus.expired: 'EXPIRED',
  PlanOptionStatus.confirmed: 'CONFIRMED',
};

_PlanningResult _$PlanningResultFromJson(Map<String, dynamic> json) =>
    _PlanningResult(
      tripRequest: TripRequest.fromJson(
        json['trip_request'] as Map<String, dynamic>,
      ),
      options: (json['options'] as List<dynamic>)
          .map((e) => PlanOption.fromJson(e as Map<String, dynamic>))
          .toList(),
      noMatchReasons: (json['no_match_reasons'] as List<dynamic>)
          .map((e) => e as String)
          .toList(),
      providerWarnings: (json['provider_warnings'] as List<dynamic>)
          .map((e) => e as String)
          .toList(),
    );

Map<String, dynamic> _$PlanningResultToJson(_PlanningResult instance) =>
    <String, dynamic>{
      'trip_request': instance.tripRequest,
      'options': instance.options,
      'no_match_reasons': instance.noMatchReasons,
      'provider_warnings': instance.providerWarnings,
    };

_TripRequestDetail _$TripRequestDetailFromJson(Map<String, dynamic> json) =>
    _TripRequestDetail(
      tripRequest: TripRequest.fromJson(
        json['trip_request'] as Map<String, dynamic>,
      ),
      planOptions: (json['plan_options'] as List<dynamic>)
          .map((e) => PlanOption.fromJson(e as Map<String, dynamic>))
          .toList(),
    );

Map<String, dynamic> _$TripRequestDetailToJson(_TripRequestDetail instance) =>
    <String, dynamic>{
      'trip_request': instance.tripRequest,
      'plan_options': instance.planOptions,
    };
