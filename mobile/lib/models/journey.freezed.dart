// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint, type=warning, deprecated_member_use, deprecated_member_use_from_same_package
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'journey.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// dart format off
T _$identity<T>(T value) => value;

/// @nodoc
mixin _$Journey {

 String get id;@JsonKey(name: 'trip_request_id') String get tripRequestId; JourneyStatus get status;@JsonKey(name: 'active_itinerary_version') int get activeItineraryVersion;@JsonKey(name: 'created_at') DateTime get createdAt;@JsonKey(name: 'updated_at') DateTime get updatedAt;
/// Create a copy of Journey
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$JourneyCopyWith<Journey> get copyWith => _$JourneyCopyWithImpl<Journey>(this as Journey, _$identity);

  /// Serializes this Journey to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Journey&&(identical(other.id, id) || other.id == id)&&(identical(other.tripRequestId, tripRequestId) || other.tripRequestId == tripRequestId)&&(identical(other.status, status) || other.status == status)&&(identical(other.activeItineraryVersion, activeItineraryVersion) || other.activeItineraryVersion == activeItineraryVersion)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt)&&(identical(other.updatedAt, updatedAt) || other.updatedAt == updatedAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,tripRequestId,status,activeItineraryVersion,createdAt,updatedAt);

@override
String toString() {
  return 'Journey(id: $id, tripRequestId: $tripRequestId, status: $status, activeItineraryVersion: $activeItineraryVersion, createdAt: $createdAt, updatedAt: $updatedAt)';
}


}

