/// Centralized form validators.
///
/// Keeps validation rules and error messages in one place so login, register,
/// and future forms stay consistent. Pass these straight to
/// `AppTextField.validator` / `TextFormField.validator`.
abstract final class AppValidators {
  static final RegExp _emailRegex = RegExp(r'^[^@\s]+@[^@\s]+\.[^@\s]+$');

  /// Returns an error message when [value] is empty/whitespace, else `null`.
  static String? required(
    String? value, {
    String message = 'This field is required',
  }) {
    return (value ?? '').trim().isEmpty ? message : null;
  }

  /// Validates an email address (required + format).
  static String? email(String? value) {
    final v = value?.trim() ?? '';
    if (v.isEmpty) {
      return 'Email is required';
    }
    if (!_emailRegex.hasMatch(v)) {
      return 'Enter a valid email';
    }
    return null;
  }

  /// Validates a password (required + optional minimum length).
  static String? password(String? value, {int minLength = 6}) {
    final v = value ?? '';
    if (v.isEmpty) {
      return 'Password is required';
    }
    if (minLength > 0 && v.length < minLength) {
      return 'Password must be at least $minLength characters';
    }
    return null;
  }

  /// Validates a confirmation field matches [expected].
  static String? confirmPassword(String? value, {required String expected}) {
    if ((value ?? '').isEmpty) {
      return 'Confirm your password';
    }
    if (value != expected) {
      return 'Passwords do not match';
    }
    return null;
  }
}
