// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint, type=warning, deprecated_member_use, deprecated_member_use_from_same_package
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'structured_intent.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// dart format off
T _$identity<T>(T value) => value;

/// @nodoc
mixin _$IntentPreferences {

 String? get language;@JsonKey(name: 'hotel_tier') String? get hotelTier; List<String>? get accessibility;
/// Create a copy of IntentPreferences
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$IntentPreferencesCopyWith<IntentPreferences> get copyWith => _$IntentPreferencesCopyWithImpl<IntentPreferences>(this as IntentPreferences, _$identity);

  /// Serializes this IntentPreferences to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is IntentPreferences&&(identical(other.language, language) || other.language == language)&&(identical(other.hotelTier, hotelTier) || other.hotelTier == hotelTier)&&const DeepCollectionEquality().equals(other.accessibility, accessibility));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,language,hotelTier,const DeepCollectionEquality().hash(accessibility));

@override
String toString() {
  return 'IntentPreferences(language: $language, hotelTier: $hotelTier, accessibility: $accessibility)';
}


}

/// @nodoc
abstract mixin class $IntentPreferencesCopyWith<$Res>  {
  factory $IntentPreferencesCopyWith(IntentPreferences value, $Res Function(IntentPreferences) _then) = _$IntentPreferencesCopyWithImpl;
@useResult
$Res call({
 String? language,@JsonKey(name: 'hotel_tier') String? hotelTier, List<String>? accessibility
});




}
/// @nodoc
class _$IntentPreferencesCopyWithImpl<$Res>
    implements $IntentPreferencesCopyWith<$Res> {
  _$IntentPreferencesCopyWithImpl(this._self, this._then);

  final IntentPreferences _self;
  final $Res Function(IntentPreferences) _then;

/// Create a copy of IntentPreferences
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? language = freezed,Object? hotelTier = freezed,Object? accessibility = freezed,}) {
  return _then(IntentPreferences(
language: freezed == language ? _self.language : language // ignore: cast_nullable_to_non_nullable
as String?,hotelTier: freezed == hotelTier ? _self.hotelTier : hotelTier // ignore: cast_nullable_to_non_nullable
as String?,accessibility: freezed == accessibility ? _self.accessibility : accessibility // ignore: cast_nullable_to_non_nullable
as List<String>?,
  ));
}

}


/// Adds pattern-matching-related methods to [IntentPreferences].
extension IntentPreferencesPatterns on IntentPreferences {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _IntentPreferences value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _IntentPreferences() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _IntentPreferences value)  $default,){
final _that = this;
switch (_that) {
case _IntentPreferences():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _IntentPreferences value)?  $default,){
final _that = this;
switch (_that) {
case _IntentPreferences() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String? language, @JsonKey(name: 'hotel_tier')  String? hotelTier,  List<String>? accessibility)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _IntentPreferences() when $default != null:
return $default(_that.language,_that.hotelTier,_that.accessibility);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String? language, @JsonKey(name: 'hotel_tier')  String? hotelTier,  List<String>? accessibility)  $default,) {final _that = this;
switch (_that) {
case _IntentPreferences():
return $default(_that.language,_that.hotelTier,_that.accessibility);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String? language, @JsonKey(name: 'hotel_tier')  String? hotelTier,  List<String>? accessibility)?  $default,) {final _that = this;
switch (_that) {
case _IntentPreferences() when $default != null:
return $default(_that.language,_that.hotelTier,_that.accessibility);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _IntentPreferences implements IntentPreferences {
  const _IntentPreferences({this.language, @JsonKey(name: 'hotel_tier') this.hotelTier,  List<String>? accessibility}): _accessibility = accessibility;
  factory _IntentPreferences.fromJson(Map<String, dynamic> json) => _$IntentPreferencesFromJson(json);

@override final  String? language;
@override@JsonKey(name: 'hotel_tier') final  String? hotelTier;
 final  List<String>? _accessibility;
@override List<String>? get accessibility {
  final value = _accessibility;
  if (value == null) return null;
  if (_accessibility is EqualUnmodifiableListView) return _accessibility;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(value);
}


/// Create a copy of IntentPreferences
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$IntentPreferencesCopyWith<_IntentPreferences> get copyWith => __$IntentPreferencesCopyWithImpl<_IntentPreferences>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$IntentPreferencesToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _IntentPreferences&&(identical(other.language, language) || other.language == language)&&(identical(other.hotelTier, hotelTier) || other.hotelTier == hotelTier)&&const DeepCollectionEquality().equals(other._accessibility, _accessibility));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,language,hotelTier,const DeepCollectionEquality().hash(_accessibility));

@override
String toString() {
  return 'IntentPreferences(language: $language, hotelTier: $hotelTier, accessibility: $accessibility)';
}


}

/// @nodoc
abstract mixin class _$IntentPreferencesCopyWith<$Res> implements $IntentPreferencesCopyWith<$Res> {
  factory _$IntentPreferencesCopyWith(_IntentPreferences value, $Res Function(_IntentPreferences) _then) = __$IntentPreferencesCopyWithImpl;
@override @useResult
$Res call({
 String? language,@JsonKey(name: 'hotel_tier') String? hotelTier, List<String>? accessibility
});




}
/// @nodoc
class __$IntentPreferencesCopyWithImpl<$Res>
    implements _$IntentPreferencesCopyWith<$Res> {
  __$IntentPreferencesCopyWithImpl(this._self, this._then);

  final _IntentPreferences _self;
  final $Res Function(_IntentPreferences) _then;

/// Create a copy of IntentPreferences
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? language = freezed,Object? hotelTier = freezed,Object? accessibility = freezed,}) {
  return _then(_IntentPreferences(
language: freezed == language ? _self.language : language // ignore: cast_nullable_to_non_nullable
as String?,hotelTier: freezed == hotelTier ? _self.hotelTier : hotelTier // ignore: cast_nullable_to_non_nullable
as String?,accessibility: freezed == accessibility ? _self._accessibility : accessibility // ignore: cast_nullable_to_non_nullable
as List<String>?,
  ));
}


}