/// @nodoc
abstract mixin class $JourneyCopyWith<$Res>  {
  factory $JourneyCopyWith(Journey value, $Res Function(Journey) _then) = _$JourneyCopyWithImpl;
@useResult
$Res call({
 String id,@JsonKey(name: 'trip_request_id') String tripRequestId, JourneyStatus status,@JsonKey(name: 'active_itinerary_version') int activeItineraryVersion,@JsonKey(name: 'created_at') DateTime createdAt,@JsonKey(name: 'updated_at') DateTime updatedAt
});




}
/// @nodoc
class _$JourneyCopyWithImpl<$Res>
    implements $JourneyCopyWith<$Res> {
  _$JourneyCopyWithImpl(this._self, this._then);

  final Journey _self;
  final $Res Function(Journey) _then;

/// Create a copy of Journey
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? tripRequestId = null,Object? status = null,Object? activeItineraryVersion = null,Object? createdAt = null,Object? updatedAt = null,}) {
  return _then(Journey(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,tripRequestId: null == tripRequestId ? _self.tripRequestId : tripRequestId // ignore: cast_nullable_to_non_nullable
as String,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as JourneyStatus,activeItineraryVersion: null == activeItineraryVersion ? _self.activeItineraryVersion : activeItineraryVersion // ignore: cast_nullable_to_non_nullable
as int,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,updatedAt: null == updatedAt ? _self.updatedAt : updatedAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}

}


/// Adds pattern-matching-related methods to [Journey].
extension JourneyPatterns on Journey {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Journey value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Journey() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Journey value)  $default,){
final _that = this;
switch (_that) {
case _Journey():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Journey value)?  $default,){
final _that = this;
switch (_that) {
case _Journey() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id, @JsonKey(name: 'trip_request_id')  String tripRequestId,  JourneyStatus status, @JsonKey(name: 'active_itinerary_version')  int activeItineraryVersion, @JsonKey(name: 'created_at')  DateTime createdAt, @JsonKey(name: 'updated_at')  DateTime updatedAt)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Journey() when $default != null:
return $default(_that.id,_that.tripRequestId,_that.status,_that.activeItineraryVersion,_that.createdAt,_that.updatedAt);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id, @JsonKey(name: 'trip_request_id')  String tripRequestId,  JourneyStatus status, @JsonKey(name: 'active_itinerary_version')  int activeItineraryVersion, @JsonKey(name: 'created_at')  DateTime createdAt, @JsonKey(name: 'updated_at')  DateTime updatedAt)  $default,) {final _that = this;
switch (_that) {
case _Journey():
return $default(_that.id,_that.tripRequestId,_that.status,_that.activeItineraryVersion,_that.createdAt,_that.updatedAt);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id, @JsonKey(name: 'trip_request_id')  String tripRequestId,  JourneyStatus status, @JsonKey(name: 'active_itinerary_version')  int activeItineraryVersion, @JsonKey(name: 'created_at')  DateTime createdAt, @JsonKey(name: 'updated_at')  DateTime updatedAt)?  $default,) {final _that = this;
switch (_that) {
case _Journey() when $default != null:
return $default(_that.id,_that.tripRequestId,_that.status,_that.activeItineraryVersion,_that.createdAt,_that.updatedAt);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _Journey implements Journey {
  const _Journey({required this.id, @JsonKey(name: 'trip_request_id') required this.tripRequestId, required this.status, @JsonKey(name: 'active_itinerary_version') required this.activeItineraryVersion, @JsonKey(name: 'created_at') required this.createdAt, @JsonKey(name: 'updated_at') required this.updatedAt});
  factory _Journey.fromJson(Map<String, dynamic> json) => _$JourneyFromJson(json);

@override final  String id;
@override@JsonKey(name: 'trip_request_id') final  String tripRequestId;
@override final  JourneyStatus status;
@override@JsonKey(name: 'active_itinerary_version') final  int activeItineraryVersion;
@override@JsonKey(name: 'created_at') final  DateTime createdAt;
@override@JsonKey(name: 'updated_at') final  DateTime updatedAt;

/// Create a copy of Journey
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$JourneyCopyWith<_Journey> get copyWith => __$JourneyCopyWithImpl<_Journey>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$JourneyToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _Journey&&(identical(other.id, id) || other.id == id)&&(identical(other.tripRequestId, tripRequestId) || other.tripRequestId == tripRequestId)&&(identical(other.status, status) || other.status == status)&&(identical(other.activeItineraryVersion, activeItineraryVersion) || other.activeItineraryVersion == activeItineraryVersion)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt)&&(identical(other.updatedAt, updatedAt) || other.updatedAt == updatedAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,tripRequestId,status,activeItineraryVersion,createdAt,updatedAt);

@override
String toString() {
  return 'Journey(id: $id, tripRequestId: $tripRequestId, status: $status, activeItineraryVersion: $activeItineraryVersion, createdAt: $createdAt, updatedAt: $updatedAt)';
}


}

/// @nodoc
abstract mixin class _$JourneyCopyWith<$Res> implements $JourneyCopyWith<$Res> {
  factory _$JourneyCopyWith(_Journey value, $Res Function(_Journey) _then) = __$JourneyCopyWithImpl;
@override @useResult
$Res call({
 String id,@JsonKey(name: 'trip_request_id') String tripRequestId, JourneyStatus status,@JsonKey(name: 'active_itinerary_version') int activeItineraryVersion,@JsonKey(name: 'created_at') DateTime createdAt,@JsonKey(name: 'updated_at') DateTime updatedAt
});




}
/// @nodoc
class __$JourneyCopyWithImpl<$Res>
    implements _$JourneyCopyWith<$Res> {
  __$JourneyCopyWithImpl(this._self, this._then);

  final _Journey _self;
  final $Res Function(_Journey) _then;

/// Create a copy of Journey
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? tripRequestId = null,Object? status = null,Object? activeItineraryVersion = null,Object? createdAt = null,Object? updatedAt = null,}) {
  return _then(_Journey(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,tripRequestId: null == tripRequestId ? _self.tripRequestId : tripRequestId // ignore: cast_nullable_to_non_nullable
as String,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as JourneyStatus,activeItineraryVersion: null == activeItineraryVersion ? _self.activeItineraryVersion : activeItineraryVersion // ignore: cast_nullable_to_non_nullable
as int,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,updatedAt: null == updatedAt ? _self.updatedAt : updatedAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}


}


/// @nodoc
mixin _$ItineraryItem {

 String get id;@JsonKey(name: 'item_type') ItemType get itemType;@JsonKey(name: 'provider_id') String? get providerId;@JsonKey(name: 'external_reservation_id') String? get externalReservationId; String get title; ItineraryItemStatus get status;@JsonKey(name: 'time_window') TimeWindow get timeWindow;@JsonKey(name: 'origin_code') String? get originCode;@JsonKey(name: 'destination_code') String? get destinationCode; ConvertedMoney? get price;@JsonKey(name: 'operational_notes') List<String> get operationalNotes; bool get synthetic; String get source;
/// Create a copy of ItineraryItem
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$ItineraryItemCopyWith<ItineraryItem> get copyWith => _$ItineraryItemCopyWithImpl<ItineraryItem>(this as ItineraryItem, _$identity);

  /// Serializes this ItineraryItem to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is ItineraryItem&&(identical(other.id, id) || other.id == id)&&(identical(other.itemType, itemType) || other.itemType == itemType)&&(identical(other.providerId, providerId) || other.providerId == providerId)&&(identical(other.externalReservationId, externalReservationId) || other.externalReservationId == externalReservationId)&&(identical(other.title, title) || other.title == title)&&(identical(other.status, status) || other.status == status)&&(identical(other.timeWindow, timeWindow) || other.timeWindow == timeWindow)&&(identical(other.originCode, originCode) || other.originCode == originCode)&&(identical(other.destinationCode, destinationCode) || other.destinationCode == destinationCode)&&(identical(other.price, price) || other.price == price)&&const DeepCollectionEquality().equals(other.operationalNotes, operationalNotes)&&(identical(other.synthetic, synthetic) || other.synthetic == synthetic)&&(identical(other.source, source) || other.source == source));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,itemType,providerId,externalReservationId,title,status,timeWindow,originCode,destinationCode,price,const DeepCollectionEquality().hash(operationalNotes),synthetic,source);

@override
String toString() {
  return 'ItineraryItem(id: $id, itemType: $itemType, providerId: $providerId, externalReservationId: $externalReservationId, title: $title, status: $status, timeWindow: $timeWindow, originCode: $originCode, destinationCode: $destinationCode, price: $price, operationalNotes: $operationalNotes, synthetic: $synthetic, source: $source)';
}


}

/// @nodoc
abstract mixin class $ItineraryItemCopyWith<$Res>  {
  factory $ItineraryItemCopyWith(ItineraryItem value, $Res Function(ItineraryItem) _then) = _$ItineraryItemCopyWithImpl;
@useResult
$Res call({
 String id,@JsonKey(name: 'item_type') ItemType itemType,@JsonKey(name: 'provider_id') String? providerId,@JsonKey(name: 'external_reservation_id') String? externalReservationId, String title, ItineraryItemStatus status,@JsonKey(name: 'time_window') TimeWindow timeWindow,@JsonKey(name: 'origin_code') String? originCode,@JsonKey(name: 'destination_code') String? destinationCode, ConvertedMoney? price,@JsonKey(name: 'operational_notes') List<String> operationalNotes, bool synthetic, String source
});


$TimeWindowCopyWith<$Res> get timeWindow;$ConvertedMoneyCopyWith<$Res>? get price;

}
/// @nodoc
class _$ItineraryItemCopyWithImpl<$Res>
    implements $ItineraryItemCopyWith<$Res> {
  _$ItineraryItemCopyWithImpl(this._self, this._then);

  final ItineraryItem _self;
  final $Res Function(ItineraryItem) _then;

/// Create a copy of ItineraryItem
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? itemType = null,Object? providerId = freezed,Object? externalReservationId = freezed,Object? title = null,Object? status = null,Object? timeWindow = null,Object? originCode = freezed,Object? destinationCode = freezed,Object? price = freezed,Object? operationalNotes = null,Object? synthetic = null,Object? source = null,}) {
  return _then(ItineraryItem(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,itemType: null == itemType ? _self.itemType : itemType // ignore: cast_nullable_to_non_nullable
as ItemType,providerId: freezed == providerId ? _self.providerId : providerId // ignore: cast_nullable_to_non_nullable
as String?,externalReservationId: freezed == externalReservationId ? _self.externalReservationId : externalReservationId // ignore: cast_nullable_to_non_nullable
as String?,title: null == title ? _self.title : title // ignore: cast_nullable_to_non_nullable
as String,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as ItineraryItemStatus,timeWindow: null == timeWindow ? _self.timeWindow : timeWindow // ignore: cast_nullable_to_non_nullable
as TimeWindow,originCode: freezed == originCode ? _self.originCode : originCode // ignore: cast_nullable_to_non_nullable
as String?,destinationCode: freezed == destinationCode ? _self.destinationCode : destinationCode // ignore: cast_nullable_to_non_nullable
as String?,price: freezed == price ? _self.price : price // ignore: cast_nullable_to_non_nullable
as ConvertedMoney?,operationalNotes: null == operationalNotes ? _self.operationalNotes : operationalNotes // ignore: cast_nullable_to_non_nullable
as List<String>,synthetic: null == synthetic ? _self.synthetic : synthetic // ignore: cast_nullable_to_non_nullable
as bool,source: null == source ? _self.source : source // ignore: cast_nullable_to_non_nullable
as String,
  ));
}
/// Create a copy of ItineraryItem
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$TimeWindowCopyWith<$Res> get timeWindow {
  
  return $TimeWindowCopyWith<$Res>(_self.timeWindow, (value) {
    return _then(_self.copyWith(timeWindow: value));
  });
}/// Create a copy of ItineraryItem
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$ConvertedMoneyCopyWith<$Res>? get price {
    if (_self.price == null) {
    return null;
  }

  return $ConvertedMoneyCopyWith<$Res>(_self.price!, (value) {
    return _then(_self.copyWith(price: value));
  });
}
}


