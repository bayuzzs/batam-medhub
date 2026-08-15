import 'package:flutter/material.dart';

import 'package:mobile/ui/core/app_assets.dart';
import 'package:mobile/ui/core/app_container.dart';
import 'package:mobile/ui/core/app_router.dart';
import 'package:mobile/ui/core/primary_radial_gradient.dart';

/// First screen shown when the app launches.
///
/// Shows the app logo, a tagline, and a "Get Started" button that routes to
/// the login page. No app logic yet — pure UI.
class OnboardingPage extends StatelessWidget {
  const OnboardingPage({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      body: Stack(
        children: [
          const PrimaryRadialGradient(),
          AppContainer(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                Spacer(),
                Center(child: Image.asset(AppAssets.logo, height: 100)),
                const SizedBox(height: 40),
                Text(
                  'Your Health Journey, Orchestrated',
                  textAlign: TextAlign.center,
                  style: theme.textTheme.bodyMedium?.copyWith(
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const Spacer(),
                FilledButton(
                  onPressed: () => context.pushLogin(),
                  child: const Text('Get Started'),
                ),
                SizedBox(height: 32),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