/// @nodoc
mixin _$StructuredIntent {

@JsonKey(name: 'schema_version') String get schemaVersion; IntentResolution get resolution;@JsonKey(name: 'intent_category') String? get intentCategory;@JsonKey(name: 'requested_service_text') String get requestedServiceText;@JsonKey(name: 'service_code') String? get serviceCode;@JsonKey(name: 'candidate_service_codes') List<String> get candidateServiceCodes;@JsonKey(name: 'origin_port') String? get originPort;@JsonKey(name: 'date_window') DateWindow? get dateWindow;@JsonKey(name: 'patient_count') int? get patientCount;@JsonKey(name: 'companion_count') int? get companionCount;@JsonKey(name: 'stay_type') StayType? get stayType; Money? get budget; IntentPreferences get preferences;@JsonKey(name: 'missing_fields') List<String> get missingFields;@JsonKey(name: 'clarification_question') String? get clarificationQuestion;@JsonKey(name: 'out_of_scope_reason') String? get outOfScopeReason;@JsonKey(name: 'unsupported_reason') String? get unsupportedReason;
/// Create a copy of StructuredIntent
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$StructuredIntentCopyWith<StructuredIntent> get copyWith => _$StructuredIntentCopyWithImpl<StructuredIntent>(this as StructuredIntent, _$identity);

  /// Serializes this StructuredIntent to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is StructuredIntent&&(identical(other.schemaVersion, schemaVersion) || other.schemaVersion == schemaVersion)&&(identical(other.resolution, resolution) || other.resolution == resolution)&&(identical(other.intentCategory, intentCategory) || other.intentCategory == intentCategory)&&(identical(other.requestedServiceText, requestedServiceText) || other.requestedServiceText == requestedServiceText)&&(identical(other.serviceCode, serviceCode) || other.serviceCode == serviceCode)&&const DeepCollectionEquality().equals(other.candidateServiceCodes, candidateServiceCodes)&&(identical(other.originPort, originPort) || other.originPort == originPort)&&(identical(other.dateWindow, dateWindow) || other.dateWindow == dateWindow)&&(identical(other.patientCount, patientCount) || other.patientCount == patientCount)&&(identical(other.companionCount, companionCount) || other.companionCount == companionCount)&&(identical(other.stayType, stayType) || other.stayType == stayType)&&(identical(other.budget, budget) || other.budget == budget)&&(identical(other.preferences, preferences) || other.preferences == preferences)&&const DeepCollectionEquality().equals(other.missingFields, missingFields)&&(identical(other.clarificationQuestion, clarificationQuestion) || other.clarificationQuestion == clarificationQuestion)&&(identical(other.outOfScopeReason, outOfScopeReason) || other.outOfScopeReason == outOfScopeReason)&&(identical(other.unsupportedReason, unsupportedReason) || other.unsupportedReason == unsupportedReason));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,schemaVersion,resolution,intentCategory,requestedServiceText,serviceCode,const DeepCollectionEquality().hash(candidateServiceCodes),originPort,dateWindow,patientCount,companionCount,stayType,budget,preferences,const DeepCollectionEquality().hash(missingFields),clarificationQuestion,outOfScopeReason,unsupportedReason);

@override
String toString() {
  return 'StructuredIntent(schemaVersion: $schemaVersion, resolution: $resolution, intentCategory: $intentCategory, requestedServiceText: $requestedServiceText, serviceCode: $serviceCode, candidateServiceCodes: $candidateServiceCodes, originPort: $originPort, dateWindow: $dateWindow, patientCount: $patientCount, companionCount: $companionCount, stayType: $stayType, budget: $budget, preferences: $preferences, missingFields: $missingFields, clarificationQuestion: $clarificationQuestion, outOfScopeReason: $outOfScopeReason, unsupportedReason: $unsupportedReason)';
}


}

