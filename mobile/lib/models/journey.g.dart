// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'journey.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_Journey _$JourneyFromJson(Map<String, dynamic> json) => _Journey(
  id: json['id'] as String,
  tripRequestId: json['trip_request_id'] as String,
  status: $enumDecode(_$JourneyStatusEnumMap, json['status']),
  activeItineraryVersion: (json['active_itinerary_version'] as num).toInt(),
  createdAt: DateTime.parse(json['created_at'] as String),
  updatedAt: DateTime.parse(json['updated_at'] as String),
);

Map<String, dynamic> _$JourneyToJson(_Journey instance) => <String, dynamic>{
  'id': instance.id,
  'trip_request_id': instance.tripRequestId,
  'status': _$JourneyStatusEnumMap[instance.status]!,
  'active_itinerary_version': instance.activeItineraryVersion,
  'created_at': instance.createdAt.toIso8601String(),
  'updated_at': instance.updatedAt.toIso8601String(),
};

const _$JourneyStatusEnumMap = {
  JourneyStatus.active: 'ACTIVE',
  JourneyStatus.manualReview: 'MANUAL_REVIEW',
};

_ItineraryItem _$ItineraryItemFromJson(Map<String, dynamic> json) =>
    _ItineraryItem(
      id: json['id'] as String,
      itemType: $enumDecode(_$ItemTypeEnumMap, json['item_type']),
      providerId: json['provider_id'] as String?,
      externalReservationId: json['external_reservation_id'] as String?,
      title: json['title'] as String,
      status: $enumDecode(_$ItineraryItemStatusEnumMap, json['status']),
      timeWindow: TimeWindow.fromJson(
        json['time_window'] as Map<String, dynamic>,
      ),
      originCode: json['origin_code'] as String?,
      destinationCode: json['destination_code'] as String?,
      price: json['price'] == null
          ? null
          : ConvertedMoney.fromJson(json['price'] as Map<String, dynamic>),
      operationalNotes: (json['operational_notes'] as List<dynamic>)
          .map((e) => e as String)
          .toList(),
      synthetic: json['synthetic'] as bool,
      source: json['source'] as String,
    );

Map<String, dynamic> _$ItineraryItemToJson(_ItineraryItem instance) =>
    <String, dynamic>{
      'id': instance.id,
      'item_type': _$ItemTypeEnumMap[instance.itemType]!,
      'provider_id': instance.providerId,
      'external_reservation_id': instance.externalReservationId,
      'title': instance.title,
      'status': _$ItineraryItemStatusEnumMap[instance.status]!,
      'time_window': instance.timeWindow,
      'origin_code': instance.originCode,
      'destination_code': instance.destinationCode,
      'price': instance.price,
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

const _$ItineraryItemStatusEnumMap = {
  ItineraryItemStatus.confirmed: 'CONFIRMED',
  ItineraryItemStatus.buffer: 'BUFFER',
  ItineraryItemStatus.superseded: 'SUPERSEDED',
};

_ItineraryVersion _$ItineraryVersionFromJson(Map<String, dynamic> json) =>
    _ItineraryVersion(
      id: json['id'] as String,
      journeyId: json['journey_id'] as String,
      version: (json['version'] as num).toInt(),
      status: $enumDecode(_$ItineraryVersionStatusEnumMap, json['status']),
      basedOnDisruptionId: json['based_on_disruption_id'] as String?,
      totalPrice: PriceSummary.fromJson(
        json['total_price'] as Map<String, dynamic>,
      ),
      items: (json['items'] as List<dynamic>)
          .map((e) => ItineraryItem.fromJson(e as Map<String, dynamic>))
          .toList(),
      createdAt: DateTime.parse(json['created_at'] as String),
    );

Map<String, dynamic> _$ItineraryVersionToJson(_ItineraryVersion instance) =>
    <String, dynamic>{
      'id': instance.id,
      'journey_id': instance.journeyId,
      'version': instance.version,
      'status': _$ItineraryVersionStatusEnumMap[instance.status]!,
      'based_on_disruption_id': instance.basedOnDisruptionId,
      'total_price': instance.totalPrice,
      'items': instance.items,
      'created_at': instance.createdAt.toIso8601String(),
    };

const _$ItineraryVersionStatusEnumMap = {
  ItineraryVersionStatus.active: 'ACTIVE',
  ItineraryVersionStatus.superseded: 'SUPERSEDED',
  ItineraryVersionStatus.abandoned: 'ABANDONED',
};

_ItineraryVersionSummary _$ItineraryVersionSummaryFromJson(
  Map<String, dynamic> json,
) => _ItineraryVersionSummary(
  id: json['id'] as String,
  version: (json['version'] as num).toInt(),
  status: $enumDecode(_$ItineraryVersionStatusEnumMap, json['status']),
  basedOnDisruptionId: json['based_on_disruption_id'] as String?,
  createdAt: DateTime.parse(json['created_at'] as String),
);

Map<String, dynamic> _$ItineraryVersionSummaryToJson(
  _ItineraryVersionSummary instance,
) => <String, dynamic>{
  'id': instance.id,
  'version': instance.version,
  'status': _$ItineraryVersionStatusEnumMap[instance.status]!,
  'based_on_disruption_id': instance.basedOnDisruptionId,
  'created_at': instance.createdAt.toIso8601String(),
};

_JourneyDetail _$JourneyDetailFromJson(Map<String, dynamic> json) =>
    _JourneyDetail(
      journey: Journey.fromJson(json['journey'] as Map<String, dynamic>),
      activeItinerary: ItineraryVersion.fromJson(
        json['active_itinerary'] as Map<String, dynamic>,
      ),
      itineraryVersions: (json['itinerary_versions'] as List<dynamic>)
          .map(
            (e) => ItineraryVersionSummary.fromJson(e as Map<String, dynamic>),
          )
          .toList(),
    );

Map<String, dynamic> _$JourneyDetailToJson(_JourneyDetail instance) =>
    <String, dynamic>{
      'journey': instance.journey,
      'active_itinerary': instance.activeItinerary,
      'itinerary_versions': instance.itineraryVersions,
    };
