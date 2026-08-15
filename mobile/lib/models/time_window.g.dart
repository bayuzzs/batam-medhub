// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'time_window.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_DateWindow _$DateWindowFromJson(Map<String, dynamic> json) => _DateWindow(
  from: DateTime.parse(json['from'] as String),
  to: DateTime.parse(json['to'] as String),
);

Map<String, dynamic> _$DateWindowToJson(_DateWindow instance) =>
    <String, dynamic>{
      'from': instance.from.toIso8601String(),
      'to': instance.to.toIso8601String(),
    };

_TimeWindow _$TimeWindowFromJson(Map<String, dynamic> json) => _TimeWindow(
  startsAt: DateTime.parse(json['starts_at'] as String),
  endsAt: DateTime.parse(json['ends_at'] as String),
  startTimeZone: json['start_time_zone'] as String,
  endTimeZone: json['end_time_zone'] as String,
);

Map<String, dynamic> _$TimeWindowToJson(_TimeWindow instance) =>
    <String, dynamic>{
      'starts_at': instance.startsAt.toIso8601String(),
      'ends_at': instance.endsAt.toIso8601String(),
      'start_time_zone': instance.startTimeZone,
      'end_time_zone': instance.endTimeZone,
    };