/// Adds pattern-matching-related methods to [ItineraryItem].
extension ItineraryItemPatterns on ItineraryItem {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _ItineraryItem value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _ItineraryItem() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _ItineraryItem value)  $default,){
final _that = this;
switch (_that) {
case _ItineraryItem():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _ItineraryItem value)?  $default,){
final _that = this;
switch (_that) {
case _ItineraryItem() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id, @JsonKey(name: 'item_type')  ItemType itemType, @JsonKey(name: 'provider_id')  String? providerId, @JsonKey(name: 'external_reservation_id')  String? externalReservationId,  String title,  ItineraryItemStatus status, @JsonKey(name: 'time_window')  TimeWindow timeWindow, @JsonKey(name: 'origin_code')  String? originCode, @JsonKey(name: 'destination_code')  String? destinationCode,  ConvertedMoney? price, @JsonKey(name: 'operational_notes')  List<String> operationalNotes,  bool synthetic,  String source)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _ItineraryItem() when $default != null:
return $default(_that.id,_that.itemType,_that.providerId,_that.externalReservationId,_that.title,_that.status,_that.timeWindow,_that.originCode,_that.destinationCode,_that.price,_that.operationalNotes,_that.synthetic,_that.source);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id, @JsonKey(name: 'item_type')  ItemType itemType, @JsonKey(name: 'provider_id')  String? providerId, @JsonKey(name: 'external_reservation_id')  String? externalReservationId,  String title,  ItineraryItemStatus status, @JsonKey(name: 'time_window')  TimeWindow timeWindow, @JsonKey(name: 'origin_code')  String? originCode, @JsonKey(name: 'destination_code')  String? destinationCode,  ConvertedMoney? price, @JsonKey(name: 'operational_notes')  List<String> operationalNotes,  bool synthetic,  String source)  $default,) {final _that = this;
switch (_that) {
case _ItineraryItem():
return $default(_that.id,_that.itemType,_that.providerId,_that.externalReservationId,_that.title,_that.status,_that.timeWindow,_that.originCode,_that.destinationCode,_that.price,_that.operationalNotes,_that.synthetic,_that.source);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id, @JsonKey(name: 'item_type')  ItemType itemType, @JsonKey(name: 'provider_id')  String? providerId, @JsonKey(name: 'external_reservation_id')  String? externalReservationId,  String title,  ItineraryItemStatus status, @JsonKey(name: 'time_window')  TimeWindow timeWindow, @JsonKey(name: 'origin_code')  String? originCode, @JsonKey(name: 'destination_code')  String? destinationCode,  ConvertedMoney? price, @JsonKey(name: 'operational_notes')  List<String> operationalNotes,  bool synthetic,  String source)?  $default,) {final _that = this;
switch (_that) {
case _ItineraryItem() when $default != null:
return $default(_that.id,_that.itemType,_that.providerId,_that.externalReservationId,_that.title,_that.status,_that.timeWindow,_that.originCode,_that.destinationCode,_that.price,_that.operationalNotes,_that.synthetic,_that.source);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _ItineraryItem implements ItineraryItem {
  const _ItineraryItem({required this.id, @JsonKey(name: 'item_type') required this.itemType, @JsonKey(name: 'provider_id') this.providerId, @JsonKey(name: 'external_reservation_id') this.externalReservationId, required this.title, required this.status, @JsonKey(name: 'time_window') required this.timeWindow, @JsonKey(name: 'origin_code') this.originCode, @JsonKey(name: 'destination_code') this.destinationCode, this.price, @JsonKey(name: 'operational_notes') required  List<String> operationalNotes, required this.synthetic, required this.source}): _operationalNotes = operationalNotes;
  factory _ItineraryItem.fromJson(Map<String, dynamic> json) => _$ItineraryItemFromJson(json);

@override final  String id;
@override@JsonKey(name: 'item_type') final  ItemType itemType;
@override@JsonKey(name: 'provider_id') final  String? providerId;
@override@JsonKey(name: 'external_reservation_id') final  String? externalReservationId;
@override final  String title;
@override final  ItineraryItemStatus status;
@override@JsonKey(name: 'time_window') final  TimeWindow timeWindow;
@override@JsonKey(name: 'origin_code') final  String? originCode;
@override@JsonKey(name: 'destination_code') final  String? destinationCode;
@override final  ConvertedMoney? price;
 final  List<String> _operationalNotes;
@override@JsonKey(name: 'operational_notes') List<String> get operationalNotes {
  if (_operationalNotes is EqualUnmodifiableListView) return _operationalNotes;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_operationalNotes);
}

@override final  bool synthetic;
@override final  String source;

/// Create a copy of ItineraryItem
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$ItineraryItemCopyWith<_ItineraryItem> get copyWith => __$ItineraryItemCopyWithImpl<_ItineraryItem>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$ItineraryItemToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _ItineraryItem&&(identical(other.id, id) || other.id == id)&&(identical(other.itemType, itemType) || other.itemType == itemType)&&(identical(other.providerId, providerId) || other.providerId == providerId)&&(identical(other.externalReservationId, externalReservationId) || other.externalReservationId == externalReservationId)&&(identical(other.title, title) || other.title == title)&&(identical(other.status, status) || other.status == status)&&(identical(other.timeWindow, timeWindow) || other.timeWindow == timeWindow)&&(identical(other.originCode, originCode) || other.originCode == originCode)&&(identical(other.destinationCode, destinationCode) || other.destinationCode == destinationCode)&&(identical(other.price, price) || other.price == price)&&const DeepCollectionEquality().equals(other._operationalNotes, _operationalNotes)&&(identical(other.synthetic, synthetic) || other.synthetic == synthetic)&&(identical(other.source, source) || other.source == source));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,itemType,providerId,externalReservationId,title,status,timeWindow,originCode,destinationCode,price,const DeepCollectionEquality().hash(_operationalNotes),synthetic,source);

