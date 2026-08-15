import 'package:freezed_annotation/freezed_annotation.dart';

part 'time_window.freezed.dart';
part 'time_window.g.dart';

/// An inclusive date window (`DateWindow` schema). `from`/`to` are
/// `format: date` calendar dates without a time component.
@freezed
abstract class DateWindow with _$DateWindow {
  const factory DateWindow({required DateTime from, required DateTime to}) =
      _DateWindow;

  factory DateWindow.fromJson(Map<String, dynamic> json) =>
      _$DateWindowFromJson(json);
}

/// A local-time window with explicit IANA time zones for each end
/// (`TimeWindow` schema). Instants are UTC; zones are kept for local display.
@freezed
abstract class TimeWindow with _$TimeWindow {
  const factory TimeWindow({
    @JsonKey(name: 'starts_at') required DateTime startsAt,
    @JsonKey(name: 'ends_at') required DateTime endsAt,
    @JsonKey(name: 'start_time_zone') required String startTimeZone,
    @JsonKey(name: 'end_time_zone') required String endTimeZone,
  }) = _TimeWindow;

  factory TimeWindow.fromJson(Map<String, dynamic> json) =>
      _$TimeWindowFromJson(json);
}
