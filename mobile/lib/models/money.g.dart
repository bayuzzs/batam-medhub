// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'money.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_Money _$MoneyFromJson(Map<String, dynamic> json) => _Money(
  amountMinor: (json['amount_minor'] as num).toInt(),
  currency: json['currency'] as String,
);

Map<String, dynamic> _$MoneyToJson(_Money instance) => <String, dynamic>{
  'amount_minor': instance.amountMinor,
  'currency': instance.currency,
};

_ConvertedMoney _$ConvertedMoneyFromJson(Map<String, dynamic> json) =>
    _ConvertedMoney(
      source: Money.fromJson(json['source'] as Map<String, dynamic>),
      display: Money.fromJson(json['display'] as Map<String, dynamic>),
      fxRate: json['fx_rate'] as String,
      fxSource: json['fx_source'] as String,
      fxEffectiveAt: DateTime.parse(json['fx_effective_at'] as String),
      estimated: json['estimated'] as bool,
    );

Map<String, dynamic> _$ConvertedMoneyToJson(_ConvertedMoney instance) =>
    <String, dynamic>{
      'source': instance.source,
      'display': instance.display,
      'fx_rate': instance.fxRate,
      'fx_source': instance.fxSource,
      'fx_effective_at': instance.fxEffectiveAt.toIso8601String(),
      'estimated': instance.estimated,
    };

_PriceSummary _$PriceSummaryFromJson(Map<String, dynamic> json) =>
    _PriceSummary(
      sourceTotals: (json['source_totals'] as List<dynamic>)
          .map((e) => Money.fromJson(e as Map<String, dynamic>))
          .toList(),
      displayTotal: Money.fromJson(
        json['display_total'] as Map<String, dynamic>,
      ),
      estimated: json['estimated'] as bool,
    );

Map<String, dynamic> _$PriceSummaryToJson(_PriceSummary instance) =>
    <String, dynamic>{
      'source_totals': instance.sourceTotals,
      'display_total': instance.displayTotal,
      'estimated': instance.estimated,
    };