/// @nodoc
abstract mixin class $StructuredIntentCopyWith<$Res>  {
  factory $StructuredIntentCopyWith(StructuredIntent value, $Res Function(StructuredIntent) _then) = _$StructuredIntentCopyWithImpl;
@useResult
$Res call({
@JsonKey(name: 'schema_version') String schemaVersion, IntentResolution resolution,@JsonKey(name: 'intent_category') String? intentCategory,@JsonKey(name: 'requested_service_text') String requestedServiceText,@JsonKey(name: 'service_code') String? serviceCode,@JsonKey(name: 'candidate_service_codes') List<String> candidateServiceCodes,@JsonKey(name: 'origin_port') String? originPort,@JsonKey(name: 'date_window') DateWindow? dateWindow,@JsonKey(name: 'patient_count') int? patientCount,@JsonKey(name: 'companion_count') int? companionCount,@JsonKey(name: 'stay_type') StayType? stayType, Money? budget, IntentPreferences preferences,@JsonKey(name: 'missing_fields') List<String> missingFields,@JsonKey(name: 'clarification_question') String? clarificationQuestion,@JsonKey(name: 'out_of_scope_reason') String? outOfScopeReason,@JsonKey(name: 'unsupported_reason') String? unsupportedReason
});


$DateWindowCopyWith<$Res>? get dateWindow;$MoneyCopyWith<$Res>? get budget;$IntentPreferencesCopyWith<$Res> get preferences;

}
/// @nodoc
class _$StructuredIntentCopyWithImpl<$Res>
    implements $StructuredIntentCopyWith<$Res> {
  _$StructuredIntentCopyWithImpl(this._self, this._then);

  final StructuredIntent _self;
  final $Res Function(StructuredIntent) _then;

/// Create a copy of StructuredIntent
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? schemaVersion = null,Object? resolution = null,Object? intentCategory = freezed,Object? requestedServiceText = null,Object? serviceCode = freezed,Object? candidateServiceCodes = null,Object? originPort = freezed,Object? dateWindow = freezed,Object? patientCount = freezed,Object? companionCount = freezed,Object? stayType = freezed,Object? budget = freezed,Object? preferences = null,Object? missingFields = null,Object? clarificationQuestion = freezed,Object? outOfScopeReason = freezed,Object? unsupportedReason = freezed,}) {
  return _then(StructuredIntent(
schemaVersion: null == schemaVersion ? _self.schemaVersion : schemaVersion // ignore: cast_nullable_to_non_nullable
as String,resolution: null == resolution ? _self.resolution : resolution // ignore: cast_nullable_to_non_nullable
as IntentResolution,intentCategory: freezed == intentCategory ? _self.intentCategory : intentCategory // ignore: cast_nullable_to_non_nullable
as String?,requestedServiceText: null == requestedServiceText ? _self.requestedServiceText : requestedServiceText // ignore: cast_nullable_to_non_nullable
as String,serviceCode: freezed == serviceCode ? _self.serviceCode : serviceCode // ignore: cast_nullable_to_non_nullable
as String?,candidateServiceCodes: null == candidateServiceCodes ? _self.candidateServiceCodes : candidateServiceCodes // ignore: cast_nullable_to_non_nullable
as List<String>,originPort: freezed == originPort ? _self.originPort : originPort // ignore: cast_nullable_to_non_nullable
as String?,dateWindow: freezed == dateWindow ? _self.dateWindow : dateWindow // ignore: cast_nullable_to_non_nullable
as DateWindow?,patientCount: freezed == patientCount ? _self.patientCount : patientCount // ignore: cast_nullable_to_non_nullable
as int?,companionCount: freezed == companionCount ? _self.companionCount : companionCount // ignore: cast_nullable_to_non_nullable
as int?,stayType: freezed == stayType ? _self.stayType : stayType // ignore: cast_nullable_to_non_nullable
as StayType?,budget: freezed == budget ? _self.budget : budget // ignore: cast_nullable_to_non_nullable
as Money?,preferences: null == preferences ? _self.preferences : preferences // ignore: cast_nullable_to_non_nullable
as IntentPreferences,missingFields: null == missingFields ? _self.missingFields : missingFields // ignore: cast_nullable_to_non_nullable
as List<String>,clarificationQuestion: freezed == clarificationQuestion ? _self.clarificationQuestion : clarificationQuestion // ignore: cast_nullable_to_non_nullable
as String?,outOfScopeReason: freezed == outOfScopeReason ? _self.outOfScopeReason : outOfScopeReason // ignore: cast_nullable_to_non_nullable
as String?,unsupportedReason: freezed == unsupportedReason ? _self.unsupportedReason : unsupportedReason // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}
/// Create a copy of StructuredIntent
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$DateWindowCopyWith<$Res>? get dateWindow {
    if (_self.dateWindow == null) {
    return null;
  }

  return $DateWindowCopyWith<$Res>(_self.dateWindow!, (value) {
    return _then(_self.copyWith(dateWindow: value));
  });
}/// Create a copy of StructuredIntent
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$MoneyCopyWith<$Res>? get budget {
    if (_self.budget == null) {
    return null;
  }

  return $MoneyCopyWith<$Res>(_self.budget!, (value) {
    return _then(_self.copyWith(budget: value));
  });
}/// Create a copy of StructuredIntent
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$IntentPreferencesCopyWith<$Res> get preferences {
  
  return $IntentPreferencesCopyWith<$Res>(_self.preferences, (value) {
    return _then(_self.copyWith(preferences: value));
  });
}
}


/// Adds pattern-matching-related methods to [StructuredIntent].
extension StructuredIntentPatterns on StructuredIntent {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _StructuredIntent value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _StructuredIntent() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _StructuredIntent value)  $default,){
final _that = this;
switch (_that) {
case _StructuredIntent():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _StructuredIntent value)?  $default,){
final _that = this;
switch (_that) {
case _StructuredIntent() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function(@JsonKey(name: 'schema_version')  String schemaVersion,  IntentResolution resolution, @JsonKey(name: 'intent_category')  String? intentCategory, @JsonKey(name: 'requested_service_text')  String requestedServiceText, @JsonKey(name: 'service_code')  String? serviceCode, @JsonKey(name: 'candidate_service_codes')  List<String> candidateServiceCodes, @JsonKey(name: 'origin_port')  String? originPort, @JsonKey(name: 'date_window')  DateWindow? dateWindow, @JsonKey(name: 'patient_count')  int? patientCount, @JsonKey(name: 'companion_count')  int? companionCount, @JsonKey(name: 'stay_type')  StayType? stayType,  Money? budget,  IntentPreferences preferences, @JsonKey(name: 'missing_fields')  List<String> missingFields, @JsonKey(name: 'clarification_question')  String? clarificationQuestion, @JsonKey(name: 'out_of_scope_reason')  String? outOfScopeReason, @JsonKey(name: 'unsupported_reason')  String? unsupportedReason)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _StructuredIntent() when $default != null:
return $default(_that.schemaVersion,_that.resolution,_that.intentCategory,_that.requestedServiceText,_that.serviceCode,_that.candidateServiceCodes,_that.originPort,_that.dateWindow,_that.patientCount,_that.companionCount,_that.stayType,_that.budget,_that.preferences,_that.missingFields,_that.clarificationQuestion,_that.outOfScopeReason,_that.unsupportedReason);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function(@JsonKey(name: 'schema_version')  String schemaVersion,  IntentResolution resolution, @JsonKey(name: 'intent_category')  String? intentCategory, @JsonKey(name: 'requested_service_text')  String requestedServiceText, @JsonKey(name: 'service_code')  String? serviceCode, @JsonKey(name: 'candidate_service_codes')  List<String> candidateServiceCodes, @JsonKey(name: 'origin_port')  String? originPort, @JsonKey(name: 'date_window')  DateWindow? dateWindow, @JsonKey(name: 'patient_count')  int? patientCount, @JsonKey(name: 'companion_count')  int? companionCount, @JsonKey(name: 'stay_type')  StayType? stayType,  Money? budget,  IntentPreferences preferences, @JsonKey(name: 'missing_fields')  List<String> missingFields, @JsonKey(name: 'clarification_question')  String? clarificationQuestion, @JsonKey(name: 'out_of_scope_reason')  String? outOfScopeReason, @JsonKey(name: 'unsupported_reason')  String? unsupportedReason)  $default,) {final _that = this;
switch (_that) {
case _StructuredIntent():
return $default(_that.schemaVersion,_that.resolution,_that.intentCategory,_that.requestedServiceText,_that.serviceCode,_that.candidateServiceCodes,_that.originPort,_that.dateWindow,_that.patientCount,_that.companionCount,_that.stayType,_that.budget,_that.preferences,_that.missingFields,_that.clarificationQuestion,_that.outOfScopeReason,_that.unsupportedReason);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function(@JsonKey(name: 'schema_version')  String schemaVersion,  IntentResolution resolution, @JsonKey(name: 'intent_category')  String? intentCategory, @JsonKey(name: 'requested_service_text')  String requestedServiceText, @JsonKey(name: 'service_code')  String? serviceCode, @JsonKey(name: 'candidate_service_codes')  List<String> candidateServiceCodes, @JsonKey(name: 'origin_port')  String? originPort, @JsonKey(name: 'date_window')  DateWindow? dateWindow, @JsonKey(name: 'patient_count')  int? patientCount, @JsonKey(name: 'companion_count')  int? companionCount, @JsonKey(name: 'stay_type')  StayType? stayType,  Money? budget,  IntentPreferences preferences, @JsonKey(name: 'missing_fields')  List<String> missingFields, @JsonKey(name: 'clarification_question')  String? clarificationQuestion, @JsonKey(name: 'out_of_scope_reason')  String? outOfScopeReason, @JsonKey(name: 'unsupported_reason')  String? unsupportedReason)?  $default,) {final _that = this;
switch (_that) {
case _StructuredIntent() when $default != null:
return $default(_that.schemaVersion,_that.resolution,_that.intentCategory,_that.requestedServiceText,_that.serviceCode,_that.candidateServiceCodes,_that.originPort,_that.dateWindow,_that.patientCount,_that.companionCount,_that.stayType,_that.budget,_that.preferences,_that.missingFields,_that.clarificationQuestion,_that.outOfScopeReason,_that.unsupportedReason);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _StructuredIntent implements StructuredIntent {
  const _StructuredIntent({@JsonKey(name: 'schema_version') required this.schemaVersion, required this.resolution, @JsonKey(name: 'intent_category') this.intentCategory, @JsonKey(name: 'requested_service_text') required this.requestedServiceText, @JsonKey(name: 'service_code') this.serviceCode, @JsonKey(name: 'candidate_service_codes') required  List<String> candidateServiceCodes, @JsonKey(name: 'origin_port') this.originPort, @JsonKey(name: 'date_window') this.dateWindow, @JsonKey(name: 'patient_count') this.patientCount, @JsonKey(name: 'companion_count') this.companionCount, @JsonKey(name: 'stay_type') this.stayType, this.budget, required this.preferences, @JsonKey(name: 'missing_fields') required  List<String> missingFields, @JsonKey(name: 'clarification_question') this.clarificationQuestion, @JsonKey(name: 'out_of_scope_reason') this.outOfScopeReason, @JsonKey(name: 'unsupported_reason') this.unsupportedReason}): _candidateServiceCodes = candidateServiceCodes,_missingFields = missingFields;
  factory _StructuredIntent.fromJson(Map<String, dynamic> json) => _$StructuredIntentFromJson(json);

@override@JsonKey(name: 'schema_version') final  String schemaVersion;
@override final  IntentResolution resolution;
@override@JsonKey(name: 'intent_category') final  String? intentCategory;
@override@JsonKey(name: 'requested_service_text') final  String requestedServiceText;
@override@JsonKey(name: 'service_code') final  String? serviceCode;
 final  List<String> _candidateServiceCodes;
@override@JsonKey(name: 'candidate_service_codes') List<String> get candidateServiceCodes {
  if (_candidateServiceCodes is EqualUnmodifiableListView) return _candidateServiceCodes;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_candidateServiceCodes);
}

@override@JsonKey(name: 'origin_port') final  String? originPort;
@override@JsonKey(name: 'date_window') final  DateWindow? dateWindow;
@override@JsonKey(name: 'patient_count') final  int? patientCount;
@override@JsonKey(name: 'companion_count') final  int? companionCount;
@override@JsonKey(name: 'stay_type') final  StayType? stayType;
@override final  Money? budget;
@override final  IntentPreferences preferences;
 final  List<String> _missingFields;
@override@JsonKey(name: 'missing_fields') List<String> get missingFields {
  if (_missingFields is EqualUnmodifiableListView) return _missingFields;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_missingFields);
}

@override@JsonKey(name: 'clarification_question') final  String? clarificationQuestion;
@override@JsonKey(name: 'out_of_scope_reason') final  String? outOfScopeReason;
@override@JsonKey(name: 'unsupported_reason') final  String? unsupportedReason;

/// Create a copy of StructuredIntent
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$StructuredIntentCopyWith<_StructuredIntent> get copyWith => __$StructuredIntentCopyWithImpl<_StructuredIntent>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$StructuredIntentToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _StructuredIntent&&(identical(other.schemaVersion, schemaVersion) || other.schemaVersion == schemaVersion)&&(identical(other.resolution, resolution) || other.resolution == resolution)&&(identical(other.intentCategory, intentCategory) || other.intentCategory == intentCategory)&&(identical(other.requestedServiceText, requestedServiceText) || other.requestedServiceText == requestedServiceText)&&(identical(other.serviceCode, serviceCode) || other.serviceCode == serviceCode)&&const DeepCollectionEquality().equals(other._candidateServiceCodes, _candidateServiceCodes)&&(identical(other.originPort, originPort) || other.originPort == originPort)&&(identical(other.dateWindow, dateWindow) || other.dateWindow == dateWindow)&&(identical(other.patientCount, patientCount) || other.patientCount == patientCount)&&(identical(other.companionCount, companionCount) || other.companionCount == companionCount)&&(identical(other.stayType, stayType) || other.stayType == stayType)&&(identical(other.budget, budget) || other.budget == budget)&&(identical(other.preferences, preferences) || other.preferences == preferences)&&const DeepCollectionEquality().equals(other._missingFields, _missingFields)&&(identical(other.clarificationQuestion, clarificationQuestion) || other.clarificationQuestion == clarificationQuestion)&&(identical(other.outOfScopeReason, outOfScopeReason) || other.outOfScopeReason == outOfScopeReason)&&(identical(other.unsupportedReason, unsupportedReason) || other.unsupportedReason == unsupportedReason));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,schemaVersion,resolution,intentCategory,requestedServiceText,serviceCode,const DeepCollectionEquality().hash(_candidateServiceCodes),originPort,dateWindow,patientCount,companionCount,stayType,budget,preferences,const DeepCollectionEquality().hash(_missingFields),clarificationQuestion,outOfScopeReason,unsupportedReason);

@override
String toString() {
  return 'StructuredIntent(schemaVersion: $schemaVersion, resolution: $resolution, intentCategory: $intentCategory, requestedServiceText: $requestedServiceText, serviceCode: $serviceCode, candidateServiceCodes: $candidateServiceCodes, originPort: $originPort, dateWindow: $dateWindow, patientCount: $patientCount, companionCount: $companionCount, stayType: $stayType, budget: $budget, preferences: $preferences, missingFields: $missingFields, clarificationQuestion: $clarificationQuestion, outOfScopeReason: $outOfScopeReason, unsupportedReason: $unsupportedReason)';
}


}

/// @nodoc
abstract mixin class _$StructuredIntentCopyWith<$Res> implements $StructuredIntentCopyWith<$Res> {
  factory _$StructuredIntentCopyWith(_StructuredIntent value, $Res Function(_StructuredIntent) _then) = __$StructuredIntentCopyWithImpl;
@override @useResult
$Res call({
@JsonKey(name: 'schema_version') String schemaVersion, IntentResolution resolution,@JsonKey(name: 'intent_category') String? intentCategory,@JsonKey(name: 'requested_service_text') String requestedServiceText,@JsonKey(name: 'service_code') String? serviceCode,@JsonKey(name: 'candidate_service_codes') List<String> candidateServiceCodes,@JsonKey(name: 'origin_port') String? originPort,@JsonKey(name: 'date_window') DateWindow? dateWindow,@JsonKey(name: 'patient_count') int? patientCount,@JsonKey(name: 'companion_count') int? companionCount,@JsonKey(name: 'stay_type') StayType? stayType, Money? budget, IntentPreferences preferences,@JsonKey(name: 'missing_fields') List<String> missingFields,@JsonKey(name: 'clarification_question') String? clarificationQuestion,@JsonKey(name: 'out_of_scope_reason') String? outOfScopeReason,@JsonKey(name: 'unsupported_reason') String? unsupportedReason
});


@override $DateWindowCopyWith<$Res>? get dateWindow;@override $MoneyCopyWith<$Res>? get budget;@override $IntentPreferencesCopyWith<$Res> get preferences;

}
/// @nodoc
class __$StructuredIntentCopyWithImpl<$Res>
    implements _$StructuredIntentCopyWith<$Res> {
  __$StructuredIntentCopyWithImpl(this._self, this._then);

  final _StructuredIntent _self;
  final $Res Function(_StructuredIntent) _then;

/// Create a copy of StructuredIntent
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? schemaVersion = null,Object? resolution = null,Object? intentCategory = freezed,Object? requestedServiceText = null,Object? serviceCode = freezed,Object? candidateServiceCodes = null,Object? originPort = freezed,Object? dateWindow = freezed,Object? patientCount = freezed,Object? companionCount = freezed,Object? stayType = freezed,Object? budget = freezed,Object? preferences = null,Object? missingFields = null,Object? clarificationQuestion = freezed,Object? outOfScopeReason = freezed,Object? unsupportedReason = freezed,}) {
  return _then(_StructuredIntent(
schemaVersion: null == schemaVersion ? _self.schemaVersion : schemaVersion // ignore: cast_nullable_to_non_nullable
as String,resolution: null == resolution ? _self.resolution : resolution // ignore: cast_nullable_to_non_nullable
as IntentResolution,intentCategory: freezed == intentCategory ? _self.intentCategory : intentCategory // ignore: cast_nullable_to_non_nullable
as String?,requestedServiceText: null == requestedServiceText ? _self.requestedServiceText : requestedServiceText // ignore: cast_nullable_to_non_nullable
as String,serviceCode: freezed == serviceCode ? _self.serviceCode : serviceCode // ignore: cast_nullable_to_non_nullable
as String?,candidateServiceCodes: null == candidateServiceCodes ? _self._candidateServiceCodes : candidateServiceCodes // ignore: cast_nullable_to_non_nullable
as List<String>,originPort: freezed == originPort ? _self.originPort : originPort // ignore: cast_nullable_to_non_nullable
as String?,dateWindow: freezed == dateWindow ? _self.dateWindow : dateWindow // ignore: cast_nullable_to_non_nullable
as DateWindow?,patientCount: freezed == patientCount ? _self.patientCount : patientCount // ignore: cast_nullable_to_non_nullable
as int?,companionCount: freezed == companionCount ? _self.companionCount : companionCount // ignore: cast_nullable_to_non_nullable
as int?,stayType: freezed == stayType ? _self.stayType : stayType // ignore: cast_nullable_to_non_nullable
as StayType?,budget: freezed == budget ? _self.budget : budget // ignore: cast_nullable_to_non_nullable
as Money?,preferences: null == preferences ? _self.preferences : preferences // ignore: cast_nullable_to_non_nullable
as IntentPreferences,missingFields: null == missingFields ? _self._missingFields : missingFields // ignore: cast_nullable_to_non_nullable
as List<String>,clarificationQuestion: freezed == clarificationQuestion ? _self.clarificationQuestion : clarificationQuestion // ignore: cast_nullable_to_non_nullable
as String?,outOfScopeReason: freezed == outOfScopeReason ? _self.outOfScopeReason : outOfScopeReason // ignore: cast_nullable_to_non_nullable
as String?,unsupportedReason: freezed == unsupportedReason ? _self.unsupportedReason : unsupportedReason // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}

/// Create a copy of StructuredIntent
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$DateWindowCopyWith<$Res>? get dateWindow {
    if (_self.dateWindow == null) {
    return null;
  }

  return $DateWindowCopyWith<$Res>(_self.dateWindow!, (value) {
    return _then(_self.copyWith(dateWindow: value));
  });
}/// Create a copy of StructuredIntent
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$MoneyCopyWith<$Res>? get budget {
    if (_self.budget == null) {
    return null;
  }

  return $MoneyCopyWith<$Res>(_self.budget!, (value) {
    return _then(_self.copyWith(budget: value));
  });
}/// Create a copy of StructuredIntent
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$IntentPreferencesCopyWith<$Res> get preferences {
  
  return $IntentPreferencesCopyWith<$Res>(_self.preferences, (value) {
    return _then(_self.copyWith(preferences: value));
  });
}
}


/// @nodoc
mixin _$IntentCorrections {

@JsonKey(name: 'service_code') String? get serviceCode;@JsonKey(name: 'origin_port') String? get originPort;@JsonKey(name: 'date_window') DateWindow? get dateWindow;@JsonKey(name: 'patient_count') int? get patientCount;@JsonKey(name: 'companion_count') int? get companionCount;@JsonKey(name: 'stay_type') StayType? get stayType; Money? get budget; IntentPreferences? get preferences;
/// Create a copy of IntentCorrections
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$IntentCorrectionsCopyWith<IntentCorrections> get copyWith => _$IntentCorrectionsCopyWithImpl<IntentCorrections>(this as IntentCorrections, _$identity);

  /// Serializes this IntentCorrections to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is IntentCorrections&&(identical(other.serviceCode, serviceCode) || other.serviceCode == serviceCode)&&(identical(other.originPort, originPort) || other.originPort == originPort)&&(identical(other.dateWindow, dateWindow) || other.dateWindow == dateWindow)&&(identical(other.patientCount, patientCount) || other.patientCount == patientCount)&&(identical(other.companionCount, companionCount) || other.companionCount == companionCount)&&(identical(other.stayType, stayType) || other.stayType == stayType)&&(identical(other.budget, budget) || other.budget == budget)&&(identical(other.preferences, preferences) || other.preferences == preferences));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,serviceCode,originPort,dateWindow,patientCount,companionCount,stayType,budget,preferences);

@override
String toString() {
  return 'IntentCorrections(serviceCode: $serviceCode, originPort: $originPort, dateWindow: $dateWindow, patientCount: $patientCount, companionCount: $companionCount, stayType: $stayType, budget: $budget, preferences: $preferences)';
}


}

/// @nodoc
abstract mixin class $IntentCorrectionsCopyWith<$Res>  {
  factory $IntentCorrectionsCopyWith(IntentCorrections value, $Res Function(IntentCorrections) _then) = _$IntentCorrectionsCopyWithImpl;
@useResult
$Res call({
@JsonKey(name: 'service_code') String? serviceCode,@JsonKey(name: 'origin_port') String? originPort,@JsonKey(name: 'date_window') DateWindow? dateWindow,@JsonKey(name: 'patient_count') int? patientCount,@JsonKey(name: 'companion_count') int? companionCount,@JsonKey(name: 'stay_type') StayType? stayType, Money? budget, IntentPreferences? preferences
});


$DateWindowCopyWith<$Res>? get dateWindow;$MoneyCopyWith<$Res>? get budget;$IntentPreferencesCopyWith<$Res>? get preferences;

}
/// @nodoc
class _$IntentCorrectionsCopyWithImpl<$Res>
    implements $IntentCorrectionsCopyWith<$Res> {
  _$IntentCorrectionsCopyWithImpl(this._self, this._then);

  final IntentCorrections _self;
  final $Res Function(IntentCorrections) _then;

/// Create a copy of IntentCorrections
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? serviceCode = freezed,Object? originPort = freezed,Object? dateWindow = freezed,Object? patientCount = freezed,Object? companionCount = freezed,Object? stayType = freezed,Object? budget = freezed,Object? preferences = freezed,}) {
  return _then(IntentCorrections(
serviceCode: freezed == serviceCode ? _self.serviceCode : serviceCode // ignore: cast_nullable_to_non_nullable
as String?,originPort: freezed == originPort ? _self.originPort : originPort // ignore: cast_nullable_to_non_nullable
as String?,dateWindow: freezed == dateWindow ? _self.dateWindow : dateWindow // ignore: cast_nullable_to_non_nullable
as DateWindow?,patientCount: freezed == patientCount ? _self.patientCount : patientCount // ignore: cast_nullable_to_non_nullable
as int?,companionCount: freezed == companionCount ? _self.companionCount : companionCount // ignore: cast_nullable_to_non_nullable
as int?,stayType: freezed == stayType ? _self.stayType : stayType // ignore: cast_nullable_to_non_nullable
as StayType?,budget: freezed == budget ? _self.budget : budget // ignore: cast_nullable_to_non_nullable
as Money?,preferences: freezed == preferences ? _self.preferences : preferences // ignore: cast_nullable_to_non_nullable
as IntentPreferences?,
  ));
}
/// Create a copy of IntentCorrections
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$DateWindowCopyWith<$Res>? get dateWindow {
    if (_self.dateWindow == null) {
    return null;
  }

  return $DateWindowCopyWith<$Res>(_self.dateWindow!, (value) {
    return _then(_self.copyWith(dateWindow: value));
  });
}/// Create a copy of IntentCorrections
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$MoneyCopyWith<$Res>? get budget {
    if (_self.budget == null) {
    return null;
  }

  return $MoneyCopyWith<$Res>(_self.budget!, (value) {
    return _then(_self.copyWith(budget: value));
  });
}/// Create a copy of IntentCorrections
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$IntentPreferencesCopyWith<$Res>? get preferences {
    if (_self.preferences == null) {
    return null;
  }

  return $IntentPreferencesCopyWith<$Res>(_self.preferences!, (value) {
    return _then(_self.copyWith(preferences: value));
  });
}
}


