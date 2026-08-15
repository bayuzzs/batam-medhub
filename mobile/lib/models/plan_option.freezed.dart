// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint, type=warning, deprecated_member_use, deprecated_member_use_from_same_package
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'plan_option.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// dart format off
T _$identity<T>(T value) => value;

/// @nodoc
mixin _$PlanItem {

 String get id;@JsonKey(name: 'item_type') ItemType get itemType;@JsonKey(name: 'provider_id') String? get providerId;@JsonKey(name: 'external_offer_id') String? get externalOfferId; String get title;@JsonKey(name: 'time_window') TimeWindow get timeWindow;@JsonKey(name: 'origin_code') String? get originCode;@JsonKey(name: 'destination_code') String? get destinationCode; ConvertedMoney? get price;@JsonKey(name: 'offer_expires_at') DateTime? get offerExpiresAt;@JsonKey(name: 'operational_notes') List<String> get operationalNotes; bool get synthetic; String get source;
/// Create a copy of PlanItem
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$PlanItemCopyWith<PlanItem> get copyWith => _$PlanItemCopyWithImpl<PlanItem>(this as PlanItem, _$identity);

  /// Serializes this PlanItem to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is PlanItem&&(identical(other.id, id) || other.id == id)&&(identical(other.itemType, itemType) || other.itemType == itemType)&&(identical(other.providerId, providerId) || other.providerId == providerId)&&(identical(other.externalOfferId, externalOfferId) || other.externalOfferId == externalOfferId)&&(identical(other.title, title) || other.title == title)&&(identical(other.timeWindow, timeWindow) || other.timeWindow == timeWindow)&&(identical(other.originCode, originCode) || other.originCode == originCode)&&(identical(other.destinationCode, destinationCode) || other.destinationCode == destinationCode)&&(identical(other.price, price) || other.price == price)&&(identical(other.offerExpiresAt, offerExpiresAt) || other.offerExpiresAt == offerExpiresAt)&&const DeepCollectionEquality().equals(other.operationalNotes, operationalNotes)&&(identical(other.synthetic, synthetic) || other.synthetic == synthetic)&&(identical(other.source, source) || other.source == source));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,itemType,providerId,externalOfferId,title,timeWindow,originCode,destinationCode,price,offerExpiresAt,const DeepCollectionEquality().hash(operationalNotes),synthetic,source);

@override
String toString() {
  return 'PlanItem(id: $id, itemType: $itemType, providerId: $providerId, externalOfferId: $externalOfferId, title: $title, timeWindow: $timeWindow, originCode: $originCode, destinationCode: $destinationCode, price: $price, offerExpiresAt: $offerExpiresAt, operationalNotes: $operationalNotes, synthetic: $synthetic, source: $source)';
}


}

/// @nodoc
abstract mixin class $PlanItemCopyWith<$Res>  {
  factory $PlanItemCopyWith(PlanItem value, $Res Function(PlanItem) _then) = _$PlanItemCopyWithImpl;
@useResult
$Res call({
 String id,@JsonKey(name: 'item_type') ItemType itemType,@JsonKey(name: 'provider_id') String? providerId,@JsonKey(name: 'external_offer_id') String? externalOfferId, String title,@JsonKey(name: 'time_window') TimeWindow timeWindow,@JsonKey(name: 'origin_code') String? originCode,@JsonKey(name: 'destination_code') String? destinationCode, ConvertedMoney? price,@JsonKey(name: 'offer_expires_at') DateTime? offerExpiresAt,@JsonKey(name: 'operational_notes') List<String> operationalNotes, bool synthetic, String source
});


$TimeWindowCopyWith<$Res> get timeWindow;$ConvertedMoneyCopyWith<$Res>? get price;

}
/// @nodoc
class _$PlanItemCopyWithImpl<$Res>
    implements $PlanItemCopyWith<$Res> {
  _$PlanItemCopyWithImpl(this._self, this._then);

  final PlanItem _self;
  final $Res Function(PlanItem) _then;

/// Create a copy of PlanItem
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? itemType = null,Object? providerId = freezed,Object? externalOfferId = freezed,Object? title = null,Object? timeWindow = null,Object? originCode = freezed,Object? destinationCode = freezed,Object? price = freezed,Object? offerExpiresAt = freezed,Object? operationalNotes = null,Object? synthetic = null,Object? source = null,}) {
  return _then(PlanItem(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,itemType: null == itemType ? _self.itemType : itemType // ignore: cast_nullable_to_non_nullable
as ItemType,providerId: freezed == providerId ? _self.providerId : providerId // ignore: cast_nullable_to_non_nullable
as String?,externalOfferId: freezed == externalOfferId ? _self.externalOfferId : externalOfferId // ignore: cast_nullable_to_non_nullable
as String?,title: null == title ? _self.title : title // ignore: cast_nullable_to_non_nullable
as String,timeWindow: null == timeWindow ? _self.timeWindow : timeWindow // ignore: cast_nullable_to_non_nullable
as TimeWindow,originCode: freezed == originCode ? _self.originCode : originCode // ignore: cast_nullable_to_non_nullable
as String?,destinationCode: freezed == destinationCode ? _self.destinationCode : destinationCode // ignore: cast_nullable_to_non_nullable
as String?,price: freezed == price ? _self.price : price // ignore: cast_nullable_to_non_nullable
as ConvertedMoney?,offerExpiresAt: freezed == offerExpiresAt ? _self.offerExpiresAt : offerExpiresAt // ignore: cast_nullable_to_non_nullable
as DateTime?,operationalNotes: null == operationalNotes ? _self.operationalNotes : operationalNotes // ignore: cast_nullable_to_non_nullable
as List<String>,synthetic: null == synthetic ? _self.synthetic : synthetic // ignore: cast_nullable_to_non_nullable
as bool,source: null == source ? _self.source : source // ignore: cast_nullable_to_non_nullable
as String,
  ));
}
/// Create a copy of PlanItem
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$TimeWindowCopyWith<$Res> get timeWindow {
  
  return $TimeWindowCopyWith<$Res>(_self.timeWindow, (value) {
    return _then(_self.copyWith(timeWindow: value));
  });
}/// Create a copy of PlanItem
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


