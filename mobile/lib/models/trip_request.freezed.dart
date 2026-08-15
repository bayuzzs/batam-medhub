// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint, type=warning, deprecated_member_use, deprecated_member_use_from_same_package
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'trip_request.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// dart format off
T _$identity<T>(T value) => value;

/// @nodoc
mixin _$TripRequest {

 String get id; TripRequestStatus get status; StructuredIntent get intent;@JsonKey(name: 'planning_revision') int get planningRevision;@JsonKey(name: 'reference_currency') String get referenceCurrency;@JsonKey(name: 'journey_id') String? get journeyId;@JsonKey(name: 'created_at') DateTime get createdAt;@JsonKey(name: 'updated_at') DateTime get updatedAt;
/// Create a copy of TripRequest
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$TripRequestCopyWith<TripRequest> get copyWith => _$TripRequestCopyWithImpl<TripRequest>(this as TripRequest, _$identity);

  /// Serializes this TripRequest to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is TripRequest&&(identical(other.id, id) || other.id == id)&&(identical(other.status, status) || other.status == status)&&(identical(other.intent, intent) || other.intent == intent)&&(identical(other.planningRevision, planningRevision) || other.planningRevision == planningRevision)&&(identical(other.referenceCurrency, referenceCurrency) || other.referenceCurrency == referenceCurrency)&&(identical(other.journeyId, journeyId) || other.journeyId == journeyId)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt)&&(identical(other.updatedAt, updatedAt) || other.updatedAt == updatedAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,status,intent,planningRevision,referenceCurrency,journeyId,createdAt,updatedAt);

@override
String toString() {
  return 'TripRequest(id: $id, status: $status, intent: $intent, planningRevision: $planningRevision, referenceCurrency: $referenceCurrency, journeyId: $journeyId, createdAt: $createdAt, updatedAt: $updatedAt)';
}


}

