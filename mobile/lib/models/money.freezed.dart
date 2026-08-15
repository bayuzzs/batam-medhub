// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint, type=warning, deprecated_member_use, deprecated_member_use_from_same_package
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'money.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// dart format off
T _$identity<T>(T value) => value;

/// @nodoc
mixin _$Money {

@JsonKey(name: 'amount_minor') int get amountMinor; String get currency;
/// Create a copy of Money
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$MoneyCopyWith<Money> get copyWith => _$MoneyCopyWithImpl<Money>(this as Money, _$identity);

  /// Serializes this Money to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Money&&(identical(other.amountMinor, amountMinor) || other.amountMinor == amountMinor)&&(identical(other.currency, currency) || other.currency == currency));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,amountMinor,currency);

@override
String toString() {
  return 'Money(amountMinor: $amountMinor, currency: $currency)';
}


}

/// @nodoc
abstract mixin class $MoneyCopyWith<$Res>  {
  factory $MoneyCopyWith(Money value, $Res Function(Money) _then) = _$MoneyCopyWithImpl;
@useResult
$Res call({
@JsonKey(name: 'amount_minor') int amountMinor, String currency
});




}
/// @nodoc
class _$MoneyCopyWithImpl<$Res>
    implements $MoneyCopyWith<$Res> {
  _$MoneyCopyWithImpl(this._self, this._then);

  final Money _self;
  final $Res Function(Money) _then;

/// Create a copy of Money
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? amountMinor = null,Object? currency = null,}) {
  return _then(Money(
amountMinor: null == amountMinor ? _self.amountMinor : amountMinor // ignore: cast_nullable_to_non_nullable
as int,currency: null == currency ? _self.currency : currency // ignore: cast_nullable_to_non_nullable
as String,
  ));
}

}


/// Adds pattern-matching-related methods to [Money].
extension MoneyPatterns on Money {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Money value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Money() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Money value)  $default,){
final _that = this;
switch (_that) {
case _Money():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Money value)?  $default,){
final _that = this;
switch (_that) {
case _Money() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function(@JsonKey(name: 'amount_minor')  int amountMinor,  String currency)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Money() when $default != null:
return $default(_that.amountMinor,_that.currency);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function(@JsonKey(name: 'amount_minor')  int amountMinor,  String currency)  $default,) {final _that = this;
switch (_that) {
case _Money():
return $default(_that.amountMinor,_that.currency);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function(@JsonKey(name: 'amount_minor')  int amountMinor,  String currency)?  $default,) {final _that = this;
switch (_that) {
case _Money() when $default != null:
return $default(_that.amountMinor,_that.currency);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _Money implements Money {
  const _Money({@JsonKey(name: 'amount_minor') required this.amountMinor, required this.currency});
  factory _Money.fromJson(Map<String, dynamic> json) => _$MoneyFromJson(json);

@override@JsonKey(name: 'amount_minor') final  int amountMinor;
@override final  String currency;

/// Create a copy of Money
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$MoneyCopyWith<_Money> get copyWith => __$MoneyCopyWithImpl<_Money>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$MoneyToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _Money&&(identical(other.amountMinor, amountMinor) || other.amountMinor == amountMinor)&&(identical(other.currency, currency) || other.currency == currency));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,amountMinor,currency);

@override
String toString() {
  return 'Money(amountMinor: $amountMinor, currency: $currency)';
}


}

/// @nodoc
abstract mixin class _$MoneyCopyWith<$Res> implements $MoneyCopyWith<$Res> {
  factory _$MoneyCopyWith(_Money value, $Res Function(_Money) _then) = __$MoneyCopyWithImpl;
@override @useResult
$Res call({
@JsonKey(name: 'amount_minor') int amountMinor, String currency
});




}
/// @nodoc
class __$MoneyCopyWithImpl<$Res>
    implements _$MoneyCopyWith<$Res> {
  __$MoneyCopyWithImpl(this._self, this._then);

  final _Money _self;
  final $Res Function(_Money) _then;

/// Create a copy of Money
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? amountMinor = null,Object? currency = null,}) {
  return _then(_Money(
amountMinor: null == amountMinor ? _self.amountMinor : amountMinor // ignore: cast_nullable_to_non_nullable
as int,currency: null == currency ? _self.currency : currency // ignore: cast_nullable_to_non_nullable
as String,
  ));
}


}