/// Adds pattern-matching-related methods to [PlanItem].
extension PlanItemPatterns on PlanItem {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _PlanItem value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _PlanItem() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _PlanItem value)  $default,){
final _that = this;
switch (_that) {
case _PlanItem():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _PlanItem value)?  $default,){
final _that = this;
switch (_that) {
case _PlanItem() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id, @JsonKey(name: 'item_type')  ItemType itemType, @JsonKey(name: 'provider_id')  String? providerId, @JsonKey(name: 'external_offer_id')  String? externalOfferId,  String title, @JsonKey(name: 'time_window')  TimeWindow timeWindow, @JsonKey(name: 'origin_code')  String? originCode, @JsonKey(name: 'destination_code')  String? destinationCode,  ConvertedMoney? price, @JsonKey(name: 'offer_expires_at')  DateTime? offerExpiresAt, @JsonKey(name: 'operational_notes')  List<String> operationalNotes,  bool synthetic,  String source)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _PlanItem() when $default != null:
return $default(_that.id,_that.itemType,_that.providerId,_that.externalOfferId,_that.title,_that.timeWindow,_that.originCode,_that.destinationCode,_that.price,_that.offerExpiresAt,_that.operationalNotes,_that.synthetic,_that.source);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id, @JsonKey(name: 'item_type')  ItemType itemType, @JsonKey(name: 'provider_id')  String? providerId, @JsonKey(name: 'external_offer_id')  String? externalOfferId,  String title, @JsonKey(name: 'time_window')  TimeWindow timeWindow, @JsonKey(name: 'origin_code')  String? originCode, @JsonKey(name: 'destination_code')  String? destinationCode,  ConvertedMoney? price, @JsonKey(name: 'offer_expires_at')  DateTime? offerExpiresAt, @JsonKey(name: 'operational_notes')  List<String> operationalNotes,  bool synthetic,  String source)  $default,) {final _that = this;
switch (_that) {
case _PlanItem():
return $default(_that.id,_that.itemType,_that.providerId,_that.externalOfferId,_that.title,_that.timeWindow,_that.originCode,_that.destinationCode,_that.price,_that.offerExpiresAt,_that.operationalNotes,_that.synthetic,_that.source);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id, @JsonKey(name: 'item_type')  ItemType itemType, @JsonKey(name: 'provider_id')  String? providerId, @JsonKey(name: 'external_offer_id')  String? externalOfferId,  String title, @JsonKey(name: 'time_window')  TimeWindow timeWindow, @JsonKey(name: 'origin_code')  String? originCode, @JsonKey(name: 'destination_code')  String? destinationCode,  ConvertedMoney? price, @JsonKey(name: 'offer_expires_at')  DateTime? offerExpiresAt, @JsonKey(name: 'operational_notes')  List<String> operationalNotes,  bool synthetic,  String source)?  $default,) {final _that = this;
switch (_that) {
case _PlanItem() when $default != null:
return $default(_that.id,_that.itemType,_that.providerId,_that.externalOfferId,_that.title,_that.timeWindow,_that.originCode,_that.destinationCode,_that.price,_that.offerExpiresAt,_that.operationalNotes,_that.synthetic,_that.source);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _PlanItem implements PlanItem {
  const _PlanItem({required this.id, @JsonKey(name: 'item_type') required this.itemType, @JsonKey(name: 'provider_id') this.providerId, @JsonKey(name: 'external_offer_id') this.externalOfferId, required this.title, @JsonKey(name: 'time_window') required this.timeWindow, @JsonKey(name: 'origin_code') this.originCode, @JsonKey(name: 'destination_code') this.destinationCode, this.price, @JsonKey(name: 'offer_expires_at') this.offerExpiresAt, @JsonKey(name: 'operational_notes') required  List<String> operationalNotes, required this.synthetic, required this.source}): _operationalNotes = operationalNotes;
  factory _PlanItem.fromJson(Map<String, dynamic> json) => _$PlanItemFromJson(json);

@override final  String id;
@override@JsonKey(name: 'item_type') final  ItemType itemType;
@override@JsonKey(name: 'provider_id') final  String? providerId;
@override@JsonKey(name: 'external_offer_id') final  String? externalOfferId;
@override final  String title;
@override@JsonKey(name: 'time_window') final  TimeWindow timeWindow;
@override@JsonKey(name: 'origin_code') final  String? originCode;
@override@JsonKey(name: 'destination_code') final  String? destinationCode;
@override final  ConvertedMoney? price;
@override@JsonKey(name: 'offer_expires_at') final  DateTime? offerExpiresAt;
 final  List<String> _operationalNotes;
@override@JsonKey(name: 'operational_notes') List<String> get operationalNotes {
  if (_operationalNotes is EqualUnmodifiableListView) return _operationalNotes;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_operationalNotes);
}

@override final  bool synthetic;
@override final  String source;

/// Create a copy of PlanItem
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$PlanItemCopyWith<_PlanItem> get copyWith => __$PlanItemCopyWithImpl<_PlanItem>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$PlanItemToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _PlanItem&&(identical(other.id, id) || other.id == id)&&(identical(other.itemType, itemType) || other.itemType == itemType)&&(identical(other.providerId, providerId) || other.providerId == providerId)&&(identical(other.externalOfferId, externalOfferId) || other.externalOfferId == externalOfferId)&&(identical(other.title, title) || other.title == title)&&(identical(other.timeWindow, timeWindow) || other.timeWindow == timeWindow)&&(identical(other.originCode, originCode) || other.originCode == originCode)&&(identical(other.destinationCode, destinationCode) || other.destinationCode == destinationCode)&&(identical(other.price, price) || other.price == price)&&(identical(other.offerExpiresAt, offerExpiresAt) || other.offerExpiresAt == offerExpiresAt)&&const DeepCollectionEquality().equals(other._operationalNotes, _operationalNotes)&&(identical(other.synthetic, synthetic) || other.synthetic == synthetic)&&(identical(other.source, source) || other.source == source));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,itemType,providerId,externalOfferId,title,timeWindow,originCode,destinationCode,price,offerExpiresAt,const DeepCollectionEquality().hash(_operationalNotes),synthetic,source);

