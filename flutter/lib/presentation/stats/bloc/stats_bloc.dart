import 'package:bloc/bloc.dart';
import 'package:equatable/equatable.dart';

import '../../../../data/datasources/remote/api_service.dart';
import '../../../../injection_container.dart' as di;
import '../../../../core/utils/date_utils.dart' as date_utils;

part 'stats_event.dart';
part 'stats_state.dart';

class StatsBloc extends Bloc<StatsEvent, StatsState> {
  final ApiService _apiService;

  StatsBloc({ApiService? apiService})
      : _apiService = apiService ?? di.sl<ApiService>(),
        super(const StatsState()) {
    on<StatsLoadRequested>(_onLoadRequested);
    on<StatsDateRangeChanged>(_onDateRangeChanged);
    on<StatsPeriodSelected>(_onPeriodSelected);
  }

  Future<void> _onLoadRequested(
    StatsLoadRequested event,
    Emitter<StatsState> emit,
  ) async {
    // If no date range set, use default (last 6 months)
    if (state.startDate == null || state.endDate == null) {
      final now = DateTime.now();
      final start = now.subtract(const Duration(days: 180));
      await _loadStats(emit, start, now);
    } else {
      await _loadStats(emit, state.startDate!, state.endDate!);
    }
  }

  Future<void> _onDateRangeChanged(
    StatsDateRangeChanged event,
    Emitter<StatsState> emit,
  ) async {
    await _loadStats(emit, event.startDate, event.endDate);
  }

  Future<void> _onPeriodSelected(
    StatsPeriodSelected event,
    Emitter<StatsState> emit,
  ) async {
    final now = DateTime.now();
    DateTime start;
    DateTime end = now;

    switch (event.period) {
      case StatsPeriod.thisWeek:
        start = date_utils.DateUtils.startOfDay(
          now.subtract(Duration(days: now.weekday - 1)),
        );
        break;
      case StatsPeriod.thisMonth:
        start = DateTime(now.year, now.month, 1);
        break;
      case StatsPeriod.lastMonth:
        final lastMonth = DateTime(now.year, now.month, 0);
        start = DateTime(lastMonth.year, lastMonth.month, 1);
        end = DateTime(lastMonth.year, lastMonth.month + 1, 0, 23, 59, 59);
        break;
      case StatsPeriod.thisQuarter:
        final quarter = ((now.month - 1) / 3).floor();
        start = DateTime(now.year, quarter * 3 + 1, 1);
        break;
      case StatsPeriod.thisYear:
        start = DateTime(now.year, 1, 1);
        break;
      case StatsPeriod.lastSixMonths:
        start = now.subtract(const Duration(days: 180));
        break;
      case StatsPeriod.custom:
        // Keep current range
        return;
    }

    await _loadStats(emit, start, end);
  }

  Future<void> _loadStats(
    Emitter<StatsState> emit,
    DateTime start,
    DateTime end,
  ) async {
    emit(state.copyWith(
      status: StatsStatus.loading,
      startDate: start,
      endDate: end,
    ));

    try {
      final data = await _apiService.getDetailedStats(start, end);

      final summary = StatsSummary.fromJson(data['summary']);
      final categoryStats = (data['by_category'] as List)
          .map((json) => CategoryStats.fromJson(json))
          .toList();
      final monthlyTrend = (data['monthly_trend'] as List)
          .map((json) => MonthlyStats.fromJson(json))
          .toList();

      emit(state.copyWith(
        status: StatsStatus.loaded,
        summary: summary,
        categoryStats: categoryStats,
        monthlyTrend: monthlyTrend,
      ));
    } catch (e) {
      emit(state.copyWith(
        status: StatsStatus.error,
        errorMessage: e.toString(),
      ));
    }
  }
}