/// Adds pattern-matching-related methods to [IntentCorrections].
extension IntentCorrectionsPatterns on IntentCorrections {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _IntentCorrections value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _IntentCorrections() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _IntentCorrections value)  $default,){
final _that = this;
switch (_that) {
case _IntentCorrections():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _IntentCorrections value)?  $default,){
final _that = this;
switch (_that) {
case _IntentCorrections() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function(@JsonKey(name: 'service_code')  String? serviceCode, @JsonKey(name: 'origin_port')  String? originPort, @JsonKey(name: 'date_window')  DateWindow? dateWindow, @JsonKey(name: 'patient_count')  int? patientCount, @JsonKey(name: 'companion_count')  int? companionCount, @JsonKey(name: 'stay_type')  StayType? stayType,  Money? budget,  IntentPreferences? preferences)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _IntentCorrections() when $default != null:
return $default(_that.serviceCode,_that.originPort,_that.dateWindow,_that.patientCount,_that.companionCount,_that.stayType,_that.budget,_that.preferences);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function(@JsonKey(name: 'service_code')  String? serviceCode, @JsonKey(name: 'origin_port')  String? originPort, @JsonKey(name: 'date_window')  DateWindow? dateWindow, @JsonKey(name: 'patient_count')  int? patientCount, @JsonKey(name: 'companion_count')  int? companionCount, @JsonKey(name: 'stay_type')  StayType? stayType,  Money? budget,  IntentPreferences? preferences)  $default,) {final _that = this;
switch (_that) {
case _IntentCorrections():
return $default(_that.serviceCode,_that.originPort,_that.dateWindow,_that.patientCount,_that.companionCount,_that.stayType,_that.budget,_that.preferences);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function(@JsonKey(name: 'service_code')  String? serviceCode, @JsonKey(name: 'origin_port')  String? originPort, @JsonKey(name: 'date_window')  DateWindow? dateWindow, @JsonKey(name: 'patient_count')  int? patientCount, @JsonKey(name: 'companion_count')  int? companionCount, @JsonKey(name: 'stay_type')  StayType? stayType,  Money? budget,  IntentPreferences? preferences)?  $default,) {final _that = this;
switch (_that) {
case _IntentCorrections() when $default != null:
return $default(_that.serviceCode,_that.originPort,_that.dateWindow,_that.patientCount,_that.companionCount,_that.stayType,_that.budget,_that.preferences);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _IntentCorrections implements IntentCorrections {
  const _IntentCorrections({@JsonKey(name: 'service_code') this.serviceCode, @JsonKey(name: 'origin_port') this.originPort, @JsonKey(name: 'date_window') this.dateWindow, @JsonKey(name: 'patient_count') this.patientCount, @JsonKey(name: 'companion_count') this.companionCount, @JsonKey(name: 'stay_type') this.stayType, this.budget, this.preferences});
  factory _IntentCorrections.fromJson(Map<String, dynamic> json) => _$IntentCorrectionsFromJson(json);

@override@JsonKey(name: 'service_code') final  String? serviceCode;
@override@JsonKey(name: 'origin_port') final  String? originPort;
@override@JsonKey(name: 'date_window') final  DateWindow? dateWindow;
@override@JsonKey(name: 'patient_count') final  int? patientCount;
@override@JsonKey(name: 'companion_count') final  int? companionCount;
@override@JsonKey(name: 'stay_type') final  StayType? stayType;
@override final  Money? budget;
@override final  IntentPreferences? preferences;

/// Create a copy of IntentCorrections
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$IntentCorrectionsCopyWith<_IntentCorrections> get copyWith => __$IntentCorrectionsCopyWithImpl<_IntentCorrections>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$IntentCorrectionsToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _IntentCorrections&&(identical(other.serviceCode, serviceCode) || other.serviceCode == serviceCode)&&(identical(other.originPort, originPort) || other.originPort == originPort)&&(identical(other.dateWindow, dateWindow) || other.dateWindow == dateWindow)&&(identical(other.patientCount, patientCount) || other.patientCount == patientCount)&&(identical(other.companionCount, companionCount) || other.companionCount == companionCount)&&(identical(other.stayType, stayType) || other.stayType == stayType)&&(identical(other.budget, budget) || other.budget == budget)&&(identical(other.preferences, preferences) || other.preferences == preferences));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,serviceCode,originPort,dateWindow,patientCount,companionCount,stayType,budget,preferences);

@override
String toString() {
  return 'IntentCorrections(serviceCode: $serviceCode, originPort: $originPort, dateWindow: $dateWindow, patientCount: $patientCount, companionCount: $companionCount, stayType: $stayType, budget: $budget, preferences: $preferences)';
}


}

/// @nodoc
abstract mixin class _$IntentCorrectionsCopyWith<$Res> implements $IntentCorrectionsCopyWith<$Res> {
  factory _$IntentCorrectionsCopyWith(_IntentCorrections value, $Res Function(_IntentCorrections) _then) = __$IntentCorrectionsCopyWithImpl;
@override @useResult
$Res call({
@JsonKey(name: 'service_code') String? serviceCode,@JsonKey(name: 'origin_port') String? originPort,@JsonKey(name: 'date_window') DateWindow? dateWindow,@JsonKey(name: 'patient_count') int? patientCount,@JsonKey(name: 'companion_count') int? companionCount,@JsonKey(name: 'stay_type') StayType? stayType, Money? budget, IntentPreferences? preferences
});


@override $DateWindowCopyWith<$Res>? get dateWindow;@override $MoneyCopyWith<$Res>? get budget;@override $IntentPreferencesCopyWith<$Res>? get preferences;

}
/// @nodoc
class __$IntentCorrectionsCopyWithImpl<$Res>
    implements _$IntentCorrectionsCopyWith<$Res> {
  __$IntentCorrectionsCopyWithImpl(this._self, this._then);

  final _IntentCorrections _self;
  final $Res Function(_IntentCorrections) _then;

/// Create a copy of IntentCorrections
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? serviceCode = freezed,Object? originPort = freezed,Object? dateWindow = freezed,Object? patientCount = freezed,Object? companionCount = freezed,Object? stayType = freezed,Object? budget = freezed,Object? preferences = freezed,}) {
  return _then(_IntentCorrections(
serviceCode: freezed == serviceCode ? _self.serviceCode : serviceCode // ignore: cast_nullable_to_non_nullable
as String?,originPort: freezed == originPort ? _self.originPort : originPort // ignore: cast_nullable_to_non_nullable
as String?,dateWindow: freezed == dateWindow ? _self.dateWindow : dateWindow // ignore: cast_nullable_to_non_nullable
as DateWindow?,patientCount: freezed == patientCount ? _self.patientCount : patientCount // ignore: cast_nullable_to_non_nullable
as int?,companionCount: freezed == companionCount ? _self.companionCount : companionCount // ignore: cast_nullable_to_non_nullable
as int?,stayType: freezed == stayType ? _self.stayType : stayType // ignore: cast_nullable_to_non_nullable
as StayType?,budget: freezed == budget ? _self.budget : budget // ignore: cast_nullable_to_non_nullable
as Money?,preferences: freezed == preferences ? _self.preferences : preferences // ignore: cast_nullable_to_non_nullable
as IntentPreferences?,
  ));
}

/// Create a copy of IntentCorrections
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$DateWindowCopyWith<$Res>? get dateWindow {
    if (_self.dateWindow == null) {
    return null;
  }

  return $DateWindowCopyWith<$Res>(_self.dateWindow!, (value) {
    return _then(_self.copyWith(dateWindow: value));
  });
}/// Create a copy of IntentCorrections
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$MoneyCopyWith<$Res>? get budget {
    if (_self.budget == null) {
    return null;
  }

  return $MoneyCopyWith<$Res>(_self.budget!, (value) {
    return _then(_self.copyWith(budget: value));
  });
}/// Create a copy of IntentCorrections
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$IntentPreferencesCopyWith<$Res>? get preferences {
    if (_self.preferences == null) {
    return null;
  }

  return $IntentPreferencesCopyWith<$Res>(_self.preferences!, (value) {
    return _then(_self.copyWith(preferences: value));
  });
}
}