@override
String toString() {
  return 'PlanItem(id: $id, itemType: $itemType, providerId: $providerId, externalOfferId: $externalOfferId, title: $title, timeWindow: $timeWindow, originCode: $originCode, destinationCode: $destinationCode, price: $price, offerExpiresAt: $offerExpiresAt, operationalNotes: $operationalNotes, synthetic: $synthetic, source: $source)';
}


}

/// @nodoc
abstract mixin class _$PlanItemCopyWith<$Res> implements $PlanItemCopyWith<$Res> {
  factory _$PlanItemCopyWith(_PlanItem value, $Res Function(_PlanItem) _then) = __$PlanItemCopyWithImpl;
@override @useResult
$Res call({
 String id,@JsonKey(name: 'item_type') ItemType itemType,@JsonKey(name: 'provider_id') String? providerId,@JsonKey(name: 'external_offer_id') String? externalOfferId, String title,@JsonKey(name: 'time_window') TimeWindow timeWindow,@JsonKey(name: 'origin_code') String? originCode,@JsonKey(name: 'destination_code') String? destinationCode, ConvertedMoney? price,@JsonKey(name: 'offer_expires_at') DateTime? offerExpiresAt,@JsonKey(name: 'operational_notes') List<String> operationalNotes, bool synthetic, String source
});


@override $TimeWindowCopyWith<$Res> get timeWindow;@override $ConvertedMoneyCopyWith<$Res>? get price;

}
/// @nodoc
class __$PlanItemCopyWithImpl<$Res>
    implements _$PlanItemCopyWith<$Res> {
  __$PlanItemCopyWithImpl(this._self, this._then);

  final _PlanItem _self;
  final $Res Function(_PlanItem) _then;

/// Create a copy of PlanItem
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? itemType = null,Object? providerId = freezed,Object? externalOfferId = freezed,Object? title = null,Object? timeWindow = null,Object? originCode = freezed,Object? destinationCode = freezed,Object? price = freezed,Object? offerExpiresAt = freezed,Object? operationalNotes = null,Object? synthetic = null,Object? source = null,}) {
  return _then(_PlanItem(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,itemType: null == itemType ? _self.itemType : itemType // ignore: cast_nullable_to_non_nullable
as ItemType,providerId: freezed == providerId ? _self.providerId : providerId // ignore: cast_nullable_to_non_nullable
as String?,externalOfferId: freezed == externalOfferId ? _self.externalOfferId : externalOfferId // ignore: cast_nullable_to_non_nullable
as String?,title: null == title ? _self.title : title // ignore: cast_nullable_to_non_nullable
as String,timeWindow: null == timeWindow ? _self.timeWindow : timeWindow // ignore: cast_nullable_to_non_nullable
as TimeWindow,originCode: freezed == originCode ? _self.originCode : originCode // ignore: cast_nullable_to_non_nullable
as String?,destinationCode: freezed == destinationCode ? _self.destinationCode : destinationCode // ignore: cast_nullable_to_non_nullable
as String?,price: freezed == price ? _self.price : price // ignore: cast_nullable_to_non_nullable
as ConvertedMoney?,offerExpiresAt: freezed == offerExpiresAt ? _self.offerExpiresAt : offerExpiresAt // ignore: cast_nullable_to_non_nullable
as DateTime?,operationalNotes: null == operationalNotes ? _self._operationalNotes : operationalNotes // ignore: cast_nullable_to_non_nullable
as List<String>,synthetic: null == synthetic ? _self.synthetic : synthetic // ignore: cast_nullable_to_non_nullable
as bool,source: null == source ? _self.source : source // ignore: cast_nullable_to_non_nullable
as String,
  ));
}

