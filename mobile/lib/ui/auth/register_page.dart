import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lucide_icons_flutter/lucide_icons.dart';

import 'package:mobile/application/auth/providers.dart';
import 'package:mobile/ui/core/theme/app_assets.dart';
import 'package:mobile/ui/core/widgets/app_container.dart';
import 'package:mobile/ui/core/navigation/app_router.dart';
import 'package:mobile/ui/core/widgets/app_text_field.dart';
import 'package:mobile/ui/core/widgets/app_validators.dart';
import 'package:mobile/ui/core/widgets/primary_radial_gradient.dart';

/// Register screen.
///
/// Similar layout to [LoginPage] (see `login_page.dart`) but with full name,
/// email, password and confirm password fields. On success the
/// [authControllerProvider] transitions to `authenticated` and the router
/// redirects to the chat shell; failures surface as a snackbar.
class RegisterPage extends ConsumerStatefulWidget {
  const RegisterPage({super.key});

  @override
  ConsumerState<RegisterPage> createState() => _RegisterPageState();
}

class _RegisterPageState extends ConsumerState<RegisterPage> {
  final _formKey = GlobalKey<FormState>();
  final _fullNameController = TextEditingController();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  final _confirmPasswordController = TextEditingController();
  bool _obscurePassword = true;
  bool _obscureConfirmPassword = true;

  @override
  void dispose() {
    _fullNameController.dispose();
    _emailController.dispose();
    _passwordController.dispose();
    _confirmPasswordController.dispose();
    super.dispose();
  }

  Future<void> _register() async {
    if (!(_formKey.currentState?.validate() ?? false)) {
      return;
    }
    final notifier = ref.read(authControllerProvider.notifier);
    final ok = await notifier.register(
      fullName: _fullNameController.text.trim(),
      email: _emailController.text.trim(),
      password: _passwordController.text,
    );
    if (!ok && mounted) {
      final message =
          ref.read(authControllerProvider).errorMessage ??
          'Unable to register. Please try again.';
      ScaffoldMessenger.of(context)
        ..hideCurrentSnackBar()
        ..showSnackBar(SnackBar(content: Text(message)));
    }
  }

  void _goToLogin() {
    context.goLogin();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isSubmitting = ref.watch(
      authControllerProvider.select((state) => state.isSubmitting),
    );

    return Scaffold(
      body: Stack(
        children: [
          const PrimaryRadialGradient(startOpacity: 0.20),
          AppContainer(
            scrollable: true,
            vertical: 100,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                Center(child: Image.asset(AppAssets.logo, height: 64)),
                const SizedBox(height: 32),
                Text(
                  'Create Your Account',
                  textAlign: TextAlign.center,
                  style: theme.textTheme.headlineMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 8),
                Text(
                  'Join Batam MedHub to get started',
                  textAlign: TextAlign.center,
                  style: theme.textTheme.bodyLarge?.copyWith(
                    color: theme.colorScheme.onSurfaceVariant,
                  ),
                ),
                const SizedBox(height: 32),
                Form(
                  key: _formKey,
                  child: Column(
                    children: [
                      AppTextField(
                        label: 'Full Name',
                        controller: _fullNameController,
                        textCapitalization: TextCapitalization.words,
                        textInputAction: TextInputAction.next,
                        prefixIcon: const Icon(Icons.person_outline),
                        validator: (value) => AppValidators.required(
                          value,
                          message: 'Full name is required',
                        ),
                      ),
                      const SizedBox(height: 16),
                      AppTextField(
                        label: 'Email',
                        controller: _emailController,
                        keyboardType: TextInputType.emailAddress,
                        textInputAction: TextInputAction.next,
                        autocorrect: false,
                        prefixIcon: const Icon(LucideIcons.mail),
                        validator: AppValidators.email,
                      ),
                      const SizedBox(height: 16),
                      AppTextField(
                        label: 'Password',
                        controller: _passwordController,
                        obscureText: _obscurePassword,
                        textInputAction: TextInputAction.next,
                        prefixIcon: const Icon(LucideIcons.lock),
                        suffixIcon: IconButton(
                          color: theme.colorScheme.onSurfaceVariant,
                          icon: Icon(
                            _obscurePassword
                                ? LucideIcons.eyeOff
                                : LucideIcons.eye,
                          ),
                          onPressed: () {
                            setState(() {
                              _obscurePassword = !_obscurePassword;
                            });
                          },
                        ),
                        validator: AppValidators.password,
                      ),
                      const SizedBox(height: 16),
                      AppTextField(
                        label: 'Confirm Password',
                        controller: _confirmPasswordController,
                        obscureText: _obscureConfirmPassword,
                        textInputAction: TextInputAction.done,
                        prefixIcon: const Icon(LucideIcons.lock),
                        suffixIcon: IconButton(
                          color: theme.colorScheme.onSurfaceVariant,
                          icon: Icon(
                            _obscureConfirmPassword
                                ? LucideIcons.eyeOff
                                : LucideIcons.eye,
                          ),
                          onPressed: () {
                            setState(() {
                              _obscureConfirmPassword =
                                  !_obscureConfirmPassword;
                            });
                          },
                        ),
                        validator: (value) => AppValidators.confirmPassword(
                          value,
                          expected: _passwordController.text,
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 24),
                FilledButton(
                  onPressed: isSubmitting ? null : _register,
                  child: isSubmitting
                      ? const SizedBox(
                          width: 22,
                          height: 22,
                          child: CircularProgressIndicator(strokeWidth: 2.5),
                        )
                      : const Text('Register'),
                ),
                const SizedBox(height: 16),
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text(
                      'Already have an account?',
                      style: theme.textTheme.bodyMedium?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant,
                      ),
                    ),
                    TextButton(
                      onPressed: _goToLogin,
                      child: const Text('Login'),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