/// @nodoc
mixin _$AmendIntentRequest {

 String? get answer; IntentCorrections? get corrections;
/// Create a copy of AmendIntentRequest
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$AmendIntentRequestCopyWith<AmendIntentRequest> get copyWith => _$AmendIntentRequestCopyWithImpl<AmendIntentRequest>(this as AmendIntentRequest, _$identity);

  /// Serializes this AmendIntentRequest to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is AmendIntentRequest&&(identical(other.answer, answer) || other.answer == answer)&&(identical(other.corrections, corrections) || other.corrections == corrections));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,answer,corrections);

@override
String toString() {
  return 'AmendIntentRequest(answer: $answer, corrections: $corrections)';
}


}

/// @nodoc
abstract mixin class $AmendIntentRequestCopyWith<$Res>  {
  factory $AmendIntentRequestCopyWith(AmendIntentRequest value, $Res Function(AmendIntentRequest) _then) = _$AmendIntentRequestCopyWithImpl;
@useResult
$Res call({
 String? answer, IntentCorrections? corrections
});


$IntentCorrectionsCopyWith<$Res>? get corrections;

}
/// @nodoc
class _$AmendIntentRequestCopyWithImpl<$Res>
    implements $AmendIntentRequestCopyWith<$Res> {
  _$AmendIntentRequestCopyWithImpl(this._self, this._then);

  final AmendIntentRequest _self;
  final $Res Function(AmendIntentRequest) _then;

/// Create a copy of AmendIntentRequest
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? answer = freezed,Object? corrections = freezed,}) {
  return _then(AmendIntentRequest(
answer: freezed == answer ? _self.answer : answer // ignore: cast_nullable_to_non_nullable
as String?,corrections: freezed == corrections ? _self.corrections : corrections // ignore: cast_nullable_to_non_nullable
as IntentCorrections?,
  ));
}
/// Create a copy of AmendIntentRequest
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$IntentCorrectionsCopyWith<$Res>? get corrections {
    if (_self.corrections == null) {
    return null;
  }

  return $IntentCorrectionsCopyWith<$Res>(_self.corrections!, (value) {
    return _then(_self.copyWith(corrections: value));
  });
}
}