/// Create a copy of PlanItem
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$TimeWindowCopyWith<$Res> get timeWindow {
  
  return $TimeWindowCopyWith<$Res>(_self.timeWindow, (value) {
    return _then(_self.copyWith(timeWindow: value));
  });
}/// Create a copy of PlanItem
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
mixin _$PlanOption {

 String get id;@JsonKey(name: 'trip_request_id') String get tripRequestId;@JsonKey(name: 'planning_revision') int get planningRevision; int get rank; PlanOptionStatus get status;@JsonKey(name: 'expires_at') DateTime get expiresAt; List<String> get explanation; List<PlanItem> get items;@JsonKey(name: 'total_price') PriceSummary get totalPrice;
/// Create a copy of PlanOption
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$PlanOptionCopyWith<PlanOption> get copyWith => _$PlanOptionCopyWithImpl<PlanOption>(this as PlanOption, _$identity);

  /// Serializes this PlanOption to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is PlanOption&&(identical(other.id, id) || other.id == id)&&(identical(other.tripRequestId, tripRequestId) || other.tripRequestId == tripRequestId)&&(identical(other.planningRevision, planningRevision) || other.planningRevision == planningRevision)&&(identical(other.rank, rank) || other.rank == rank)&&(identical(other.status, status) || other.status == status)&&(identical(other.expiresAt, expiresAt) || other.expiresAt == expiresAt)&&const DeepCollectionEquality().equals(other.explanation, explanation)&&const DeepCollectionEquality().equals(other.items, items)&&(identical(other.totalPrice, totalPrice) || other.totalPrice == totalPrice));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,tripRequestId,planningRevision,rank,status,expiresAt,const DeepCollectionEquality().hash(explanation),const DeepCollectionEquality().hash(items),totalPrice);

@override
String toString() {
  return 'PlanOption(id: $id, tripRequestId: $tripRequestId, planningRevision: $planningRevision, rank: $rank, status: $status, expiresAt: $expiresAt, explanation: $explanation, items: $items, totalPrice: $totalPrice)';
}


}

/// @nodoc
abstract mixin class $PlanOptionCopyWith<$Res>  {
  factory $PlanOptionCopyWith(PlanOption value, $Res Function(PlanOption) _then) = _$PlanOptionCopyWithImpl;
@useResult
$Res call({
 String id,@JsonKey(name: 'trip_request_id') String tripRequestId,@JsonKey(name: 'planning_revision') int planningRevision, int rank, PlanOptionStatus status,@JsonKey(name: 'expires_at') DateTime expiresAt, List<String> explanation, List<PlanItem> items,@JsonKey(name: 'total_price') PriceSummary totalPrice
});


$PriceSummaryCopyWith<$Res> get totalPrice;

}
/// @nodoc
class _$PlanOptionCopyWithImpl<$Res>
    implements $PlanOptionCopyWith<$Res> {
  _$PlanOptionCopyWithImpl(this._self, this._then);

  final PlanOption _self;
  final $Res Function(PlanOption) _then;

/// Create a copy of PlanOption
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? tripRequestId = null,Object? planningRevision = null,Object? rank = null,Object? status = null,Object? expiresAt = null,Object? explanation = null,Object? items = null,Object? totalPrice = null,}) {
  return _then(PlanOption(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,tripRequestId: null == tripRequestId ? _self.tripRequestId : tripRequestId // ignore: cast_nullable_to_non_nullable
as String,planningRevision: null == planningRevision ? _self.planningRevision : planningRevision // ignore: cast_nullable_to_non_nullable
as int,rank: null == rank ? _self.rank : rank // ignore: cast_nullable_to_non_nullable
as int,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as PlanOptionStatus,expiresAt: null == expiresAt ? _self.expiresAt : expiresAt // ignore: cast_nullable_to_non_nullable
as DateTime,explanation: null == explanation ? _self.explanation : explanation // ignore: cast_nullable_to_non_nullable
as List<String>,items: null == items ? _self.items : items // ignore: cast_nullable_to_non_nullable
as List<PlanItem>,totalPrice: null == totalPrice ? _self.totalPrice : totalPrice // ignore: cast_nullable_to_non_nullable
as PriceSummary,
  ));
}
/// Create a copy of PlanOption
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$PriceSummaryCopyWith<$Res> get totalPrice {
  
  return $PriceSummaryCopyWith<$Res>(_self.totalPrice, (value) {
    return _then(_self.copyWith(totalPrice: value));
  });
}
}