@override
String toString() {
  return 'ItineraryItem(id: $id, itemType: $itemType, providerId: $providerId, externalReservationId: $externalReservationId, title: $title, status: $status, timeWindow: $timeWindow, originCode: $originCode, destinationCode: $destinationCode, price: $price, operationalNotes: $operationalNotes, synthetic: $synthetic, source: $source)';
}


}

/// @nodoc
abstract mixin class _$ItineraryItemCopyWith<$Res> implements $ItineraryItemCopyWith<$Res> {
  factory _$ItineraryItemCopyWith(_ItineraryItem value, $Res Function(_ItineraryItem) _then) = __$ItineraryItemCopyWithImpl;
@override @useResult
$Res call({
 String id,@JsonKey(name: 'item_type') ItemType itemType,@JsonKey(name: 'provider_id') String? providerId,@JsonKey(name: 'external_reservation_id') String? externalReservationId, String title, ItineraryItemStatus status,@JsonKey(name: 'time_window') TimeWindow timeWindow,@JsonKey(name: 'origin_code') String? originCode,@JsonKey(name: 'destination_code') String? destinationCode, ConvertedMoney? price,@JsonKey(name: 'operational_notes') List<String> operationalNotes, bool synthetic, String source
});


@override $TimeWindowCopyWith<$Res> get timeWindow;@override $ConvertedMoneyCopyWith<$Res>? get price;

}
/// @nodoc
class __$ItineraryItemCopyWithImpl<$Res>
    implements _$ItineraryItemCopyWith<$Res> {
  __$ItineraryItemCopyWithImpl(this._self, this._then);

  final _ItineraryItem _self;
  final $Res Function(_ItineraryItem) _then;

/// Create a copy of ItineraryItem
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? itemType = null,Object? providerId = freezed,Object? externalReservationId = freezed,Object? title = null,Object? status = null,Object? timeWindow = null,Object? originCode = freezed,Object? destinationCode = freezed,Object? price = freezed,Object? operationalNotes = null,Object? synthetic = null,Object? source = null,}) {
  return _then(_ItineraryItem(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,itemType: null == itemType ? _self.itemType : itemType // ignore: cast_nullable_to_non_nullable
as ItemType,providerId: freezed == providerId ? _self.providerId : providerId // ignore: cast_nullable_to_non_nullable
as String?,externalReservationId: freezed == externalReservationId ? _self.externalReservationId : externalReservationId // ignore: cast_nullable_to_non_nullable
as String?,title: null == title ? _self.title : title // ignore: cast_nullable_to_non_nullable
as String,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as ItineraryItemStatus,timeWindow: null == timeWindow ? _self.timeWindow : timeWindow // ignore: cast_nullable_to_non_nullable
as TimeWindow,originCode: freezed == originCode ? _self.originCode : originCode // ignore: cast_nullable_to_non_nullable
as String?,destinationCode: freezed == destinationCode ? _self.destinationCode : destinationCode // ignore: cast_nullable_to_non_nullable
as String?,price: freezed == price ? _self.price : price // ignore: cast_nullable_to_non_nullable
as ConvertedMoney?,operationalNotes: null == operationalNotes ? _self._operationalNotes : operationalNotes // ignore: cast_nullable_to_non_nullable
as List<String>,synthetic: null == synthetic ? _self.synthetic : synthetic // ignore: cast_nullable_to_non_nullable
as bool,source: null == source ? _self.source : source // ignore: cast_nullable_to_non_nullable
as String,
  ));
}

/// Create a copy of ItineraryItem
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$TimeWindowCopyWith<$Res> get timeWindow {
  
  return $TimeWindowCopyWith<$Res>(_self.timeWindow, (value) {
    return _then(_self.copyWith(timeWindow: value));
  });
}/// Create a copy of ItineraryItem
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$ConvertedMoneyCopyWith<$Res>? get price {
    if (_self.price == null) {
    return null;
  }

  return $ConvertedMoneyCopyWith<$Res>(_self.price!, (value) {
    return _then(_self.copyWith(price: value));
  });
}
}


/// @nodoc
mixin _$ItineraryVersion {

 String get id;@JsonKey(name: 'journey_id') String get journeyId; int get version; ItineraryVersionStatus get status;@JsonKey(name: 'based_on_disruption_id') String? get basedOnDisruptionId;@JsonKey(name: 'total_price') PriceSummary get totalPrice; List<ItineraryItem> get items;@JsonKey(name: 'created_at') DateTime get createdAt;
/// Create a copy of ItineraryVersion
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$ItineraryVersionCopyWith<ItineraryVersion> get copyWith => _$ItineraryVersionCopyWithImpl<ItineraryVersion>(this as ItineraryVersion, _$identity);

  /// Serializes this ItineraryVersion to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is ItineraryVersion&&(identical(other.id, id) || other.id == id)&&(identical(other.journeyId, journeyId) || other.journeyId == journeyId)&&(identical(other.version, version) || other.version == version)&&(identical(other.status, status) || other.status == status)&&(identical(other.basedOnDisruptionId, basedOnDisruptionId) || other.basedOnDisruptionId == basedOnDisruptionId)&&(identical(other.totalPrice, totalPrice) || other.totalPrice == totalPrice)&&const DeepCollectionEquality().equals(other.items, items)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,journeyId,version,status,basedOnDisruptionId,totalPrice,const DeepCollectionEquality().hash(items),createdAt);

@override
String toString() {
  return 'ItineraryVersion(id: $id, journeyId: $journeyId, version: $version, status: $status, basedOnDisruptionId: $basedOnDisruptionId, totalPrice: $totalPrice, items: $items, createdAt: $createdAt)';
}


}