/// @nodoc
mixin _$ConvertedMoney {

 Money get source; Money get display;@JsonKey(name: 'fx_rate') String get fxRate;@JsonKey(name: 'fx_source') String get fxSource;@JsonKey(name: 'fx_effective_at') DateTime get fxEffectiveAt; bool get estimated;
/// Create a copy of ConvertedMoney
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$ConvertedMoneyCopyWith<ConvertedMoney> get copyWith => _$ConvertedMoneyCopyWithImpl<ConvertedMoney>(this as ConvertedMoney, _$identity);

  /// Serializes this ConvertedMoney to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is ConvertedMoney&&(identical(other.source, source) || other.source == source)&&(identical(other.display, display) || other.display == display)&&(identical(other.fxRate, fxRate) || other.fxRate == fxRate)&&(identical(other.fxSource, fxSource) || other.fxSource == fxSource)&&(identical(other.fxEffectiveAt, fxEffectiveAt) || other.fxEffectiveAt == fxEffectiveAt)&&(identical(other.estimated, estimated) || other.estimated == estimated));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,source,display,fxRate,fxSource,fxEffectiveAt,estimated);

@override
String toString() {
  return 'ConvertedMoney(source: $source, display: $display, fxRate: $fxRate, fxSource: $fxSource, fxEffectiveAt: $fxEffectiveAt, estimated: $estimated)';
}


}

/// @nodoc
abstract mixin class $ConvertedMoneyCopyWith<$Res>  {
  factory $ConvertedMoneyCopyWith(ConvertedMoney value, $Res Function(ConvertedMoney) _then) = _$ConvertedMoneyCopyWithImpl;
@useResult
$Res call({
 Money source, Money display,@JsonKey(name: 'fx_rate') String fxRate,@JsonKey(name: 'fx_source') String fxSource,@JsonKey(name: 'fx_effective_at') DateTime fxEffectiveAt, bool estimated
});


$MoneyCopyWith<$Res> get source;$MoneyCopyWith<$Res> get display;

}
/// @nodoc
class _$ConvertedMoneyCopyWithImpl<$Res>
    implements $ConvertedMoneyCopyWith<$Res> {
  _$ConvertedMoneyCopyWithImpl(this._self, this._then);

  final ConvertedMoney _self;
  final $Res Function(ConvertedMoney) _then;

/// Create a copy of ConvertedMoney
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? source = null,Object? display = null,Object? fxRate = null,Object? fxSource = null,Object? fxEffectiveAt = null,Object? estimated = null,}) {
  return _then(ConvertedMoney(
source: null == source ? _self.source : source // ignore: cast_nullable_to_non_nullable
as Money,display: null == display ? _self.display : display // ignore: cast_nullable_to_non_nullable
as Money,fxRate: null == fxRate ? _self.fxRate : fxRate // ignore: cast_nullable_to_non_nullable
as String,fxSource: null == fxSource ? _self.fxSource : fxSource // ignore: cast_nullable_to_non_nullable
as String,fxEffectiveAt: null == fxEffectiveAt ? _self.fxEffectiveAt : fxEffectiveAt // ignore: cast_nullable_to_non_nullable
as DateTime,estimated: null == estimated ? _self.estimated : estimated // ignore: cast_nullable_to_non_nullable
as bool,
  ));
}
/// Create a copy of ConvertedMoney
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$MoneyCopyWith<$Res> get source {
  
  return $MoneyCopyWith<$Res>(_self.source, (value) {
    return _then(_self.copyWith(source: value));
  });
}/// Create a copy of ConvertedMoney
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$MoneyCopyWith<$Res> get display {
  
  return $MoneyCopyWith<$Res>(_self.display, (value) {
    return _then(_self.copyWith(display: value));
  });
}
}