/// Adds pattern-matching-related methods to [PlanOption].
extension PlanOptionPatterns on PlanOption {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _PlanOption value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _PlanOption() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _PlanOption value)  $default,){
final _that = this;
switch (_that) {
case _PlanOption():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _PlanOption value)?  $default,){
final _that = this;
switch (_that) {
case _PlanOption() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id, @JsonKey(name: 'trip_request_id')  String tripRequestId, @JsonKey(name: 'planning_revision')  int planningRevision,  int rank,  PlanOptionStatus status, @JsonKey(name: 'expires_at')  DateTime expiresAt,  List<String> explanation,  List<PlanItem> items, @JsonKey(name: 'total_price')  PriceSummary totalPrice)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _PlanOption() when $default != null:
return $default(_that.id,_that.tripRequestId,_that.planningRevision,_that.rank,_that.status,_that.expiresAt,_that.explanation,_that.items,_that.totalPrice);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id, @JsonKey(name: 'trip_request_id')  String tripRequestId, @JsonKey(name: 'planning_revision')  int planningRevision,  int rank,  PlanOptionStatus status, @JsonKey(name: 'expires_at')  DateTime expiresAt,  List<String> explanation,  List<PlanItem> items, @JsonKey(name: 'total_price')  PriceSummary totalPrice)  $default,) {final _that = this;
switch (_that) {
case _PlanOption():
return $default(_that.id,_that.tripRequestId,_that.planningRevision,_that.rank,_that.status,_that.expiresAt,_that.explanation,_that.items,_that.totalPrice);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id, @JsonKey(name: 'trip_request_id')  String tripRequestId, @JsonKey(name: 'planning_revision')  int planningRevision,  int rank,  PlanOptionStatus status, @JsonKey(name: 'expires_at')  DateTime expiresAt,  List<String> explanation,  List<PlanItem> items, @JsonKey(name: 'total_price')  PriceSummary totalPrice)?  $default,) {final _that = this;
switch (_that) {
case _PlanOption() when $default != null:
return $default(_that.id,_that.tripRequestId,_that.planningRevision,_that.rank,_that.status,_that.expiresAt,_that.explanation,_that.items,_that.totalPrice);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _PlanOption implements PlanOption {
  const _PlanOption({required this.id, @JsonKey(name: 'trip_request_id') required this.tripRequestId, @JsonKey(name: 'planning_revision') required this.planningRevision, required this.rank, required this.status, @JsonKey(name: 'expires_at') required this.expiresAt, required  List<String> explanation, required  List<PlanItem> items, @JsonKey(name: 'total_price') required this.totalPrice}): _explanation = explanation,_items = items;
  factory _PlanOption.fromJson(Map<String, dynamic> json) => _$PlanOptionFromJson(json);

@override final  String id;
@override@JsonKey(name: 'trip_request_id') final  String tripRequestId;
@override@JsonKey(name: 'planning_revision') final  int planningRevision;
@override final  int rank;
@override final  PlanOptionStatus status;
@override@JsonKey(name: 'expires_at') final  DateTime expiresAt;
 final  List<String> _explanation;
@override List<String> get explanation {
  if (_explanation is EqualUnmodifiableListView) return _explanation;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_explanation);
}

 final  List<PlanItem> _items;
@override List<PlanItem> get items {
  if (_items is EqualUnmodifiableListView) return _items;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_items);
}

@override@JsonKey(name: 'total_price') final  PriceSummary totalPrice;

/// Create a copy of PlanOption
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$PlanOptionCopyWith<_PlanOption> get copyWith => __$PlanOptionCopyWithImpl<_PlanOption>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$PlanOptionToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _PlanOption&&(identical(other.id, id) || other.id == id)&&(identical(other.tripRequestId, tripRequestId) || other.tripRequestId == tripRequestId)&&(identical(other.planningRevision, planningRevision) || other.planningRevision == planningRevision)&&(identical(other.rank, rank) || other.rank == rank)&&(identical(other.status, status) || other.status == status)&&(identical(other.expiresAt, expiresAt) || other.expiresAt == expiresAt)&&const DeepCollectionEquality().equals(other._explanation, _explanation)&&const DeepCollectionEquality().equals(other._items, _items)&&(identical(other.totalPrice, totalPrice) || other.totalPrice == totalPrice));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,tripRequestId,planningRevision,rank,status,expiresAt,const DeepCollectionEquality().hash(_explanation),const DeepCollectionEquality().hash(_items),totalPrice);

@override
String toString() {
  return 'PlanOption(id: $id, tripRequestId: $tripRequestId, planningRevision: $planningRevision, rank: $rank, status: $status, expiresAt: $expiresAt, explanation: $explanation, items: $items, totalPrice: $totalPrice)';
}


}

/// @nodoc
abstract mixin class _$PlanOptionCopyWith<$Res> implements $PlanOptionCopyWith<$Res> {
  factory _$PlanOptionCopyWith(_PlanOption value, $Res Function(_PlanOption) _then) = __$PlanOptionCopyWithImpl;
@override @useResult
$Res call({
 String id,@JsonKey(name: 'trip_request_id') String tripRequestId,@JsonKey(name: 'planning_revision') int planningRevision, int rank, PlanOptionStatus status,@JsonKey(name: 'expires_at') DateTime expiresAt, List<String> explanation, List<PlanItem> items,@JsonKey(name: 'total_price') PriceSummary totalPrice
});


@override $PriceSummaryCopyWith<$Res> get totalPrice;

}
/// @nodoc
class __$PlanOptionCopyWithImpl<$Res>
    implements _$PlanOptionCopyWith<$Res> {
  __$PlanOptionCopyWithImpl(this._self, this._then);

  final _PlanOption _self;
  final $Res Function(_PlanOption) _then;

/// Create a copy of PlanOption
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? tripRequestId = null,Object? planningRevision = null,Object? rank = null,Object? status = null,Object? expiresAt = null,Object? explanation = null,Object? items = null,Object? totalPrice = null,}) {
  return _then(_PlanOption(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,tripRequestId: null == tripRequestId ? _self.tripRequestId : tripRequestId // ignore: cast_nullable_to_non_nullable
as String,planningRevision: null == planningRevision ? _self.planningRevision : planningRevision // ignore: cast_nullable_to_non_nullable
as int,rank: null == rank ? _self.rank : rank // ignore: cast_nullable_to_non_nullable
as int,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as PlanOptionStatus,expiresAt: null == expiresAt ? _self.expiresAt : expiresAt // ignore: cast_nullable_to_non_nullable
as DateTime,explanation: null == explanation ? _self._explanation : explanation // ignore: cast_nullable_to_non_nullable
as List<String>,items: null == items ? _self._items : items // ignore: cast_nullable_to_non_nullable
as List<PlanItem>,totalPrice: null == totalPrice ? _self.totalPrice : totalPrice // ignore: cast_nullable_to_non_nullable
as PriceSummary,
  ));
}

