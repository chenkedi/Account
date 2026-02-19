import 'package:flutter/material.dart';

import '../../../../core/utils/date_utils.dart' as date_utils;
import '../bloc/stats_bloc.dart';

class StatsFilter extends StatelessWidget {
  final DateTime? startDate;
  final DateTime? endDate;
  final StatsPeriod selectedPeriod;
  final ValueChanged<StatsPeriod> onPeriodSelected;
  final VoidCallback onCustomDateTap;

  const StatsFilter({
    super.key,
    this.startDate,
    this.endDate,
    this.selectedPeriod = StatsPeriod.lastSixMonths,
    required this.onPeriodSelected,
    required this.onCustomDateTap,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.all(16),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Row(
              children: [
                Icon(
                  Icons.date_range,
                  color: Theme.of(context).colorScheme.primary,
                ),
                const SizedBox(width: 8),
                Text(
                  '选择时间范围',
                  style: Theme.of(context).textTheme.titleMedium,
                ),
              ],
            ),
            const SizedBox(height: 16),
            Wrap(
              spacing: 8,
              runSpacing: 8,
              children: [
                _PeriodChip(
                  label: '本周',
                  period: StatsPeriod.thisWeek,
                  selected: selectedPeriod == StatsPeriod.thisWeek,
                  onSelected: onPeriodSelected,
                ),
                _PeriodChip(
                  label: '本月',
                  period: StatsPeriod.thisMonth,
                  selected: selectedPeriod == StatsPeriod.thisMonth,
                  onSelected: onPeriodSelected,
                ),
                _PeriodChip(
                  label: '上月',
                  period: StatsPeriod.lastMonth,
                  selected: selectedPeriod == StatsPeriod.lastMonth,
                  onSelected: onPeriodSelected,
                ),
                _PeriodChip(
                  label: '本季度',
                  period: StatsPeriod.thisQuarter,
                  selected: selectedPeriod == StatsPeriod.thisQuarter,
                  onSelected: onPeriodSelected,
                ),
                _PeriodChip(
                  label: '本年',
                  period: StatsPeriod.thisYear,
                  selected: selectedPeriod == StatsPeriod.thisYear,
                  onSelected: onPeriodSelected,
                ),
                _PeriodChip(
                  label: '近6个月',
                  period: StatsPeriod.lastSixMonths,
                  selected: selectedPeriod == StatsPeriod.lastSixMonths,
                  onSelected: onPeriodSelected,
                ),
              ],
            ),
            const SizedBox(height: 16),
            OutlinedButton.icon(
              onPressed: onCustomDateTap,
              icon: const Icon(Icons.edit_calendar),
              label: const Text('自定义日期'),
            ),
            if (startDate != null && endDate != null) ...[
              const SizedBox(height: 12),
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: Theme.of(context).colorScheme.surfaceVariant,
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Text(
                  '${date_utils.DateUtils.formatDate(startDate!)} - ${date_utils.DateUtils.formatDate(endDate!)}',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    color: Theme.of(context).colorScheme.onSurfaceVariant,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }
}

class _PeriodChip extends StatelessWidget {
  final String label;
  final StatsPeriod period;
  final bool selected;
  final ValueChanged<StatsPeriod> onSelected;

  const _PeriodChip({
    required this.label,
    required this.period,
    required this.selected,
    required this.onSelected,
  });

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return FilterChip(
      label: Text(label),
      selected: selected,
      onSelected: (isSelected) {
        if (isSelected) {
          onSelected(period);
        }
      },
      selectedColor: colorScheme.primaryContainer,
      checkmarkColor: colorScheme.onPrimaryContainer,
      labelStyle: TextStyle(
        color: selected
            ? colorScheme.onPrimaryContainer
            : colorScheme.onSurface,
      ),
      backgroundColor: colorScheme.surface,
      side: BorderSide(
        color: selected ? colorScheme.primary : colorScheme.outline,
      ),
    );
  }
}
