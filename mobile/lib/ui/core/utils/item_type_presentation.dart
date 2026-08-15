import 'package:flutter/material.dart';
import 'package:lucide_icons_flutter/lucide_icons.dart';

import 'package:mobile/models/plan_option.dart';

/// Human label for an [ItemType], e.g. `Ferry outbound`, `Hospital
/// appointment`, `Arrival buffer`. Shared by the plan detail and active
/// itinerary screens.
String itemTypeLabel(ItemType type) {
  switch (type) {
    case ItemType.ferryOutbound:
      return 'Ferry outbound';
    case ItemType.ferryReturn:
      return 'Ferry return';
    case ItemType.transportPickup:
      return 'Transport pickup';
    case ItemType.transportDropoff:
      return 'Transport drop-off';
    case ItemType.hospitalAppointment:
      return 'Hospital appointment';
    case ItemType.additionalCare:
      return 'Additional care';
    case ItemType.hotelStay:
      return 'Hotel stay';
    case ItemType.arrivalBuffer:
      return 'Arrival buffer';
    case ItemType.departureBuffer:
      return 'Departure buffer';
  }
}

/// Leading icon for an [ItemType], shared by the plan detail and active
/// itinerary screens.
IconData itemTypeIcon(ItemType type) {
  switch (type) {
    case ItemType.ferryOutbound:
    case ItemType.ferryReturn:
      return LucideIcons.ship;
    case ItemType.transportPickup:
    case ItemType.transportDropoff:
      return LucideIcons.car;
    case ItemType.hospitalAppointment:
      return LucideIcons.hospital;
    case ItemType.additionalCare:
      return LucideIcons.heartPulse;
    case ItemType.hotelStay:
      return LucideIcons.bedDouble;
    case ItemType.arrivalBuffer:
    case ItemType.departureBuffer:
      return LucideIcons.hourglass;
  }
}