/// Adds pattern-matching-related methods to [ConvertedMoney].
extension ConvertedMoneyPatterns on ConvertedMoney {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _ConvertedMoney value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _ConvertedMoney() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _ConvertedMoney value)  $default,){
final _that = this;
switch (_that) {
case _ConvertedMoney():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _ConvertedMoney value)?  $default,){
final _that = this;
switch (_that) {
case _ConvertedMoney() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( Money source,  Money display, @JsonKey(name: 'fx_rate')  String fxRate, @JsonKey(name: 'fx_source')  String fxSource, @JsonKey(name: 'fx_effective_at')  DateTime fxEffectiveAt,  bool estimated)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _ConvertedMoney() when $default != null:
return $default(_that.source,_that.display,_that.fxRate,_that.fxSource,_that.fxEffectiveAt,_that.estimated);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( Money source,  Money display, @JsonKey(name: 'fx_rate')  String fxRate, @JsonKey(name: 'fx_source')  String fxSource, @JsonKey(name: 'fx_effective_at')  DateTime fxEffectiveAt,  bool estimated)  $default,) {final _that = this;
switch (_that) {
case _ConvertedMoney():
return $default(_that.source,_that.display,_that.fxRate,_that.fxSource,_that.fxEffectiveAt,_that.estimated);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( Money source,  Money display, @JsonKey(name: 'fx_rate')  String fxRate, @JsonKey(name: 'fx_source')  String fxSource, @JsonKey(name: 'fx_effective_at')  DateTime fxEffectiveAt,  bool estimated)?  $default,) {final _that = this;
switch (_that) {
case _ConvertedMoney() when $default != null:
return $default(_that.source,_that.display,_that.fxRate,_that.fxSource,_that.fxEffectiveAt,_that.estimated);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _ConvertedMoney implements ConvertedMoney {
  const _ConvertedMoney({required this.source, required this.display, @JsonKey(name: 'fx_rate') required this.fxRate, @JsonKey(name: 'fx_source') required this.fxSource, @JsonKey(name: 'fx_effective_at') required this.fxEffectiveAt, required this.estimated});
  factory _ConvertedMoney.fromJson(Map<String, dynamic> json) => _$ConvertedMoneyFromJson(json);

@override final  Money source;
@override final  Money display;
@override@JsonKey(name: 'fx_rate') final  String fxRate;
@override@JsonKey(name: 'fx_source') final  String fxSource;
@override@JsonKey(name: 'fx_effective_at') final  DateTime fxEffectiveAt;
@override final  bool estimated;

/// Create a copy of ConvertedMoney
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$ConvertedMoneyCopyWith<_ConvertedMoney> get copyWith => __$ConvertedMoneyCopyWithImpl<_ConvertedMoney>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$ConvertedMoneyToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _ConvertedMoney&&(identical(other.source, source) || other.source == source)&&(identical(other.display, display) || other.display == display)&&(identical(other.fxRate, fxRate) || other.fxRate == fxRate)&&(identical(other.fxSource, fxSource) || other.fxSource == fxSource)&&(identical(other.fxEffectiveAt, fxEffectiveAt) || other.fxEffectiveAt == fxEffectiveAt)&&(identical(other.estimated, estimated) || other.estimated == estimated));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,source,display,fxRate,fxSource,fxEffectiveAt,estimated);

@override
String toString() {
  return 'ConvertedMoney(source: $source, display: $display, fxRate: $fxRate, fxSource: $fxSource, fxEffectiveAt: $fxEffectiveAt, estimated: $estimated)';
}


}

/// @nodoc
abstract mixin class _$ConvertedMoneyCopyWith<$Res> implements $ConvertedMoneyCopyWith<$Res> {
  factory _$ConvertedMoneyCopyWith(_ConvertedMoney value, $Res Function(_ConvertedMoney) _then) = __$ConvertedMoneyCopyWithImpl;
@override @useResult
$Res call({
 Money source, Money display,@JsonKey(name: 'fx_rate') String fxRate,@JsonKey(name: 'fx_source') String fxSource,@JsonKey(name: 'fx_effective_at') DateTime fxEffectiveAt, bool estimated
});


@override $MoneyCopyWith<$Res> get source;@override $MoneyCopyWith<$Res> get display;

}
/// @nodoc
class __$ConvertedMoneyCopyWithImpl<$Res>
    implements _$ConvertedMoneyCopyWith<$Res> {
  __$ConvertedMoneyCopyWithImpl(this._self, this._then);

  final _ConvertedMoney _self;
  final $Res Function(_ConvertedMoney) _then;

/// Create a copy of ConvertedMoney
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? source = null,Object? display = null,Object? fxRate = null,Object? fxSource = null,Object? fxEffectiveAt = null,Object? estimated = null,}) {
  return _then(_ConvertedMoney(
source: null == source ? _self.source : source // ignore: cast_nullable_to_non_nullable
as Money,display: null == display ? _self.display : display // ignore: cast_nullable_to_non_nullable
as Money,fxRate: null == fxRate ? _self.fxRate : fxRate // ignore: cast_nullable_to_non_nullable
as String,fxSource: null == fxSource ? _self.fxSource : fxSource // ignore: cast_nullable_to_non_nullable
as String,fxEffectiveAt: null == fxEffectiveAt ? _self.fxEffectiveAt : fxEffectiveAt // ignore: cast_nullable_to_non_nullable
as DateTime,estimated: null == estimated ? _self.estimated : estimated // ignore: cast_nullable_to_non_nullable
as bool,
  ));
}

/// Create a copy of ConvertedMoney
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$MoneyCopyWith<$Res> get source {
  
  return $MoneyCopyWith<$Res>(_self.source, (value) {
    return _then(_self.copyWith(source: value));
  });
}/// Create a copy of ConvertedMoney
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$MoneyCopyWith<$Res> get display {
  
  return $MoneyCopyWith<$Res>(_self.display, (value) {
    return _then(_self.copyWith(display: value));
  });
}
}


