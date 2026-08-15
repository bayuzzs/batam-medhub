// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint, type=warning, deprecated_member_use, deprecated_member_use_from_same_package
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'medical_service.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// dart format off
T _$identity<T>(T value) => value;

/// @nodoc
mixin _$MedicalService {

 String get code; String get name; String get category; String? get description; bool get active; bool get synthetic; String get source;
/// Create a copy of MedicalService
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$MedicalServiceCopyWith<MedicalService> get copyWith => _$MedicalServiceCopyWithImpl<MedicalService>(this as MedicalService, _$identity);

  /// Serializes this MedicalService to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is MedicalService&&(identical(other.code, code) || other.code == code)&&(identical(other.name, name) || other.name == name)&&(identical(other.category, category) || other.category == category)&&(identical(other.description, description) || other.description == description)&&(identical(other.active, active) || other.active == active)&&(identical(other.synthetic, synthetic) || other.synthetic == synthetic)&&(identical(other.source, source) || other.source == source));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,code,name,category,description,active,synthetic,source);

@override
String toString() {
  return 'MedicalService(code: $code, name: $name, category: $category, description: $description, active: $active, synthetic: $synthetic, source: $source)';
}


}

/// @nodoc
abstract mixin class $MedicalServiceCopyWith<$Res>  {
  factory $MedicalServiceCopyWith(MedicalService value, $Res Function(MedicalService) _then) = _$MedicalServiceCopyWithImpl;
@useResult
$Res call({
 String code, String name, String category, String? description, bool active, bool synthetic, String source
});




}
/// @nodoc
class _$MedicalServiceCopyWithImpl<$Res>
    implements $MedicalServiceCopyWith<$Res> {
  _$MedicalServiceCopyWithImpl(this._self, this._then);

  final MedicalService _self;
  final $Res Function(MedicalService) _then;

/// Create a copy of MedicalService
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? code = null,Object? name = null,Object? category = null,Object? description = freezed,Object? active = null,Object? synthetic = null,Object? source = null,}) {
  return _then(MedicalService(
code: null == code ? _self.code : code // ignore: cast_nullable_to_non_nullable
as String,name: null == name ? _self.name : name // ignore: cast_nullable_to_non_nullable
as String,category: null == category ? _self.category : category // ignore: cast_nullable_to_non_nullable
as String,description: freezed == description ? _self.description : description // ignore: cast_nullable_to_non_nullable
as String?,active: null == active ? _self.active : active // ignore: cast_nullable_to_non_nullable
as bool,synthetic: null == synthetic ? _self.synthetic : synthetic // ignore: cast_nullable_to_non_nullable
as bool,source: null == source ? _self.source : source // ignore: cast_nullable_to_non_nullable
as String,
  ));
}

}


/// Adds pattern-matching-related methods to [MedicalService].
extension MedicalServicePatterns on MedicalService {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _MedicalService value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _MedicalService() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _MedicalService value)  $default,){
final _that = this;
switch (_that) {
case _MedicalService():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _MedicalService value)?  $default,){
final _that = this;
switch (_that) {
case _MedicalService() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String code,  String name,  String category,  String? description,  bool active,  bool synthetic,  String source)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _MedicalService() when $default != null:
return $default(_that.code,_that.name,_that.category,_that.description,_that.active,_that.synthetic,_that.source);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String code,  String name,  String category,  String? description,  bool active,  bool synthetic,  String source)  $default,) {final _that = this;
switch (_that) {
case _MedicalService():
return $default(_that.code,_that.name,_that.category,_that.description,_that.active,_that.synthetic,_that.source);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String code,  String name,  String category,  String? description,  bool active,  bool synthetic,  String source)?  $default,) {final _that = this;
switch (_that) {
case _MedicalService() when $default != null:
return $default(_that.code,_that.name,_that.category,_that.description,_that.active,_that.synthetic,_that.source);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _MedicalService implements MedicalService {
  const _MedicalService({required this.code, required this.name, required this.category, this.description, required this.active, required this.synthetic, required this.source});
  factory _MedicalService.fromJson(Map<String, dynamic> json) => _$MedicalServiceFromJson(json);

@override final  String code;
@override final  String name;
@override final  String category;
@override final  String? description;
@override final  bool active;
@override final  bool synthetic;
@override final  String source;

/// Create a copy of MedicalService
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$MedicalServiceCopyWith<_MedicalService> get copyWith => __$MedicalServiceCopyWithImpl<_MedicalService>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$MedicalServiceToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _MedicalService&&(identical(other.code, code) || other.code == code)&&(identical(other.name, name) || other.name == name)&&(identical(other.category, category) || other.category == category)&&(identical(other.description, description) || other.description == description)&&(identical(other.active, active) || other.active == active)&&(identical(other.synthetic, synthetic) || other.synthetic == synthetic)&&(identical(other.source, source) || other.source == source));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,code,name,category,description,active,synthetic,source);

@override
String toString() {
  return 'MedicalService(code: $code, name: $name, category: $category, description: $description, active: $active, synthetic: $synthetic, source: $source)';
}


}

