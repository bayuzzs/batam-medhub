/// `hospital-demo-01` → `Hospital Demo 01`; empty/null → `Medical provider`.
String formatProviderLabel(String? providerId) {
  if (providerId == null || providerId.isEmpty) return 'Medical provider';
  return providerId
      .split('-')
      .map((p) => p.isEmpty ? '' : p[0].toUpperCase() + p.substring(1))
      .join(' ');
}