/// @nodoc
abstract mixin class $TripRequestCopyWith<$Res>  {
  factory $TripRequestCopyWith(TripRequest value, $Res Function(TripRequest) _then) = _$TripRequestCopyWithImpl;
@useResult
$Res call({
 String id, TripRequestStatus status, StructuredIntent intent,@JsonKey(name: 'planning_revision') int planningRevision,@JsonKey(name: 'reference_currency') String referenceCurrency,@JsonKey(name: 'journey_id') String? journeyId,@JsonKey(name: 'created_at') DateTime createdAt,@JsonKey(name: 'updated_at') DateTime updatedAt
});


$StructuredIntentCopyWith<$Res> get intent;

}
/// @nodoc
class _$TripRequestCopyWithImpl<$Res>
    implements $TripRequestCopyWith<$Res> {
  _$TripRequestCopyWithImpl(this._self, this._then);

  final TripRequest _self;
  final $Res Function(TripRequest) _then;

/// Create a copy of TripRequest
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? status = null,Object? intent = null,Object? planningRevision = null,Object? referenceCurrency = null,Object? journeyId = freezed,Object? createdAt = null,Object? updatedAt = null,}) {
  return _then(TripRequest(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as TripRequestStatus,intent: null == intent ? _self.intent : intent // ignore: cast_nullable_to_non_nullable
as StructuredIntent,planningRevision: null == planningRevision ? _self.planningRevision : planningRevision // ignore: cast_nullable_to_non_nullable
as int,referenceCurrency: null == referenceCurrency ? _self.referenceCurrency : referenceCurrency // ignore: cast_nullable_to_non_nullable
as String,journeyId: freezed == journeyId ? _self.journeyId : journeyId // ignore: cast_nullable_to_non_nullable
as String?,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,updatedAt: null == updatedAt ? _self.updatedAt : updatedAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}
/// Create a copy of TripRequest
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$StructuredIntentCopyWith<$Res> get intent {
  
  return $StructuredIntentCopyWith<$Res>(_self.intent, (value) {
    return _then(_self.copyWith(intent: value));
  });
}
}


/// Adds pattern-matching-related methods to [TripRequest].
extension TripRequestPatterns on TripRequest {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _TripRequest value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _TripRequest() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _TripRequest value)  $default,){
final _that = this;
switch (_that) {
case _TripRequest():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _TripRequest value)?  $default,){
final _that = this;
switch (_that) {
case _TripRequest() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  TripRequestStatus status,  StructuredIntent intent, @JsonKey(name: 'planning_revision')  int planningRevision, @JsonKey(name: 'reference_currency')  String referenceCurrency, @JsonKey(name: 'journey_id')  String? journeyId, @JsonKey(name: 'created_at')  DateTime createdAt, @JsonKey(name: 'updated_at')  DateTime updatedAt)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _TripRequest() when $default != null:
return $default(_that.id,_that.status,_that.intent,_that.planningRevision,_that.referenceCurrency,_that.journeyId,_that.createdAt,_that.updatedAt);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  TripRequestStatus status,  StructuredIntent intent, @JsonKey(name: 'planning_revision')  int planningRevision, @JsonKey(name: 'reference_currency')  String referenceCurrency, @JsonKey(name: 'journey_id')  String? journeyId, @JsonKey(name: 'created_at')  DateTime createdAt, @JsonKey(name: 'updated_at')  DateTime updatedAt)  $default,) {final _that = this;
switch (_that) {
case _TripRequest():
return $default(_that.id,_that.status,_that.intent,_that.planningRevision,_that.referenceCurrency,_that.journeyId,_that.createdAt,_that.updatedAt);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  TripRequestStatus status,  StructuredIntent intent, @JsonKey(name: 'planning_revision')  int planningRevision, @JsonKey(name: 'reference_currency')  String referenceCurrency, @JsonKey(name: 'journey_id')  String? journeyId, @JsonKey(name: 'created_at')  DateTime createdAt, @JsonKey(name: 'updated_at')  DateTime updatedAt)?  $default,) {final _that = this;
switch (_that) {
case _TripRequest() when $default != null:
return $default(_that.id,_that.status,_that.intent,_that.planningRevision,_that.referenceCurrency,_that.journeyId,_that.createdAt,_that.updatedAt);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _TripRequest implements TripRequest {
  const _TripRequest({required this.id, required this.status, required this.intent, @JsonKey(name: 'planning_revision') required this.planningRevision, @JsonKey(name: 'reference_currency') required this.referenceCurrency, @JsonKey(name: 'journey_id') this.journeyId, @JsonKey(name: 'created_at') required this.createdAt, @JsonKey(name: 'updated_at') required this.updatedAt});
  factory _TripRequest.fromJson(Map<String, dynamic> json) => _$TripRequestFromJson(json);

@override final  String id;
@override final  TripRequestStatus status;
@override final  StructuredIntent intent;
@override@JsonKey(name: 'planning_revision') final  int planningRevision;
@override@JsonKey(name: 'reference_currency') final  String referenceCurrency;
@override@JsonKey(name: 'journey_id') final  String? journeyId;
@override@JsonKey(name: 'created_at') final  DateTime createdAt;
@override@JsonKey(name: 'updated_at') final  DateTime updatedAt;

/// Create a copy of TripRequest
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$TripRequestCopyWith<_TripRequest> get copyWith => __$TripRequestCopyWithImpl<_TripRequest>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$TripRequestToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _TripRequest&&(identical(other.id, id) || other.id == id)&&(identical(other.status, status) || other.status == status)&&(identical(other.intent, intent) || other.intent == intent)&&(identical(other.planningRevision, planningRevision) || other.planningRevision == planningRevision)&&(identical(other.referenceCurrency, referenceCurrency) || other.referenceCurrency == referenceCurrency)&&(identical(other.journeyId, journeyId) || other.journeyId == journeyId)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt)&&(identical(other.updatedAt, updatedAt) || other.updatedAt == updatedAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,status,intent,planningRevision,referenceCurrency,journeyId,createdAt,updatedAt);

@override
String toString() {
  return 'TripRequest(id: $id, status: $status, intent: $intent, planningRevision: $planningRevision, referenceCurrency: $referenceCurrency, journeyId: $journeyId, createdAt: $createdAt, updatedAt: $updatedAt)';
}


}

/// @nodoc
abstract mixin class _$TripRequestCopyWith<$Res> implements $TripRequestCopyWith<$Res> {
  factory _$TripRequestCopyWith(_TripRequest value, $Res Function(_TripRequest) _then) = __$TripRequestCopyWithImpl;
@override @useResult
$Res call({
 String id, TripRequestStatus status, StructuredIntent intent,@JsonKey(name: 'planning_revision') int planningRevision,@JsonKey(name: 'reference_currency') String referenceCurrency,@JsonKey(name: 'journey_id') String? journeyId,@JsonKey(name: 'created_at') DateTime createdAt,@JsonKey(name: 'updated_at') DateTime updatedAt
});


@override $StructuredIntentCopyWith<$Res> get intent;

}
/// @nodoc
class __$TripRequestCopyWithImpl<$Res>
    implements _$TripRequestCopyWith<$Res> {
  __$TripRequestCopyWithImpl(this._self, this._then);

  final _TripRequest _self;
  final $Res Function(_TripRequest) _then;

/// Create a copy of TripRequest
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? status = null,Object? intent = null,Object? planningRevision = null,Object? referenceCurrency = null,Object? journeyId = freezed,Object? createdAt = null,Object? updatedAt = null,}) {
  return _then(_TripRequest(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,status: null == status ? _self.status : status // ignore: cast_nullable_to_non_nullable
as TripRequestStatus,intent: null == intent ? _self.intent : intent // ignore: cast_nullable_to_non_nullable
as StructuredIntent,planningRevision: null == planningRevision ? _self.planningRevision : planningRevision // ignore: cast_nullable_to_non_nullable
as int,referenceCurrency: null == referenceCurrency ? _self.referenceCurrency : referenceCurrency // ignore: cast_nullable_to_non_nullable
as String,journeyId: freezed == journeyId ? _self.journeyId : journeyId // ignore: cast_nullable_to_non_nullable
as String?,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,updatedAt: null == updatedAt ? _self.updatedAt : updatedAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}

/// Create a copy of TripRequest
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$StructuredIntentCopyWith<$Res> get intent {
  
  return $StructuredIntentCopyWith<$Res>(_self.intent, (value) {
    return _then(_self.copyWith(intent: value));
  });
}
}


/// @nodoc
mixin _$CreateTripRequest {

 String get prompt; String get locale;
/// Create a copy of CreateTripRequest
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$CreateTripRequestCopyWith<CreateTripRequest> get copyWith => _$CreateTripRequestCopyWithImpl<CreateTripRequest>(this as CreateTripRequest, _$identity);

  /// Serializes this CreateTripRequest to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is CreateTripRequest&&(identical(other.prompt, prompt) || other.prompt == prompt)&&(identical(other.locale, locale) || other.locale == locale));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,prompt,locale);

@override
String toString() {
  return 'CreateTripRequest(prompt: $prompt, locale: $locale)';
}


}

/// @nodoc
abstract mixin class $CreateTripRequestCopyWith<$Res>  {
  factory $CreateTripRequestCopyWith(CreateTripRequest value, $Res Function(CreateTripRequest) _then) = _$CreateTripRequestCopyWithImpl;
@useResult
$Res call({
 String prompt, String locale
});




}
/// @nodoc
class _$CreateTripRequestCopyWithImpl<$Res>
    implements $CreateTripRequestCopyWith<$Res> {
  _$CreateTripRequestCopyWithImpl(this._self, this._then);

  final CreateTripRequest _self;
  final $Res Function(CreateTripRequest) _then;

/// Create a copy of CreateTripRequest
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? prompt = null,Object? locale = null,}) {
  return _then(CreateTripRequest(
prompt: null == prompt ? _self.prompt : prompt // ignore: cast_nullable_to_non_nullable
as String,locale: null == locale ? _self.locale : locale // ignore: cast_nullable_to_non_nullable
as String,
  ));
}

}