/// Create a copy of PlanOption
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
mixin _$PlanningResult {

@JsonKey(name: 'trip_request') TripRequest get tripRequest; List<PlanOption> get options;@JsonKey(name: 'no_match_reasons') List<String> get noMatchReasons;@JsonKey(name: 'provider_warnings') List<String> get providerWarnings;
/// Create a copy of PlanningResult
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$PlanningResultCopyWith<PlanningResult> get copyWith => _$PlanningResultCopyWithImpl<PlanningResult>(this as PlanningResult, _$identity);

  /// Serializes this PlanningResult to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is PlanningResult&&(identical(other.tripRequest, tripRequest) || other.tripRequest == tripRequest)&&const DeepCollectionEquality().equals(other.options, options)&&const DeepCollectionEquality().equals(other.noMatchReasons, noMatchReasons)&&const DeepCollectionEquality().equals(other.providerWarnings, providerWarnings));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,tripRequest,const DeepCollectionEquality().hash(options),const DeepCollectionEquality().hash(noMatchReasons),const DeepCollectionEquality().hash(providerWarnings));

@override
String toString() {
  return 'PlanningResult(tripRequest: $tripRequest, options: $options, noMatchReasons: $noMatchReasons, providerWarnings: $providerWarnings)';
}


}

/// @nodoc
abstract mixin class $PlanningResultCopyWith<$Res>  {
  factory $PlanningResultCopyWith(PlanningResult value, $Res Function(PlanningResult) _then) = _$PlanningResultCopyWithImpl;
@useResult
$Res call({
@JsonKey(name: 'trip_request') TripRequest tripRequest, List<PlanOption> options,@JsonKey(name: 'no_match_reasons') List<String> noMatchReasons,@JsonKey(name: 'provider_warnings') List<String> providerWarnings
});


$TripRequestCopyWith<$Res> get tripRequest;

}
/// @nodoc
class _$PlanningResultCopyWithImpl<$Res>
    implements $PlanningResultCopyWith<$Res> {
  _$PlanningResultCopyWithImpl(this._self, this._then);

  final PlanningResult _self;
  final $Res Function(PlanningResult) _then;

/// Create a copy of PlanningResult
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? tripRequest = null,Object? options = null,Object? noMatchReasons = null,Object? providerWarnings = null,}) {
  return _then(PlanningResult(
tripRequest: null == tripRequest ? _self.tripRequest : tripRequest // ignore: cast_nullable_to_non_nullable
as TripRequest,options: null == options ? _self.options : options // ignore: cast_nullable_to_non_nullable
as List<PlanOption>,noMatchReasons: null == noMatchReasons ? _self.noMatchReasons : noMatchReasons // ignore: cast_nullable_to_non_nullable
as List<String>,providerWarnings: null == providerWarnings ? _self.providerWarnings : providerWarnings // ignore: cast_nullable_to_non_nullable
as List<String>,
  ));
}
/// Create a copy of PlanningResult
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$TripRequestCopyWith<$Res> get tripRequest {
  
  return $TripRequestCopyWith<$Res>(_self.tripRequest, (value) {
    return _then(_self.copyWith(tripRequest: value));
  });
}
}