/// Adds pattern-matching-related methods to [AmendIntentRequest].
extension AmendIntentRequestPatterns on AmendIntentRequest {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _AmendIntentRequest value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _AmendIntentRequest() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _AmendIntentRequest value)  $default,){
final _that = this;
switch (_that) {
case _AmendIntentRequest():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _AmendIntentRequest value)?  $default,){
final _that = this;
switch (_that) {
case _AmendIntentRequest() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String? answer,  IntentCorrections? corrections)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _AmendIntentRequest() when $default != null:
return $default(_that.answer,_that.corrections);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String? answer,  IntentCorrections? corrections)  $default,) {final _that = this;
switch (_that) {
case _AmendIntentRequest():
return $default(_that.answer,_that.corrections);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String? answer,  IntentCorrections? corrections)?  $default,) {final _that = this;
switch (_that) {
case _AmendIntentRequest() when $default != null:
return $default(_that.answer,_that.corrections);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _AmendIntentRequest implements AmendIntentRequest {
  const _AmendIntentRequest({this.answer, this.corrections});
  factory _AmendIntentRequest.fromJson(Map<String, dynamic> json) => _$AmendIntentRequestFromJson(json);

@override final  String? answer;
@override final  IntentCorrections? corrections;

/// Create a copy of AmendIntentRequest
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$AmendIntentRequestCopyWith<_AmendIntentRequest> get copyWith => __$AmendIntentRequestCopyWithImpl<_AmendIntentRequest>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$AmendIntentRequestToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _AmendIntentRequest&&(identical(other.answer, answer) || other.answer == answer)&&(identical(other.corrections, corrections) || other.corrections == corrections));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,answer,corrections);

@override
String toString() {
  return 'AmendIntentRequest(answer: $answer, corrections: $corrections)';
}


}