/// @nodoc
abstract mixin class $ItineraryVersionCopyWith<$Res>  {
  factory $ItineraryVersionCopyWith(ItineraryVersion value, $Res Function(ItineraryVersion) _then) = _$ItineraryVersionCopyWithImpl;
@useResult
$Res call({
 String id,@JsonKey(name: 'journey_id') String journeyId, int version, ItineraryVersionStatus status,@JsonKey(name: 'based_on_disruption_id') String? basedOnDisruptionId,@JsonKey(name: 'total_price') PriceSummary totalPrice, List<ItineraryItem> items,@JsonKey(name: 'created_at') DateTime createdAt
});


$PriceSummaryCopyWith<$Res> get totalPrice;

}
/// @nodoc
class _$ItineraryVersionCopyWithImpl<$Res>
    implements $ItineraryVersionCopyWith<$Res> {
  _$ItineraryVersionCopyWithImpl(this._self, this._then);

  final ItineraryVersion _self;
  final $Res Function(ItineraryVersion) _then;

/// Create a copy of ItineraryVersion
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? journeyId = null,Object? version = null,Object? status = null,Object? basedOnDisruptionId = freezed,Object? totalPrice = null,Object? items = null,Object? createdAt = null,}) {
  return _then(ItineraryVersion(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,journeyId: null == journeyId ? _self.journeyId : journeyId // ignore: cast_nullable_to_non_nullable
as String,version: null == version ? _self.version : version // ignore: cast_nullable_to_non_nullable
as int,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as ItineraryVersionStatus,basedOnDisruptionId: freezed == basedOnDisruptionId ? _self.basedOnDisruptionId : basedOnDisruptionId // ignore: cast_nullable_to_non_nullable
as String?,totalPrice: null == totalPrice ? _self.totalPrice : totalPrice // ignore: cast_nullable_to_non_nullable
as PriceSummary,items: null == items ? _self.items : items // ignore: cast_nullable_to_non_nullable
as List<ItineraryItem>,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}
/// Create a copy of ItineraryVersion
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$PriceSummaryCopyWith<$Res> get totalPrice {
  
  return $PriceSummaryCopyWith<$Res>(_self.totalPrice, (value) {
    return _then(_self.copyWith(totalPrice: value));
  });
}
}


/// Adds pattern-matching-related methods to [ItineraryVersion].
extension ItineraryVersionPatterns on ItineraryVersion {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _ItineraryVersion value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _ItineraryVersion() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _ItineraryVersion value)  $default,){
final _that = this;
switch (_that) {
case _ItineraryVersion():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _ItineraryVersion value)?  $default,){
final _that = this;
switch (_that) {
case _ItineraryVersion() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id, @JsonKey(name: 'journey_id')  String journeyId,  int version,  ItineraryVersionStatus status, @JsonKey(name: 'based_on_disruption_id')  String? basedOnDisruptionId, @JsonKey(name: 'total_price')  PriceSummary totalPrice,  List<ItineraryItem> items, @JsonKey(name: 'created_at')  DateTime createdAt)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _ItineraryVersion() when $default != null:
return $default(_that.id,_that.journeyId,_that.version,_that.status,_that.basedOnDisruptionId,_that.totalPrice,_that.items,_that.createdAt);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id, @JsonKey(name: 'journey_id')  String journeyId,  int version,  ItineraryVersionStatus status, @JsonKey(name: 'based_on_disruption_id')  String? basedOnDisruptionId, @JsonKey(name: 'total_price')  PriceSummary totalPrice,  List<ItineraryItem> items, @JsonKey(name: 'created_at')  DateTime createdAt)  $default,) {final _that = this;
switch (_that) {
case _ItineraryVersion():
return $default(_that.id,_that.journeyId,_that.version,_that.status,_that.basedOnDisruptionId,_that.totalPrice,_that.items,_that.createdAt);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id, @JsonKey(name: 'journey_id')  String journeyId,  int version,  ItineraryVersionStatus status, @JsonKey(name: 'based_on_disruption_id')  String? basedOnDisruptionId, @JsonKey(name: 'total_price')  PriceSummary totalPrice,  List<ItineraryItem> items, @JsonKey(name: 'created_at')  DateTime createdAt)?  $default,) {final _that = this;
switch (_that) {
case _ItineraryVersion() when $default != null:
return $default(_that.id,_that.journeyId,_that.version,_that.status,_that.basedOnDisruptionId,_that.totalPrice,_that.items,_that.createdAt);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _ItineraryVersion implements ItineraryVersion {
  const _ItineraryVersion({required this.id, @JsonKey(name: 'journey_id') required this.journeyId, required this.version, required this.status, @JsonKey(name: 'based_on_disruption_id') this.basedOnDisruptionId, @JsonKey(name: 'total_price') required this.totalPrice, required  List<ItineraryItem> items, @JsonKey(name: 'created_at') required this.createdAt}): _items = items;
  factory _ItineraryVersion.fromJson(Map<String, dynamic> json) => _$ItineraryVersionFromJson(json);

@override final  String id;
@override@JsonKey(name: 'journey_id') final  String journeyId;
@override final  int version;
@override final  ItineraryVersionStatus status;
@override@JsonKey(name: 'based_on_disruption_id') final  String? basedOnDisruptionId;
@override@JsonKey(name: 'total_price') final  PriceSummary totalPrice;
 final  List<ItineraryItem> _items;
@override List<ItineraryItem> get items {
  if (_items is EqualUnmodifiableListView) return _items;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_items);
}

@override@JsonKey(name: 'created_at') final  DateTime createdAt;

/// Create a copy of ItineraryVersion
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$ItineraryVersionCopyWith<_ItineraryVersion> get copyWith => __$ItineraryVersionCopyWithImpl<_ItineraryVersion>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$ItineraryVersionToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _ItineraryVersion&&(identical(other.id, id) || other.id == id)&&(identical(other.journeyId, journeyId) || other.journeyId == journeyId)&&(identical(other.version, version) || other.version == version)&&(identical(other.status, status) || other.status == status)&&(identical(other.basedOnDisruptionId, basedOnDisruptionId) || other.basedOnDisruptionId == basedOnDisruptionId)&&(identical(other.totalPrice, totalPrice) || other.totalPrice == totalPrice)&&const DeepCollectionEquality().equals(other._items, _items)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,journeyId,version,status,basedOnDisruptionId,totalPrice,const DeepCollectionEquality().hash(_items),createdAt);

@override
String toString() {
  return 'ItineraryVersion(id: $id, journeyId: $journeyId, version: $version, status: $status, basedOnDisruptionId: $basedOnDisruptionId, totalPrice: $totalPrice, items: $items, createdAt: $createdAt)';
}


}