/// @nodoc
abstract mixin class _$MedicalServiceCopyWith<$Res> implements $MedicalServiceCopyWith<$Res> {
  factory _$MedicalServiceCopyWith(_MedicalService value, $Res Function(_MedicalService) _then) = __$MedicalServiceCopyWithImpl;
@override @useResult
$Res call({
 String code, String name, String category, String? description, bool active, bool synthetic, String source
});




}
/// @nodoc
class __$MedicalServiceCopyWithImpl<$Res>
    implements _$MedicalServiceCopyWith<$Res> {
  __$MedicalServiceCopyWithImpl(this._self, this._then);

  final _MedicalService _self;
  final $Res Function(_MedicalService) _then;

/// Create a copy of MedicalService
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? code = null,Object? name = null,Object? category = null,Object? description = freezed,Object? active = null,Object? synthetic = null,Object? source = null,}) {
  return _then(_MedicalService(
code: null == code ? _self.code : code // ignore: cast_nullable_to_non_nullable
as String,name: null == name ? _self.name : name // ignore: cast_nullable_to_non_nullable
as String,category: null == category ? _self.category : category // ignore: cast_nullable_to_non_nullable
as String,description: freezed == description ? _self.description : description // ignore: cast_nullable_to_non_nullable
as String?,active: null == active ? _self.active : active // ignore: cast_nullable_to_non_nullable
as bool,synthetic: null == synthetic ? _self.synthetic : synthetic // ignore: cast_nullable_to_non_nullable
as bool,source: null == source ? _self.source : source // ignore: cast_nullable_to_non_nullable
as String,
  ));
}


}


/// @nodoc
mixin _$MedicalServiceListResponse {

 List<MedicalService> get services;
/// Create a copy of MedicalServiceListResponse
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$MedicalServiceListResponseCopyWith<MedicalServiceListResponse> get copyWith => _$MedicalServiceListResponseCopyWithImpl<MedicalServiceListResponse>(this as MedicalServiceListResponse, _$identity);

  /// Serializes this MedicalServiceListResponse to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is MedicalServiceListResponse&&const DeepCollectionEquality().equals(other.services, services));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,const DeepCollectionEquality().hash(services));

@override
String toString() {
  return 'MedicalServiceListResponse(services: $services)';
}


}

/// @nodoc
abstract mixin class $MedicalServiceListResponseCopyWith<$Res>  {
  factory $MedicalServiceListResponseCopyWith(MedicalServiceListResponse value, $Res Function(MedicalServiceListResponse) _then) = _$MedicalServiceListResponseCopyWithImpl;
@useResult
$Res call({
 List<MedicalService> services
});




}
/// @nodoc
class _$MedicalServiceListResponseCopyWithImpl<$Res>
    implements $MedicalServiceListResponseCopyWith<$Res> {
  _$MedicalServiceListResponseCopyWithImpl(this._self, this._then);

  final MedicalServiceListResponse _self;
  final $Res Function(MedicalServiceListResponse) _then;

/// Create a copy of MedicalServiceListResponse
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? services = null,}) {
  return _then(MedicalServiceListResponse(
services: null == services ? _self.services : services // ignore: cast_nullable_to_non_nullable
as List<MedicalService>,
  ));
}

}


/// Adds pattern-matching-related methods to [MedicalServiceListResponse].
extension MedicalServiceListResponsePatterns on MedicalServiceListResponse {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _MedicalServiceListResponse value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _MedicalServiceListResponse() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _MedicalServiceListResponse value)  $default,){
final _that = this;
switch (_that) {
case _MedicalServiceListResponse():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _MedicalServiceListResponse value)?  $default,){
final _that = this;
switch (_that) {
case _MedicalServiceListResponse() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( List<MedicalService> services)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _MedicalServiceListResponse() when $default != null:
return $default(_that.services);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( List<MedicalService> services)  $default,) {final _that = this;
switch (_that) {
case _MedicalServiceListResponse():
return $default(_that.services);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( List<MedicalService> services)?  $default,) {final _that = this;
switch (_that) {
case _MedicalServiceListResponse() when $default != null:
return $default(_that.services);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _MedicalServiceListResponse implements MedicalServiceListResponse {
  const _MedicalServiceListResponse({required  List<MedicalService> services}): _services = services;
  factory _MedicalServiceListResponse.fromJson(Map<String, dynamic> json) => _$MedicalServiceListResponseFromJson(json);

 final  List<MedicalService> _services;
@override List<MedicalService> get services {
  if (_services is EqualUnmodifiableListView) return _services;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_services);
}


/// Create a copy of MedicalServiceListResponse
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$MedicalServiceListResponseCopyWith<_MedicalServiceListResponse> get copyWith => __$MedicalServiceListResponseCopyWithImpl<_MedicalServiceListResponse>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$MedicalServiceListResponseToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _MedicalServiceListResponse&&const DeepCollectionEquality().equals(other._services, _services));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,const DeepCollectionEquality().hash(_services));

@override
String toString() {
  return 'MedicalServiceListResponse(services: $services)';
}


}

/// @nodoc
abstract mixin class _$MedicalServiceListResponseCopyWith<$Res> implements $MedicalServiceListResponseCopyWith<$Res> {
  factory _$MedicalServiceListResponseCopyWith(_MedicalServiceListResponse value, $Res Function(_MedicalServiceListResponse) _then) = __$MedicalServiceListResponseCopyWithImpl;
@override @useResult
$Res call({
 List<MedicalService> services
});




}
/// @nodoc
class __$MedicalServiceListResponseCopyWithImpl<$Res>
    implements _$MedicalServiceListResponseCopyWith<$Res> {
  __$MedicalServiceListResponseCopyWithImpl(this._self, this._then);

  final _MedicalServiceListResponse _self;
  final $Res Function(_MedicalServiceListResponse) _then;

/// Create a copy of MedicalServiceListResponse
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? services = null,}) {
  return _then(_MedicalServiceListResponse(
services: null == services ? _self._services : services // ignore: cast_nullable_to_non_nullable
as List<MedicalService>,
  ));
}


}

// dart format on