/// @nodoc
abstract mixin class _$AmendIntentRequestCopyWith<$Res> implements $AmendIntentRequestCopyWith<$Res> {
  factory _$AmendIntentRequestCopyWith(_AmendIntentRequest value, $Res Function(_AmendIntentRequest) _then) = __$AmendIntentRequestCopyWithImpl;
@override @useResult
$Res call({
 String? answer, IntentCorrections? corrections
});


@override $IntentCorrectionsCopyWith<$Res>? get corrections;

}
/// @nodoc
class __$AmendIntentRequestCopyWithImpl<$Res>
    implements _$AmendIntentRequestCopyWith<$Res> {
  __$AmendIntentRequestCopyWithImpl(this._self, this._then);

  final _AmendIntentRequest _self;
  final $Res Function(_AmendIntentRequest) _then;

/// Create a copy of AmendIntentRequest
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? answer = freezed,Object? corrections = freezed,}) {
  return _then(_AmendIntentRequest(
answer: freezed == answer ? _self.answer : answer // ignore: cast_nullable_to_non_nullable
as String?,corrections: freezed == corrections ? _self.corrections : corrections // ignore: cast_nullable_to_non_nullable
as IntentCorrections?,
  ));
}

/// Create a copy of AmendIntentRequest
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$IntentCorrectionsCopyWith<$Res>? get corrections {
    if (_self.corrections == null) {
    return null;
  }

  return $IntentCorrectionsCopyWith<$Res>(_self.corrections!, (value) {
    return _then(_self.copyWith(corrections: value));
  });
}
}

// dart format on
