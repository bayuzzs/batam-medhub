import 'dart:convert';

/// Reads the `exp` (expiry) claim from an HS256 JWT access token.
///
/// Only the payload segment is decoded — the token is NOT signature-verified
/// (the client has no HS256 secret; it only needs to *read* claims to decide
/// when to refresh). Returns `null` when the token isn't a decodable JWT or
/// carries no integer `exp` claim, so callers can fall back to another
/// source of truth (e.g. `expires_in_seconds`).
int? tryDecodeJwtExp(String token) {
  final segments = token.split('.');
  if (segments.length != 3) {
    return null;
  }
  try {
    final decoded = utf8.decode(
      base64Url.decode(base64Url.normalize(segments[1])),
    );
    final payload = jsonDecode(decoded);
    if (payload is! Map<String, dynamic>) {
      return null;
    }
    final exp = payload['exp'];
    return exp is int ? exp : null;
  } on FormatException {
    return null;
  }
}