/// Adds pattern-matching-related methods to [PlanningResult].
extension PlanningResultPatterns on PlanningResult {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _PlanningResult value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _PlanningResult() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _PlanningResult value)  $default,){
final _that = this;
switch (_that) {
case _PlanningResult():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _PlanningResult value)?  $default,){
final _that = this;
switch (_that) {
case _PlanningResult() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function(@JsonKey(name: 'trip_request')  TripRequest tripRequest,  List<PlanOption> options, @JsonKey(name: 'no_match_reasons')  List<String> noMatchReasons, @JsonKey(name: 'provider_warnings')  List<String> providerWarnings)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _PlanningResult() when $default != null:
return $default(_that.tripRequest,_that.options,_that.noMatchReasons,_that.providerWarnings);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function(@JsonKey(name: 'trip_request')  TripRequest tripRequest,  List<PlanOption> options, @JsonKey(name: 'no_match_reasons')  List<String> noMatchReasons, @JsonKey(name: 'provider_warnings')  List<String> providerWarnings)  $default,) {final _that = this;
switch (_that) {
case _PlanningResult():
return $default(_that.tripRequest,_that.options,_that.noMatchReasons,_that.providerWarnings);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function(@JsonKey(name: 'trip_request')  TripRequest tripRequest,  List<PlanOption> options, @JsonKey(name: 'no_match_reasons')  List<String> noMatchReasons, @JsonKey(name: 'provider_warnings')  List<String> providerWarnings)?  $default,) {final _that = this;
switch (_that) {
case _PlanningResult() when $default != null:
return $default(_that.tripRequest,_that.options,_that.noMatchReasons,_that.providerWarnings);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _PlanningResult implements PlanningResult {
  const _PlanningResult({@JsonKey(name: 'trip_request') required this.tripRequest, required  List<PlanOption> options, @JsonKey(name: 'no_match_reasons') required  List<String> noMatchReasons, @JsonKey(name: 'provider_warnings') required  List<String> providerWarnings}): _options = options,_noMatchReasons = noMatchReasons,_providerWarnings = providerWarnings;
  factory _PlanningResult.fromJson(Map<String, dynamic> json) => _$PlanningResultFromJson(json);

@override@JsonKey(name: 'trip_request') final  TripRequest tripRequest;
 final  List<PlanOption> _options;
@override List<PlanOption> get options {
  if (_options is EqualUnmodifiableListView) return _options;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_options);
}

 final  List<String> _noMatchReasons;
@override@JsonKey(name: 'no_match_reasons') List<String> get noMatchReasons {
  if (_noMatchReasons is EqualUnmodifiableListView) return _noMatchReasons;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_noMatchReasons);
}

 final  List<String> _providerWarnings;
@override@JsonKey(name: 'provider_warnings') List<String> get providerWarnings {
  if (_providerWarnings is EqualUnmodifiableListView) return _providerWarnings;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_providerWarnings);
}


/// Create a copy of PlanningResult
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$PlanningResultCopyWith<_PlanningResult> get copyWith => __$PlanningResultCopyWithImpl<_PlanningResult>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$PlanningResultToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _PlanningResult&&(identical(other.tripRequest, tripRequest) || other.tripRequest == tripRequest)&&const DeepCollectionEquality().equals(other._options, _options)&&const DeepCollectionEquality().equals(other._noMatchReasons, _noMatchReasons)&&const DeepCollectionEquality().equals(other._providerWarnings, _providerWarnings));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,tripRequest,const DeepCollectionEquality().hash(_options),const DeepCollectionEquality().hash(_noMatchReasons),const DeepCollectionEquality().hash(_providerWarnings));

@override
String toString() {
  return 'PlanningResult(tripRequest: $tripRequest, options: $options, noMatchReasons: $noMatchReasons, providerWarnings: $providerWarnings)';
}


}

/// @nodoc
abstract mixin class _$PlanningResultCopyWith<$Res> implements $PlanningResultCopyWith<$Res> {
  factory _$PlanningResultCopyWith(_PlanningResult value, $Res Function(_PlanningResult) _then) = __$PlanningResultCopyWithImpl;
@override @useResult
$Res call({
@JsonKey(name: 'trip_request') TripRequest tripRequest, List<PlanOption> options,@JsonKey(name: 'no_match_reasons') List<String> noMatchReasons,@JsonKey(name: 'provider_warnings') List<String> providerWarnings
});


@override $TripRequestCopyWith<$Res> get tripRequest;

}
/// @nodoc
class __$PlanningResultCopyWithImpl<$Res>
    implements _$PlanningResultCopyWith<$Res> {
  __$PlanningResultCopyWithImpl(this._self, this._then);

  final _PlanningResult _self;
  final $Res Function(_PlanningResult) _then;

/// Create a copy of PlanningResult
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? tripRequest = null,Object? options = null,Object? noMatchReasons = null,Object? providerWarnings = null,}) {
  return _then(_PlanningResult(
tripRequest: null == tripRequest ? _self.tripRequest : tripRequest // ignore: cast_nullable_to_non_nullable
as TripRequest,options: null == options ? _self._options : options // ignore: cast_nullable_to_non_nullable
as List<PlanOption>,noMatchReasons: null == noMatchReasons ? _self._noMatchReasons : noMatchReasons // ignore: cast_nullable_to_non_nullable
as List<String>,providerWarnings: null == providerWarnings ? _self._providerWarnings : providerWarnings // ignore: cast_nullable_to_non_nullable
as List<String>,
  ));
}

/// Create a copy of PlanningResult
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$TripRequestCopyWith<$Res> get tripRequest {
  
  return $TripRequestCopyWith<$Res>(_self.tripRequest, (value) {
    return _then(_self.copyWith(tripRequest: value));
  });
}
}


/// @nodoc
mixin _$TripRequestDetail {

@JsonKey(name: 'trip_request') TripRequest get tripRequest;@JsonKey(name: 'plan_options') List<PlanOption> get planOptions;
/// Create a copy of TripRequestDetail
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$TripRequestDetailCopyWith<TripRequestDetail> get copyWith => _$TripRequestDetailCopyWithImpl<TripRequestDetail>(this as TripRequestDetail, _$identity);

  /// Serializes this TripRequestDetail to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is TripRequestDetail&&(identical(other.tripRequest, tripRequest) || other.tripRequest == tripRequest)&&const DeepCollectionEquality().equals(other.planOptions, planOptions));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,tripRequest,const DeepCollectionEquality().hash(planOptions));

@override
String toString() {
  return 'TripRequestDetail(tripRequest: $tripRequest, planOptions: $planOptions)';
}


}

