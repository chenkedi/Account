import 'package:flutter/material.dart';
import 'package:fl_chart/fl_chart.dart';

import '../../../../core/utils/amount_utils.dart' as amount_utils;
import '../bloc/stats_bloc.dart';

class CategoryChart extends StatelessWidget {
  final List<CategoryStats> categoryStats;
  final String type; // 'income' or 'expense'

  const CategoryChart({
    super.key,
    required this.categoryStats,
    required this.type,
  });

  @override
  Widget build(BuildContext context) {
    final filtered = categoryStats
        .where((stat) => stat.categoryType == type)
        .toList();

    if (filtered.isEmpty) {
      return const Center(
        child: Text('暂无数据'),
      );
    }

    final total = filtered.fold<double>(
      0.0,
      (sum, stat) => sum + stat.totalAmount,
    );

    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        _ChartSummary(
          title: type == 'income' ? '收入统计' : '支出统计',
          amount: total,
          type: type,
        ),
        const SizedBox(height: 24),
        SizedBox(
          height: 200,
          child: _PieChartWidget(
            categoryStats: filtered,
            total: total,
          ),
        ),
        const SizedBox(height: 24),
        _CategoryLegend(categoryStats: filtered),
      ],
    );
  }
}

class _ChartSummary extends StatelessWidget {
  final String title;
  final double amount;
  final String type;

  const _ChartSummary({
    required this.title,
    required this.amount,
    required this.type,
  });

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              title,
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                color: colorScheme.onSurfaceVariant,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              amount_utils.AmountUtils.formatAmount(amount),
              style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                color: type == 'income'
                    ? colorScheme.primary
                    : colorScheme.error,
                fontWeight: FontWeight.bold,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _PieChartWidget extends StatefulWidget {
  final List<CategoryStats> categoryStats;
  final double total;

  const _PieChartWidget({
    required this.categoryStats,
    required this.total,
  });

  @override
  State<_PieChartWidget> createState() => _PieChartWidgetState();
}

class _PieChartWidgetState extends State<_PieChartWidget> {
  int _touchedIndex = -1;

  @override
  Widget build(BuildContext context) {
    return PieChart(
      PieChartData(
        pieTouchData: PieTouchData(
          touchCallback: (FlTouchEvent event, pieTouchResponse) {
            setState(() {
              if (!event.isInterestedForInteractions ||
                  pieTouchResponse == null ||
                  pieTouchResponse.touchedSection == null) {
                _touchedIndex = -1;
 return;
              }
              _touchedIndex = pieTouchResponse
                  .touchedSection!.touchedSectionIndex;
            });
          },
        ),
        borderData: FlBorderData(show: false),
        sectionsSpace: 2,
        centerSpaceRadius: 40,
        sections: _buildSections(),
      ),
    );
  }

  List<PieChartSectionData> _buildSections() {
    final colors = _getChartColors();

    return List.generate(widget.categoryStats.length, (index) {
      final stat = widget.categoryStats[index];
      final isTouched = index == _touchedIndex;
      final radius = isTouched ? 60.0 : 50.0;
      final fontSize = isTouched ? 16.0 : 12.0;

      final percentage = (stat.totalAmount / widget.total * 100).toStringAsFixed(1);

      return PieChartSectionData(
        color: colors[index % colors.length],
        value: stat.totalAmount,
        title: '$percentage%',
        radius: radius,
        titleStyle: TextStyle(
          fontSize: fontSize,
          fontWeight: FontWeight.bold,
          color: Colors.white,
        ),
      );
    });
  }

  List<Color> _getChartColors() {
    return [
      const Color(0xFF4CAF50),
      const Color(0xFF2196F3),
      const Color(0xFFFF9800),
      const Color(0xFF9C27B0),
      const Color(0xFF00BCD4),
      const Color(0xFFFFEB3B),
      const Color(0xFF795548),
      const Color(0xFF607D8B),
      const Color(0xFFE91E63),
      const Color(0xFF009688),
    ];
  }
}

class _CategoryLegend extends StatelessWidget {
  final List<CategoryStats> categoryStats;

  const _CategoryLegend({required this.categoryStats});

  @override
  Widget build(BuildContext context) {
    final colors = _getChartColors();

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              '详情',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 12),
            ...List.generate(categoryStats.length, (index) {
              final stat = categoryStats[index];
              return Padding(
                padding: const EdgeInsets.only(bottom: 8.0),
                child: Row(
                  children: [
                    Container(
                      width: 16,
                      height: 16,
                      decoration: BoxDecoration(
                        color: colors[index % colors.length],
                        borderRadius: BorderRadius.circular(4),
                      ),
                    ),
                    const SizedBox(width: 8),
                    Expanded(
                      child: Text(
                        stat.categoryName,
                        style: Theme.of(context).textTheme.bodyMedium,
                      ),
                    ),
                    Text(
                      '${stat.percentage.toStringAsFixed(1)}%',
                      style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(width: 8),
                    Text(
                      amount_utils.AmountUtils.formatAmount(stat.totalAmount),
                      style: Theme.of(context).textTheme.bodyMedium,
                    ),
                  ],
                ),
              );
            }),
          ],
        ),
      ),
    );
  }

  List<Color> _getChartColors() {
    return [
      const Color(0xFF4CAF50),
      const Color(0xFF2196F3),
      const Color(0xFFFF9800),
      const Color(0xFF9C27B0),
      const Color(0xFF00BCD4),
      const Color(0xFFFFEB3B),
      const Color(0xFF795548),
      const Color(0xFF607D8B),
      const Color(0xFFE91E63),
      const Color(0xFF009688),
    ];
  }
}
