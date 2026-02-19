import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../../injection_container.dart' as di;
import '../bloc/stats_bloc.dart';
import '../widgets/stats_filter.dart';
import '../widgets/category_chart.dart';
import '../widgets/trend_chart.dart';

class StatsPage extends StatelessWidget {
  const StatsPage({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => di.sl<StatsBloc>()..add(const StatsLoadRequested()),
      child: const StatsView(),
    );
  }
}

class StatsView extends StatelessWidget {
  const StatsView({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('统计报表'),
      ),
      body: BlocBuilder<StatsBloc, StatsState>(
        builder: (context, state) {
          switch (state.status) {
            case StatsStatus.initial:
            case StatsStatus.loading:
              return const Center(child: CircularProgressIndicator());
            case StatsStatus.error:
              return _ErrorView(message: state.errorMessage);
            case StatsStatus.loaded:
              return _StatsContent(state: state);
          }
        },
      ),
    );
  }
}

class _ErrorView extends StatelessWidget {
  final String? message;

  const _ErrorView({this.message});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.error_outline,
            size: 64,
            color: Theme.of(context).colorScheme.error,
          ),
          const SizedBox(height: 16),
          Text(
            message ?? '加载失败',
            style: Theme.of(context).textTheme.titleMedium,
          ),
          const SizedBox(height: 16),
          ElevatedButton(
            onPressed: () {
              context.read<StatsBloc>().add(const StatsLoadRequested());
            },
            child: const Text('重试'),
          ),
        ],
      ),
    );
  }
}

class _StatsContent extends StatefulWidget {
  final StatsState state;

  const _StatsContent({required this.state});

  @override
  State<_StatsContent> createState() => _StatsContentView();
}

class _StatsContentView extends State<_StatsContent> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  StatsPeriod _selectedPeriod = StatsPeriod.lastSixMonths;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        StatsFilter(
          startDate: widget.state.startDate,
          endDate: widget.state.endDate,
          selectedPeriod: _selectedPeriod,
          onPeriodSelected: (period) {
            setState(() {
              _selectedPeriod = period;
            });
            context.read<StatsBloc>().add(StatsPeriodSelected(period));
          },
          onCustomDateTap: () {
            _showDateRangePicker(context);
          },
        ),
        TabBar(
          controller: _tabController,
          tabs: const [
            Tab(text: '分类统计'),
            Tab(text: '收支趋势'),
          ],
        ),
        Expanded(
          child: TabBarView(
            controller: _tabController,
            children: [
              _CategoryTab(state: widget.state),
              _TrendTab(state: widget.state),
            ],
          ),
        ),
      ],
    );
  }

  Future<void> _showDateRangePicker(BuildContext context) async {
    final now = DateTime.now();
    final initialStart = widget.state.startDate ?? now.subtract(const Duration(days: 30));
    final initialEnd = widget.state.endDate ?? now;

    final picked = await showDateRangePicker(
      context: context,
      firstDate: DateTime(2020),
      lastDate: now,
      initialDateRange: DateTimeRange(
        start: initialStart,
        end: initialEnd,
      ),
    );

    if (picked != null && mounted) {
      setState(() {
        _selectedPeriod = StatsPeriod.custom;
      });
      context.read<StatsBloc>().add(
        StatsDateRangeChanged(picked.start, picked.end),
      );
    }
  }
}

class _CategoryTab extends StatelessWidget {
  final StatsState state;

  const _CategoryTab({required this.state});

  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
      length: 2,
      child: Column(
        children: [
          const TabBar(
            tabs: [
              Tab(text: '支出'),
              Tab(text: '收入'),
            ],
          ),
          Expanded(
            child: TabBarView(
              children: [
                SingleChildScrollView(
                  padding: const EdgeInsets.all(16),
                  child: CategoryChart(
                    categoryStats: state.categoryStats,
                    type: 'expense',
                  ),
                ),
                SingleChildScrollView(
                  padding: const EdgeInsets.all(16),
                  child: CategoryChart(
                    categoryStats: state.categoryStats,
                    type: 'income',
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _TrendTab extends StatelessWidget {
  final StatsState state;

  const _TrendTab({required this.state});

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: TrendChart(monthlyTrend: state.monthlyTrend),
    );
  }
}
