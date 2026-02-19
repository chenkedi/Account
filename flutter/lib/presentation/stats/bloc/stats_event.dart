part of 'stats_bloc.dart';

abstract class StatsEvent extends Equatable {
  const StatsEvent();

  @override
  List<Object?> get props => [];
}

class StatsLoadRequested extends StatsEvent {}

class StatsDateRangeChanged extends StatsEvent {
  final DateTime startDate;
  final DateTime endDate;

  const StatsDateRangeChanged(this.startDate, this.endDate);

  @override
  List<Object?> get props => [startDate, endDate];
}

class StatsPeriodSelected extends StatsEvent {
  final StatsPeriod period;

  const StatsPeriodSelected(this.period);

  @override
  List<Object?> get props => [period];
}

enum StatsPeriod {
  thisWeek,
  thisMonth,
  lastMonth,
  thisQuarter,
  thisYear,
  lastSixMonths,
  custom,
}