/// @nodoc
mixin _$PriceSummary {

@JsonKey(name: 'source_totals') List<Money> get sourceTotals;@JsonKey(name: 'display_total') Money get displayTotal; bool get estimated;
/// Create a copy of PriceSummary
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$PriceSummaryCopyWith<PriceSummary> get copyWith => _$PriceSummaryCopyWithImpl<PriceSummary>(this as PriceSummary, _$identity);

  /// Serializes this PriceSummary to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is PriceSummary&&const DeepCollectionEquality().equals(other.sourceTotals, sourceTotals)&&(identical(other.displayTotal, displayTotal) || other.displayTotal == displayTotal)&&(identical(other.estimated, estimated) || other.estimated == estimated));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,const DeepCollectionEquality().hash(sourceTotals),displayTotal,estimated);

@override
String toString() {
  return 'PriceSummary(sourceTotals: $sourceTotals, displayTotal: $displayTotal, estimated: $estimated)';
}


}

/// @nodoc
abstract mixin class $PriceSummaryCopyWith<$Res>  {
  factory $PriceSummaryCopyWith(PriceSummary value, $Res Function(PriceSummary) _then) = _$PriceSummaryCopyWithImpl;
@useResult
$Res call({
@JsonKey(name: 'source_totals') List<Money> sourceTotals,@JsonKey(name: 'display_total') Money displayTotal, bool estimated
});


$MoneyCopyWith<$Res> get displayTotal;

}
/// @nodoc
class _$PriceSummaryCopyWithImpl<$Res>
    implements $PriceSummaryCopyWith<$Res> {
  _$PriceSummaryCopyWithImpl(this._self, this._then);

  final PriceSummary _self;
  final $Res Function(PriceSummary) _then;

/// Create a copy of PriceSummary
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? sourceTotals = null,Object? displayTotal = null,Object? estimated = null,}) {
  return _then(PriceSummary(
sourceTotals: null == sourceTotals ? _self.sourceTotals : sourceTotals // ignore: cast_nullable_to_non_nullable
as List<Money>,displayTotal: null == displayTotal ? _self.displayTotal : displayTotal // ignore: cast_nullable_to_non_nullable
as Money,estimated: null == estimated ? _self.estimated : estimated // ignore: cast_nullable_to_non_nullable
as bool,
  ));
}
/// Create a copy of PriceSummary
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$MoneyCopyWith<$Res> get displayTotal {
  
  return $MoneyCopyWith<$Res>(_self.displayTotal, (value) {
    return _then(_self.copyWith(displayTotal: value));
  });
}
}


