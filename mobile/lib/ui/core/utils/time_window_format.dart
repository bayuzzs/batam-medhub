import 'package:mobile/models/time_window.dart';

/// Formats a [TimeWindow]'s start as `22 Aug · 03:00` (device-local time).
String formatWindow(TimeWindow window) {
  const months = [
    'Jan',
    'Feb',
    'Mar',
    'Apr',
    'May',
    'Jun',
    'Jul',
    'Aug',
    'Sep',
    'Oct',
    'Nov',
    'Dec',
  ];
  final local = window.startsAt.toLocal();
  final hh = local.hour.toString().padLeft(2, '0');
  final mm = local.minute.toString().padLeft(2, '0');
  return '${local.day} ${months[local.month - 1]} · $hh:$mm';
}

/// [formatWindow] plus the IANA zone name, e.g. `22 Aug · 03:00 · Asia/Jakarta`.
String formatWindowWithZone(TimeWindow window) =>
    '${formatWindow(window)} · ${window.startTimeZone}';
