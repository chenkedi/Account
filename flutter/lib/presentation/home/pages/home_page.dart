import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../bloc/home_bloc.dart';
import '../../sync/bloc/sync_bloc.dart';
import '../../sync/widgets/sync_status_indicator.dart';
import '../../dashboard/pages/dashboard_page.dart';
import '../../transactions/pages/transactions_page.dart';
import '../../accounts/pages/accounts_page.dart';
import '../../stats/pages/stats_page.dart';
import '../../settings/pages/settings_page.dart';
import '../../widgets/common/bottom_nav_bar.dart';
import '../../widgets/common/app_bar.dart';
import '../../../injection_container.dart' as di;

class HomePage extends StatelessWidget {
  const HomePage({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiBlocProvider(
      providers: [
        BlocProvider(create: (_) => HomeBloc()),
        BlocProvider(
          create: (_) => SyncBloc(syncManager: di.sl()),
        ),
      ],
      child: const HomeView(),
    );
  }
}

class HomeView extends StatelessWidget {
  const HomeView({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<HomeBloc, HomeState>(
      builder: (context, state) {
        return Scaffold(
          appBar: _buildAppBar(context, state.currentTab),
          body: _buildPage(state.currentTab),
          bottomNavigationBar: BottomNavBar(
            currentTab: state.currentTab,
            onTabSelected: (tab) {
              context.read<HomeBloc>().add(HomeTabChanged(tab));
            },
          ),
        );
      },
    );
  }

  PreferredSizeWidget _buildAppBar(BuildContext context, HomeTab tab) {
    final actions = tab != HomeTab.settings
        ? [const SyncStatusIndicator()]
        : null;

    switch (tab) {
      case HomeTab.dashboard:
        return CommonAppBar(title: 'Dashboard', actions: actions);
      case HomeTab.transactions:
        return CommonAppBar(title: 'Transactions', actions: actions);
      case HomeTab.accounts:
        return CommonAppBar(title: 'Accounts', actions: actions);
      case HomeTab.stats:
        return CommonAppBar(title: 'Statistics', actions: actions);
      case HomeTab.settings:
        return const CommonAppBar(title: 'Settings');
    }
  }

  Widget _buildPage(HomeTab tab) {
    switch (tab) {
      case HomeTab.dashboard:
        return const DashboardPage();
      case HomeTab.transactions:
        return const TransactionsPage();
      case HomeTab.accounts:
        return const AccountsPage();
      case HomeTab.stats:
        return const StatsPage();
      case HomeTab.settings:
        return const SettingsPage();
    }
  }
}
