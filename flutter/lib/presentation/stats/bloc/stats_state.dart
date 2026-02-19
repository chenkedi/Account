part of 'stats_bloc.dart';

enum StatsStatus { initial, loading, loaded, error }

class StatsState extends Equatable {
  final StatsStatus status;
  final DateTime? startDate;
  final DateTime? endDate;
  final StatsSummary? summary;
  final List<CategoryStats> categoryStats;
  final List<MonthlyStats> monthlyTrend;
  final String? errorMessage;

  const StatsState({
    this.status = StatsStatus.initial,
    this.startDate,
    this.endDate,
    this.summary,
    this.categoryStats = const [],
    this.monthlyTrend = const [],
    this.errorMessage,
  });

  StatsState copyWith({
    StatsStatus? status,
    DateTime? startDate,
    DateTime? endDate,
    StatsSummary? summary,
    List<CategoryStats>? categoryStats,
    List<MonthlyStats>? monthlyTrend,
    String? errorMessage,
  }) {
    return StatsState(
      status: status ?? this.status,
      startDate: startDate ?? this.startDate,
      endDate: endDate ?? this.endDate,
      summary: summary ?? this.summary,
      categoryStats: categoryStats ?? this.categoryStats,
      monthlyTrend: monthlyTrend ?? this.monthlyTrend,
      errorMessage: errorMessage,
    );
  }

  @override
  List<Object?> get props => [
    status,
    startDate,
    endDate,
    summary,
    categoryStats,
    monthlyTrend,
    errorMessage,
  ];
}

class StatsSummary {
  final double incomeTotal;
  final double expenseTotal;
  final double netTotal;
  final DateTime startDate;
  final DateTime endDate;

  StatsSummary({
    required this.incomeTotal,
    required this.expenseTotal,
    required this.netTotal,
    required this.startDate,
    required this.endDate,
  });

  factory StatsSummary.fromJson(Map<String, dynamic> json) {
    return StatsSummary(
      incomeTotal: json['income_total']?.toDouble() ?? 0.0,
      expenseTotal: json['expense_total']?.toDouble() ?? 0.0,
      netTotal: json['net_total']?.toDouble() ?? 0.0,
      startDate: DateTime.parse(json['start_date']),
      endDate: DateTime.parse(json['end_date']),
    );
  }
}

class CategoryStats {
  final String categoryId;
  final String categoryName;
  final String categoryType;
  final double totalAmount;
  final int transactionCount;
  final double percentage;

  CategoryStats({
    required this.categoryId,
    required this.categoryName,
    required this.categoryType,
    required this.totalAmount,
    required this.transactionCount,
    required this.percentage,
  });

  factory CategoryStats.fromJson(Map<String, dynamic> json) {
    return CategoryStats(
      categoryId: json['category_id'] ?? '',
      categoryName: json['category_name'] ?? '未分类',
      categoryType: json['category_type'] ?? 'expense',
      totalAmount: json['total_amount']?.toDouble() ?? 0.0,
      transactionCount: json['transaction_count'] ?? 0,
      percentage: json['percentage']?.toDouble() ?? 0.0,
    );
  }
}

class MonthlyStats {
  final int year;
  final int month;
  final double incomeTotal;
  final double expenseTotal;
  final double netTotal;

  MonthlyStats({
    required this.year,
    required this.month,
    required this.incomeTotal,
    required this.expenseTotal,
    required this.netTotal,
  });

  String get monthLabel {
    return '$year年$month月';
  }

  factory MonthlyStats.fromJson(Map<String, dynamic> json) {
    return MonthlyStats(
      year: json['year'] ?? 0,
      month: json['month'] ?? 0,
      incomeTotal: json['income_total']?.toDouble() ?? 0.0,
      expenseTotal: json['expense_total']?.toDouble() ?? 0.0,
      netTotal: json['net_total']?.toDouble() ?? 0.0,
    );
  }
}