/// @nodoc
abstract mixin class $TripRequestDetailCopyWith<$Res>  {
  factory $TripRequestDetailCopyWith(TripRequestDetail value, $Res Function(TripRequestDetail) _then) = _$TripRequestDetailCopyWithImpl;
@useResult
$Res call({
@JsonKey(name: 'trip_request') TripRequest tripRequest,@JsonKey(name: 'plan_options') List<PlanOption> planOptions
});


$TripRequestCopyWith<$Res> get tripRequest;

}
/// @nodoc
class _$TripRequestDetailCopyWithImpl<$Res>
    implements $TripRequestDetailCopyWith<$Res> {
  _$TripRequestDetailCopyWithImpl(this._self, this._then);

  final TripRequestDetail _self;
  final $Res Function(TripRequestDetail) _then;

/// Create a copy of TripRequestDetail
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? tripRequest = null,Object? planOptions = null,}) {
  return _then(TripRequestDetail(
tripRequest: null == tripRequest ? _self.tripRequest : tripRequest // ignore: cast_nullable_to_non_nullable
as TripRequest,planOptions: null == planOptions ? _self.planOptions : planOptions // ignore: cast_nullable_to_non_nullable
as List<PlanOption>,
  ));
}
/// Create a copy of TripRequestDetail
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$TripRequestCopyWith<$Res> get tripRequest {
  
  return $TripRequestCopyWith<$Res>(_self.tripRequest, (value) {
    return _then(_self.copyWith(tripRequest: value));
  });
}
}


/// Adds pattern-matching-related methods to [TripRequestDetail].
extension TripRequestDetailPatterns on TripRequestDetail {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _TripRequestDetail value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _TripRequestDetail() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _TripRequestDetail value)  $default,){
final _that = this;
switch (_that) {
case _TripRequestDetail():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _TripRequestDetail value)?  $default,){
final _that = this;
switch (_that) {
case _TripRequestDetail() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function(@JsonKey(name: 'trip_request')  TripRequest tripRequest, @JsonKey(name: 'plan_options')  List<PlanOption> planOptions)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _TripRequestDetail() when $default != null:
return $default(_that.tripRequest,_that.planOptions);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function(@JsonKey(name: 'trip_request')  TripRequest tripRequest, @JsonKey(name: 'plan_options')  List<PlanOption> planOptions)  $default,) {final _that = this;
switch (_that) {
case _TripRequestDetail():
return $default(_that.tripRequest,_that.planOptions);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function(@JsonKey(name: 'trip_request')  TripRequest tripRequest, @JsonKey(name: 'plan_options')  List<PlanOption> planOptions)?  $default,) {final _that = this;
switch (_that) {
case _TripRequestDetail() when $default != null:
return $default(_that.tripRequest,_that.planOptions);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _TripRequestDetail implements TripRequestDetail {
  const _TripRequestDetail({@JsonKey(name: 'trip_request') required this.tripRequest, @JsonKey(name: 'plan_options') required  List<PlanOption> planOptions}): _planOptions = planOptions;
  factory _TripRequestDetail.fromJson(Map<String, dynamic> json) => _$TripRequestDetailFromJson(json);

@override@JsonKey(name: 'trip_request') final  TripRequest tripRequest;
 final  List<PlanOption> _planOptions;
@override@JsonKey(name: 'plan_options') List<PlanOption> get planOptions {
  if (_planOptions is EqualUnmodifiableListView) return _planOptions;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_planOptions);
}


/// Create a copy of TripRequestDetail
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$TripRequestDetailCopyWith<_TripRequestDetail> get copyWith => __$TripRequestDetailCopyWithImpl<_TripRequestDetail>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$TripRequestDetailToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _TripRequestDetail&&(identical(other.tripRequest, tripRequest) || other.tripRequest == tripRequest)&&const DeepCollectionEquality().equals(other._planOptions, _planOptions));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,tripRequest,const DeepCollectionEquality().hash(_planOptions));

@override
String toString() {
  return 'TripRequestDetail(tripRequest: $tripRequest, planOptions: $planOptions)';
}


}

/// @nodoc
abstract mixin class _$TripRequestDetailCopyWith<$Res> implements $TripRequestDetailCopyWith<$Res> {
  factory _$TripRequestDetailCopyWith(_TripRequestDetail value, $Res Function(_TripRequestDetail) _then) = __$TripRequestDetailCopyWithImpl;
@override @useResult
$Res call({
@JsonKey(name: 'trip_request') TripRequest tripRequest,@JsonKey(name: 'plan_options') List<PlanOption> planOptions
});


@override $TripRequestCopyWith<$Res> get tripRequest;

}
/// @nodoc
class __$TripRequestDetailCopyWithImpl<$Res>
    implements _$TripRequestDetailCopyWith<$Res> {
  __$TripRequestDetailCopyWithImpl(this._self, this._then);

  final _TripRequestDetail _self;
  final $Res Function(_TripRequestDetail) _then;

/// Create a copy of TripRequestDetail
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? tripRequest = null,Object? planOptions = null,}) {
  return _then(_TripRequestDetail(
tripRequest: null == tripRequest ? _self.tripRequest : tripRequest // ignore: cast_nullable_to_non_nullable
as TripRequest,planOptions: null == planOptions ? _self._planOptions : planOptions // ignore: cast_nullable_to_non_nullable
as List<PlanOption>,
  ));
}

/// Create a copy of TripRequestDetail
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$TripRequestCopyWith<$Res> get tripRequest {
  
  return $TripRequestCopyWith<$Res>(_self.tripRequest, (value) {
    return _then(_self.copyWith(tripRequest: value));
  });
}
}

// dart format on
