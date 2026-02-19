import 'package:flutter/material.dart';
import 'package:fl_chart/fl_chart.dart';

import '../../../../core/utils/amount_utils.dart' as amount_utils;
import '../bloc/stats_bloc.dart';

class TrendChart extends StatelessWidget {
  final List<MonthlyStats> monthlyTrend;

  const TrendChart({super.key, required this.monthlyTrend});

  @override
  Widget build(BuildContext context) {
    if (monthlyTrend.isEmpty) {
      return const Center(
        child: Text('暂无数据'),
      );
    }

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text(
              '收支趋势',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 24),
            SizedBox(
              height: 200,
              child: _LineChartWidget(monthlyTrend: monthlyTrend),
            ),
            const SizedBox(height: 16),
            const _Legend(),
          ],
        ),
      ),
    );
  }
}

class _LineChartWidget extends StatelessWidget {
  final List<MonthlyStats> monthlyTrend;

  const _LineChartWidget({required this.monthlyTrend});

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return LineChart(
      LineChartData(
        gridData: FlGridData(
          show: true,
          drawVerticalLine: false,
          getDrawingHorizontalLine: (value) {
            return FlLine(
              color: colorScheme.outlineVariant.withOpacity(0.3),
              strokeWidth: 1,
            );
          },
        ),
        titlesData: FlTitlesData(
          rightTitles: const AxisTitles(
            sideTitles: SideTitles(showTitles: false),
          ),
          topTitles: const AxisTitles(
            sideTitles: SideTitles(showTitles: false),
          ),
          bottomTitles: AxisTitles(
            sideTitles: SideTitles(
              showTitles: true,
              getTitlesWidget: (value, meta) {
                final index = value.toInt();
                if (index < 0 || index >= monthlyTrend.length) {
                  return const Text('');
                }
                final stat = monthlyTrend[index];
                return Padding(
                  padding: const EdgeInsets.only(top: 8.0),
                  child: Text(
                    '${stat.month}月',
                    style: Theme.of(context).textTheme.bodySmall,
                  ),
                );
              },
              reservedSize: 30,
              interval: 1,
            ),
          ),
          leftTitles: AxisTitles(
            sideTitles: SideTitles(
              showTitles: true,
              getTitlesWidget: (value, meta) {
                if (value == 0) {
                  return const Text('');
                }
                return Padding(
                  padding: const EdgeInsets.only(right: 8.0),
                  child: Text(
                    '${(value / 1000).toStringAsFixed(0)}k',
                    style: Theme.of(context).textTheme.bodySmall,
                  ),
                );
              },
              reservedSize: 40,
            ),
          ),
        ),
        borderData: FlBorderData(
          show: true,
          border: Border(
            bottom: BorderSide(color: colorScheme.outlineVariant, width: 1),
            left: BorderSide(color: colorScheme.outlineVariant, width: 1),
            right: BorderSide.none,
            top: BorderSide.none,
          ),
        ),
        minX: 0,
        maxX: (monthlyTrend.length - 1).toDouble(),
        minY: 0,
        maxY: _calculateMaxY(),
        lineBarsData: [
          _buildIncomeBar(colorScheme.primary),
          _buildExpenseBar(colorScheme.error),
        ],
        lineTouchData: LineTouchData(
          touchTooltipData: LineTouchTooltipData(
            tooltipBgColor: colorScheme.surfaceVariant,
            getTooltipItems: (touchedSpots) {
              return touchedSpots.map((spot) {
                final index = spot.x.toInt();
                final stat = monthlyTrend[index];
                final isIncome = spot.barIndex == 0;
                final amount = isIncome ? stat.incomeTotal : stat.expenseTotal;

                return LineTooltipItem(
                  amount_utils.AmountUtils.formatAmount(amount),
                  TextStyle(
                    color: isIncome ? colorScheme.primary : colorScheme.error,
                    fontWeight: FontWeight.bold,
                  ),
                );
              }).toList();
            },
          ),
        ),
      ),
    );
  }

  double _calculateMaxY() {
    double max = 0;
    for (final stat in monthlyTrend) {
      if (stat.incomeTotal > max) max = stat.incomeTotal;
      if (stat.expenseTotal > max) max = stat.expenseTotal;
    }
    // Add 20% buffer
    return max * 1.2;
  }

  LineChartBarData _buildIncomeBar(Color color) {
    return LineChartBarData(
      spots: monthlyTrend.asMap().entries.map((entry) {
        final index = entry.key;
        final stat = entry.value;
        return FlSpot(index.toDouble(), stat.incomeTotal);
      }).toList(),
      isCurved: true,
      color: color,
      barWidth: 3,
      isStrokeCapRound: true,
      dotData: const FlDotData(show: true),
      belowBarData: BarAreaData(
        show: true,
        color: color.withOpacity(0.1),
      ),
    );
  }

  LineChartBarData _buildExpenseBar(Color color) {
    return LineChartBarData(
      spots: monthlyTrend.asMap().entries.map((entry) {
        final index = entry.key;
        final stat = entry.value;
        return FlSpot(index.toDouble(), stat.expenseTotal);
      }).toList(),
      isCurved: true,
      color: color,
      barWidth: 3,
      isStrokeCapRound: true,
      dotData: const FlDotData(show: true),
      belowBarData: BarAreaData(
        show: true,
        color: color.withOpacity(0.1),
      ),
    );
  }
}

class _Legend extends StatelessWidget {
  const _Legend();

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        _LegendItem(
          color: colorScheme.primary,
          label: '收入',
        ),
        const SizedBox(width: 24),
        _LegendItem(
          color: colorScheme.error,
          label: '支出',
        ),
      ],
    );
  }
}

class _LegendItem extends StatelessWidget {
  final Color color;
  final String label;

  const _LegendItem({required this.color, required this.label});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Container(
          width: 12,
          height: 12,
          decoration: BoxDecoration(
            color: color,
            shape: BoxShape.circle,
          ),
        ),
        const SizedBox(width: 8),
        Text(
          label,
          style: Theme.of(context).textTheme.bodyMedium,
        ),
      ],
    );
  }
}