/// @nodoc
abstract mixin class _$ItineraryVersionCopyWith<$Res> implements $ItineraryVersionCopyWith<$Res> {
  factory _$ItineraryVersionCopyWith(_ItineraryVersion value, $Res Function(_ItineraryVersion) _then) = __$ItineraryVersionCopyWithImpl;
@override @useResult
$Res call({
 String id,@JsonKey(name: 'journey_id') String journeyId, int version, ItineraryVersionStatus status,@JsonKey(name: 'based_on_disruption_id') String? basedOnDisruptionId,@JsonKey(name: 'total_price') PriceSummary totalPrice, List<ItineraryItem> items,@JsonKey(name: 'created_at') DateTime createdAt
});


@override $PriceSummaryCopyWith<$Res> get totalPrice;

}
/// @nodoc
class __$ItineraryVersionCopyWithImpl<$Res>
    implements _$ItineraryVersionCopyWith<$Res> {
  __$ItineraryVersionCopyWithImpl(this._self, this._then);

  final _ItineraryVersion _self;
  final $Res Function(_ItineraryVersion) _then;

/// Create a copy of ItineraryVersion
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? journeyId = null,Object? version = null,Object? status = null,Object? basedOnDisruptionId = freezed,Object? totalPrice = null,Object? items = null,Object? createdAt = null,}) {
  return _then(_ItineraryVersion(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,journeyId: null == journeyId ? _self.journeyId : journeyId // ignore: cast_nullable_to_non_nullable
as String,version: null == version ? _self.version : version // ignore: cast_nullable_to_non_nullable
as int,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as ItineraryVersionStatus,basedOnDisruptionId: freezed == basedOnDisruptionId ? _self.basedOnDisruptionId : basedOnDisruptionId // ignore: cast_nullable_to_non_nullable
as String?,totalPrice: null == totalPrice ? _self.totalPrice : totalPrice // ignore: cast_nullable_to_non_nullable
as PriceSummary,items: null == items ? _self._items : items // ignore: cast_nullable_to_non_nullable
as List<ItineraryItem>,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}

/// Create a copy of ItineraryVersion
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$PriceSummaryCopyWith<$Res> get totalPrice {
  
  return $PriceSummaryCopyWith<$Res>(_self.totalPrice, (value) {
    return _then(_self.copyWith(totalPrice: value));
  });
}
}


/// @nodoc
mixin _$ItineraryVersionSummary {

 String get id; int get version; ItineraryVersionStatus get status;@JsonKey(name: 'based_on_disruption_id') String? get basedOnDisruptionId;@JsonKey(name: 'created_at') DateTime get createdAt;
/// Create a copy of ItineraryVersionSummary
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$ItineraryVersionSummaryCopyWith<ItineraryVersionSummary> get copyWith => _$ItineraryVersionSummaryCopyWithImpl<ItineraryVersionSummary>(this as ItineraryVersionSummary, _$identity);

  /// Serializes this ItineraryVersionSummary to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is ItineraryVersionSummary&&(identical(other.id, id) || other.id == id)&&(identical(other.version, version) || other.version == version)&&(identical(other.status, status) || other.status == status)&&(identical(other.basedOnDisruptionId, basedOnDisruptionId) || other.basedOnDisruptionId == basedOnDisruptionId)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,version,status,basedOnDisruptionId,createdAt);

@override
String toString() {
  return 'ItineraryVersionSummary(id: $id, version: $version, status: $status, basedOnDisruptionId: $basedOnDisruptionId, createdAt: $createdAt)';
}


}

/// @nodoc
abstract mixin class $ItineraryVersionSummaryCopyWith<$Res>  {
  factory $ItineraryVersionSummaryCopyWith(ItineraryVersionSummary value, $Res Function(ItineraryVersionSummary) _then) = _$ItineraryVersionSummaryCopyWithImpl;
@useResult
$Res call({
 String id, int version, ItineraryVersionStatus status,@JsonKey(name: 'based_on_disruption_id') String? basedOnDisruptionId,@JsonKey(name: 'created_at') DateTime createdAt
});




}
/// @nodoc
class _$ItineraryVersionSummaryCopyWithImpl<$Res>
    implements $ItineraryVersionSummaryCopyWith<$Res> {
  _$ItineraryVersionSummaryCopyWithImpl(this._self, this._then);

  final ItineraryVersionSummary _self;
  final $Res Function(ItineraryVersionSummary) _then;

/// Create a copy of ItineraryVersionSummary
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? version = null,Object? status = null,Object? basedOnDisruptionId = freezed,Object? createdAt = null,}) {
  return _then(ItineraryVersionSummary(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,version: null == version ? _self.version : version // ignore: cast_nullable_to_non_nullable
as int,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as ItineraryVersionStatus,basedOnDisruptionId: freezed == basedOnDisruptionId ? _self.basedOnDisruptionId : basedOnDisruptionId // ignore: cast_nullable_to_non_nullable
as String?,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}

}


/// Adds pattern-matching-related methods to [ItineraryVersionSummary].
extension ItineraryVersionSummaryPatterns on ItineraryVersionSummary {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _ItineraryVersionSummary value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _ItineraryVersionSummary() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _ItineraryVersionSummary value)  $default,){
final _that = this;
switch (_that) {
case _ItineraryVersionSummary():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _ItineraryVersionSummary value)?  $default,){
final _that = this;
switch (_that) {
case _ItineraryVersionSummary() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  int version,  ItineraryVersionStatus status, @JsonKey(name: 'based_on_disruption_id')  String? basedOnDisruptionId, @JsonKey(name: 'created_at')  DateTime createdAt)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _ItineraryVersionSummary() when $default != null:
return $default(_that.id,_that.version,_that.status,_that.basedOnDisruptionId,_that.createdAt);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  int version,  ItineraryVersionStatus status, @JsonKey(name: 'based_on_disruption_id')  String? basedOnDisruptionId, @JsonKey(name: 'created_at')  DateTime createdAt)  $default,) {final _that = this;
switch (_that) {
case _ItineraryVersionSummary():
return $default(_that.id,_that.version,_that.status,_that.basedOnDisruptionId,_that.createdAt);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  int version,  ItineraryVersionStatus status, @JsonKey(name: 'based_on_disruption_id')  String? basedOnDisruptionId, @JsonKey(name: 'created_at')  DateTime createdAt)?  $default,) {final _that = this;
switch (_that) {
case _ItineraryVersionSummary() when $default != null:
return $default(_that.id,_that.version,_that.status,_that.basedOnDisruptionId,_that.createdAt);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _ItineraryVersionSummary implements ItineraryVersionSummary {
  const _ItineraryVersionSummary({required this.id, required this.version, required this.status, @JsonKey(name: 'based_on_disruption_id') this.basedOnDisruptionId, @JsonKey(name: 'created_at') required this.createdAt});
  factory _ItineraryVersionSummary.fromJson(Map<String, dynamic> json) => _$ItineraryVersionSummaryFromJson(json);

@override final  String id;
@override final  int version;
@override final  ItineraryVersionStatus status;
@override@JsonKey(name: 'based_on_disruption_id') final  String? basedOnDisruptionId;
@override@JsonKey(name: 'created_at') final  DateTime createdAt;

/// Create a copy of ItineraryVersionSummary
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$ItineraryVersionSummaryCopyWith<_ItineraryVersionSummary> get copyWith => __$ItineraryVersionSummaryCopyWithImpl<_ItineraryVersionSummary>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$ItineraryVersionSummaryToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _ItineraryVersionSummary&&(identical(other.id, id) || other.id == id)&&(identical(other.version, version) || other.version == version)&&(identical(other.status, status) || other.status == status)&&(identical(other.basedOnDisruptionId, basedOnDisruptionId) || other.basedOnDisruptionId == basedOnDisruptionId)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,version,status,basedOnDisruptionId,createdAt);

@override
String toString() {
  return 'ItineraryVersionSummary(id: $id, version: $version, status: $status, basedOnDisruptionId: $basedOnDisruptionId, createdAt: $createdAt)';
}


}

/// @nodoc
abstract mixin class _$ItineraryVersionSummaryCopyWith<$Res> implements $ItineraryVersionSummaryCopyWith<$Res> {
  factory _$ItineraryVersionSummaryCopyWith(_ItineraryVersionSummary value, $Res Function(_ItineraryVersionSummary) _then) = __$ItineraryVersionSummaryCopyWithImpl;
@override @useResult
$Res call({
 String id, int version, ItineraryVersionStatus status,@JsonKey(name: 'based_on_disruption_id') String? basedOnDisruptionId,@JsonKey(name: 'created_at') DateTime createdAt
});




}
/// @nodoc
class __$ItineraryVersionSummaryCopyWithImpl<$Res>
    implements _$ItineraryVersionSummaryCopyWith<$Res> {
  __$ItineraryVersionSummaryCopyWithImpl(this._self, this._then);

  final _ItineraryVersionSummary _self;
  final $Res Function(_ItineraryVersionSummary) _then;

/// Create a copy of ItineraryVersionSummary
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? version = null,Object? status = null,Object? basedOnDisruptionId = freezed,Object? createdAt = null,}) {
  return _then(_ItineraryVersionSummary(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,version: null == version ? _self.version : version // ignore: cast_nullable_to_non_nullable
as int,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as ItineraryVersionStatus,basedOnDisruptionId: freezed == basedOnDisruptionId ? _self.basedOnDisruptionId : basedOnDisruptionId // ignore: cast_nullable_to_non_nullable
as String?,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}


}


/// @nodoc
mixin _$JourneyDetail {

 Journey get journey;@JsonKey(name: 'active_itinerary') ItineraryVersion get activeItinerary;@JsonKey(name: 'itinerary_versions') List<ItineraryVersionSummary> get itineraryVersions;
/// Create a copy of JourneyDetail
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$JourneyDetailCopyWith<JourneyDetail> get copyWith => _$JourneyDetailCopyWithImpl<JourneyDetail>(this as JourneyDetail, _$identity);

  /// Serializes this JourneyDetail to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is JourneyDetail&&(identical(other.journey, journey) || other.journey == journey)&&(identical(other.activeItinerary, activeItinerary) || other.activeItinerary == activeItinerary)&&const DeepCollectionEquality().equals(other.itineraryVersions, itineraryVersions));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,journey,activeItinerary,const DeepCollectionEquality().hash(itineraryVersions));

@override
String toString() {
  return 'JourneyDetail(journey: $journey, activeItinerary: $activeItinerary, itineraryVersions: $itineraryVersions)';
}


}

/// @nodoc
abstract mixin class $JourneyDetailCopyWith<$Res>  {
  factory $JourneyDetailCopyWith(JourneyDetail value, $Res Function(JourneyDetail) _then) = _$JourneyDetailCopyWithImpl;
@useResult
$Res call({
 Journey journey,@JsonKey(name: 'active_itinerary') ItineraryVersion activeItinerary,@JsonKey(name: 'itinerary_versions') List<ItineraryVersionSummary> itineraryVersions
});


$JourneyCopyWith<$Res> get journey;$ItineraryVersionCopyWith<$Res> get activeItinerary;

}
/// @nodoc
class _$JourneyDetailCopyWithImpl<$Res>
    implements $JourneyDetailCopyWith<$Res> {
  _$JourneyDetailCopyWithImpl(this._self, this._then);

  final JourneyDetail _self;
  final $Res Function(JourneyDetail) _then;

/// Create a copy of JourneyDetail
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? journey = null,Object? activeItinerary = null,Object? itineraryVersions = null,}) {
  return _then(JourneyDetail(
journey: null == journey ? _self.journey : journey // ignore: cast_nullable_to_non_nullable
as Journey,activeItinerary: null == activeItinerary ? _self.activeItinerary : activeItinerary // ignore: cast_nullable_to_non_nullable
as ItineraryVersion,itineraryVersions: null == itineraryVersions ? _self.itineraryVersions : itineraryVersions // ignore: cast_nullable_to_non_nullable
as List<ItineraryVersionSummary>,
  ));
}
/// Create a copy of JourneyDetail
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$JourneyCopyWith<$Res> get journey {
  
  return $JourneyCopyWith<$Res>(_self.journey, (value) {
    return _then(_self.copyWith(journey: value));
  });
}/// Create a copy of JourneyDetail
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$ItineraryVersionCopyWith<$Res> get activeItinerary {
  
  return $ItineraryVersionCopyWith<$Res>(_self.activeItinerary, (value) {
    return _then(_self.copyWith(activeItinerary: value));
  });
}
}


