import 'package:mobile/models/money.dart';

/// Formats [Money] (integer minor units + ISO currency) for display in the
/// UI, using each currency's ISO 4217 exponent for the minor-unit divisor.
///
/// Examples: `SGD 126.58` (minor 12658), `IDR 150.000.000` (minor
/// 150000000, exponent 0), `SGD 251.90`.
abstract final class MoneyFormatter {
  /// ISO 4217 exponent (digits after the decimal) per currency. Defaults to
  /// 2 when unknown.
  static const Map<String, int> _exponents = {
    'SGD': 2,
    'USD': 2,
    'EUR': 2,
    'MYR': 2,
    'GBP': 2,
    'IDR': 0,
    'JPY': 0,
    'KRW': 0,
    'VND': 0,
  };

  static String format(Money money) {
    final exponent = _exponents[money.currency] ?? 2;
    final divisor = _pow10(exponent);
    final units = money.amountMinor ~/ divisor;
    final fraction = (money.amountMinor % divisor).toString().padLeft(
      exponent,
      '0',
    );
    final unitsText = _thousands(units.toString());
    final fractionText = exponent > 0 ? '.$fraction' : '';
    return '${money.currency} $unitsText$fractionText';
  }

  static int _pow10(int n) {
    var result = 1;
    for (var i = 0; i < n; i++) {
      result *= 10;
    }
    return result;
  }

  /// Inserts a `.` thousands separator (Indonesian convention), e.g.
  /// `150000000` → `150.000.000`.
  static String _thousands(String digits) {
    final buffer = StringBuffer();
    for (var i = 0; i < digits.length; i++) {
      if (i > 0 && (digits.length - i) % 3 == 0) {
        buffer.write('.');
      }
      buffer.write(digits[i]);
    }
    return buffer.toString();
  }
}