/// Adds pattern-matching-related methods to [PriceSummary].
extension PriceSummaryPatterns on PriceSummary {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _PriceSummary value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _PriceSummary() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _PriceSummary value)  $default,){
final _that = this;
switch (_that) {
case _PriceSummary():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _PriceSummary value)?  $default,){
final _that = this;
switch (_that) {
case _PriceSummary() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function(@JsonKey(name: 'source_totals')  List<Money> sourceTotals, @JsonKey(name: 'display_total')  Money displayTotal,  bool estimated)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _PriceSummary() when $default != null:
return $default(_that.sourceTotals,_that.displayTotal,_that.estimated);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function(@JsonKey(name: 'source_totals')  List<Money> sourceTotals, @JsonKey(name: 'display_total')  Money displayTotal,  bool estimated)  $default,) {final _that = this;
switch (_that) {
case _PriceSummary():
return $default(_that.sourceTotals,_that.displayTotal,_that.estimated);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function(@JsonKey(name: 'source_totals')  List<Money> sourceTotals, @JsonKey(name: 'display_total')  Money displayTotal,  bool estimated)?  $default,) {final _that = this;
switch (_that) {
case _PriceSummary() when $default != null:
return $default(_that.sourceTotals,_that.displayTotal,_that.estimated);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _PriceSummary implements PriceSummary {
  const _PriceSummary({@JsonKey(name: 'source_totals') required  List<Money> sourceTotals, @JsonKey(name: 'display_total') required this.displayTotal, required this.estimated}): _sourceTotals = sourceTotals;
  factory _PriceSummary.fromJson(Map<String, dynamic> json) => _$PriceSummaryFromJson(json);

 final  List<Money> _sourceTotals;
@override@JsonKey(name: 'source_totals') List<Money> get sourceTotals {
  if (_sourceTotals is EqualUnmodifiableListView) return _sourceTotals;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_sourceTotals);
}

@override@JsonKey(name: 'display_total') final  Money displayTotal;
@override final  bool estimated;

/// Create a copy of PriceSummary
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$PriceSummaryCopyWith<_PriceSummary> get copyWith => __$PriceSummaryCopyWithImpl<_PriceSummary>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$PriceSummaryToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _PriceSummary&&const DeepCollectionEquality().equals(other._sourceTotals, _sourceTotals)&&(identical(other.displayTotal, displayTotal) || other.displayTotal == displayTotal)&&(identical(other.estimated, estimated) || other.estimated == estimated));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,const DeepCollectionEquality().hash(_sourceTotals),displayTotal,estimated);

@override
String toString() {
  return 'PriceSummary(sourceTotals: $sourceTotals, displayTotal: $displayTotal, estimated: $estimated)';
}


}

/// @nodoc
abstract mixin class _$PriceSummaryCopyWith<$Res> implements $PriceSummaryCopyWith<$Res> {
  factory _$PriceSummaryCopyWith(_PriceSummary value, $Res Function(_PriceSummary) _then) = __$PriceSummaryCopyWithImpl;
@override @useResult
$Res call({
@JsonKey(name: 'source_totals') List<Money> sourceTotals,@JsonKey(name: 'display_total') Money displayTotal, bool estimated
});


@override $MoneyCopyWith<$Res> get displayTotal;

}
/// @nodoc
class __$PriceSummaryCopyWithImpl<$Res>
    implements _$PriceSummaryCopyWith<$Res> {
  __$PriceSummaryCopyWithImpl(this._self, this._then);

  final _PriceSummary _self;
  final $Res Function(_PriceSummary) _then;

/// Create a copy of PriceSummary
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? sourceTotals = null,Object? displayTotal = null,Object? estimated = null,}) {
  return _then(_PriceSummary(
sourceTotals: null == sourceTotals ? _self._sourceTotals : sourceTotals // ignore: cast_nullable_to_non_nullable
as List<Money>,displayTotal: null == displayTotal ? _self.displayTotal : displayTotal // ignore: cast_nullable_to_non_nullable
as Money,estimated: null == estimated ? _self.estimated : estimated // ignore: cast_nullable_to_non_nullable
as bool,
  ));
}

/// Create a copy of PriceSummary
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$MoneyCopyWith<$Res> get displayTotal {
  
  return $MoneyCopyWith<$Res>(_self.displayTotal, (value) {
    return _then(_self.copyWith(displayTotal: value));
  });
}
}

// dart format on
