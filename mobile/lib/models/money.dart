import 'package:freezed_annotation/freezed_annotation.dart';

part 'money.freezed.dart';
part 'money.g.dart';

/// Integer minor units with an ISO 4217 currency code (`Money` schema).
///
/// Amounts are always integer minor units using the currency's ISO 4217
/// exponent (e.g. 12500 SGD cents); never binary floating-point values.
@freezed
abstract class Money with _$Money {
  const factory Money({
    @JsonKey(name: 'amount_minor') required int amountMinor,
    required String currency,
  }) = _Money;

  factory Money.fromJson(Map<String, dynamic> json) => _$MoneyFromJson(json);
}

/// A provider price converted into the patient's display currency
/// (`ConvertedMoney` schema). [fxRate] is a decimal string, never a binary
/// floating-point value.
@freezed
abstract class ConvertedMoney with _$ConvertedMoney {
  const factory ConvertedMoney({
    required Money source,
    required Money display,
    @JsonKey(name: 'fx_rate') required String fxRate,
    @JsonKey(name: 'fx_source') required String fxSource,
    @JsonKey(name: 'fx_effective_at') required DateTime fxEffectiveAt,
    required bool estimated,
  }) = _ConvertedMoney;

  factory ConvertedMoney.fromJson(Map<String, dynamic> json) =>
      _$ConvertedMoneyFromJson(json);
}

/// Aggregated price for a plan option or itinerary version
/// (`PriceSummary` schema): subtotals per original provider currency plus a
/// display total in the patient's reference currency.
@freezed
abstract class PriceSummary with _$PriceSummary {
  const factory PriceSummary({
    @JsonKey(name: 'source_totals') required List<Money> sourceTotals,
    @JsonKey(name: 'display_total') required Money displayTotal,
    required bool estimated,
  }) = _PriceSummary;

  factory PriceSummary.fromJson(Map<String, dynamic> json) =>
      _$PriceSummaryFromJson(json);
}