/// Adds pattern-matching-related methods to [CreateTripRequest].
extension CreateTripRequestPatterns on CreateTripRequest {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _CreateTripRequest value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _CreateTripRequest() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _CreateTripRequest value)  $default,){
final _that = this;
switch (_that) {
case _CreateTripRequest():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _CreateTripRequest value)?  $default,){
final _that = this;
switch (_that) {
case _CreateTripRequest() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String prompt,  String locale)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _CreateTripRequest() when $default != null:
return $default(_that.prompt,_that.locale);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String prompt,  String locale)  $default,) {final _that = this;
switch (_that) {
case _CreateTripRequest():
return $default(_that.prompt,_that.locale);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String prompt,  String locale)?  $default,) {final _that = this;
switch (_that) {
case _CreateTripRequest() when $default != null:
return $default(_that.prompt,_that.locale);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _CreateTripRequest implements CreateTripRequest {
  const _CreateTripRequest({required this.prompt, required this.locale});
  factory _CreateTripRequest.fromJson(Map<String, dynamic> json) => _$CreateTripRequestFromJson(json);

@override final  String prompt;
@override final  String locale;

/// Create a copy of CreateTripRequest
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$CreateTripRequestCopyWith<_CreateTripRequest> get copyWith => __$CreateTripRequestCopyWithImpl<_CreateTripRequest>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$CreateTripRequestToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _CreateTripRequest&&(identical(other.prompt, prompt) || other.prompt == prompt)&&(identical(other.locale, locale) || other.locale == locale));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,prompt,locale);

@override
String toString() {
  return 'CreateTripRequest(prompt: $prompt, locale: $locale)';
}


}

/// @nodoc
abstract mixin class _$CreateTripRequestCopyWith<$Res> implements $CreateTripRequestCopyWith<$Res> {
  factory _$CreateTripRequestCopyWith(_CreateTripRequest value, $Res Function(_CreateTripRequest) _then) = __$CreateTripRequestCopyWithImpl;
@override @useResult
$Res call({
 String prompt, String locale
});




}
/// @nodoc
class __$CreateTripRequestCopyWithImpl<$Res>
    implements _$CreateTripRequestCopyWith<$Res> {
  __$CreateTripRequestCopyWithImpl(this._self, this._then);

  final _CreateTripRequest _self;
  final $Res Function(_CreateTripRequest) _then;

/// Create a copy of CreateTripRequest
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? prompt = null,Object? locale = null,}) {
  return _then(_CreateTripRequest(
prompt: null == prompt ? _self.prompt : prompt // ignore: cast_nullable_to_non_nullable
as String,locale: null == locale ? _self.locale : locale // ignore: cast_nullable_to_non_nullable
as String,
  ));
}


}

// dart format on