/// Adds pattern-matching-related methods to [JourneyDetail].
extension JourneyDetailPatterns on JourneyDetail {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _JourneyDetail value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _JourneyDetail() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _JourneyDetail value)  $default,){
final _that = this;
switch (_that) {
case _JourneyDetail():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _JourneyDetail value)?  $default,){
final _that = this;
switch (_that) {
case _JourneyDetail() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( Journey journey, @JsonKey(name: 'active_itinerary')  ItineraryVersion activeItinerary, @JsonKey(name: 'itinerary_versions')  List<ItineraryVersionSummary> itineraryVersions)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _JourneyDetail() when $default != null:
return $default(_that.journey,_that.activeItinerary,_that.itineraryVersions);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( Journey journey, @JsonKey(name: 'active_itinerary')  ItineraryVersion activeItinerary, @JsonKey(name: 'itinerary_versions')  List<ItineraryVersionSummary> itineraryVersions)  $default,) {final _that = this;
switch (_that) {
case _JourneyDetail():
return $default(_that.journey,_that.activeItinerary,_that.itineraryVersions);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( Journey journey, @JsonKey(name: 'active_itinerary')  ItineraryVersion activeItinerary, @JsonKey(name: 'itinerary_versions')  List<ItineraryVersionSummary> itineraryVersions)?  $default,) {final _that = this;
switch (_that) {
case _JourneyDetail() when $default != null:
return $default(_that.journey,_that.activeItinerary,_that.itineraryVersions);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _JourneyDetail implements JourneyDetail {
  const _JourneyDetail({required this.journey, @JsonKey(name: 'active_itinerary') required this.activeItinerary, @JsonKey(name: 'itinerary_versions') required  List<ItineraryVersionSummary> itineraryVersions}): _itineraryVersions = itineraryVersions;
  factory _JourneyDetail.fromJson(Map<String, dynamic> json) => _$JourneyDetailFromJson(json);

@override final  Journey journey;
@override@JsonKey(name: 'active_itinerary') final  ItineraryVersion activeItinerary;
 final  List<ItineraryVersionSummary> _itineraryVersions;
@override@JsonKey(name: 'itinerary_versions') List<ItineraryVersionSummary> get itineraryVersions {
  if (_itineraryVersions is EqualUnmodifiableListView) return _itineraryVersions;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_itineraryVersions);
}


/// Create a copy of JourneyDetail
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$JourneyDetailCopyWith<_JourneyDetail> get copyWith => __$JourneyDetailCopyWithImpl<_JourneyDetail>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$JourneyDetailToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _JourneyDetail&&(identical(other.journey, journey) || other.journey == journey)&&(identical(other.activeItinerary, activeItinerary) || other.activeItinerary == activeItinerary)&&const DeepCollectionEquality().equals(other._itineraryVersions, _itineraryVersions));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,journey,activeItinerary,const DeepCollectionEquality().hash(_itineraryVersions));

@override
String toString() {
  return 'JourneyDetail(journey: $journey, activeItinerary: $activeItinerary, itineraryVersions: $itineraryVersions)';
}


}

/// @nodoc
abstract mixin class _$JourneyDetailCopyWith<$Res> implements $JourneyDetailCopyWith<$Res> {
  factory _$JourneyDetailCopyWith(_JourneyDetail value, $Res Function(_JourneyDetail) _then) = __$JourneyDetailCopyWithImpl;
@override @useResult
$Res call({
 Journey journey,@JsonKey(name: 'active_itinerary') ItineraryVersion activeItinerary,@JsonKey(name: 'itinerary_versions') List<ItineraryVersionSummary> itineraryVersions
});


@override $JourneyCopyWith<$Res> get journey;@override $ItineraryVersionCopyWith<$Res> get activeItinerary;

}
/// @nodoc
class __$JourneyDetailCopyWithImpl<$Res>
    implements _$JourneyDetailCopyWith<$Res> {
  __$JourneyDetailCopyWithImpl(this._self, this._then);

  final _JourneyDetail _self;
  final $Res Function(_JourneyDetail) _then;

/// Create a copy of JourneyDetail
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? journey = null,Object? activeItinerary = null,Object? itineraryVersions = null,}) {
  return _then(_JourneyDetail(
journey: null == journey ? _self.journey : journey // ignore: cast_nullable_to_non_nullable
as Journey,activeItinerary: null == activeItinerary ? _self.activeItinerary : activeItinerary // ignore: cast_nullable_to_non_nullable
as ItineraryVersion,itineraryVersions: null == itineraryVersions ? _self._itineraryVersions : itineraryVersions // ignore: cast_nullable_to_non_nullable
as List<ItineraryVersionSummary>,
  ));
}

/// Create a copy of JourneyDetail
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$JourneyCopyWith<$Res> get journey {
  
  return $JourneyCopyWith<$Res>(_self.journey, (value) {
    return _then(_self.copyWith(journey: value));
  });
}/// Create a copy of JourneyDetail
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$ItineraryVersionCopyWith<$Res> get activeItinerary {
  
  return $ItineraryVersionCopyWith<$Res>(_self.activeItinerary, (value) {
    return _then(_self.copyWith(activeItinerary: value));
  });
}
}

// dart format on
