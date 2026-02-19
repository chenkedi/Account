import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../data/models/import.dart' as models;
import '../bloc/import_bloc.dart';

class ImportResultWidget extends StatelessWidget {
  const ImportResultWidget({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<ImportBloc, ImportState>(
      builder: (context, state) {
        if (state.result == null) {
          return const Center(child: Text('没有结果'));
        }

        final result = state.result!;

        return Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              const Spacer(),
              _ResultIcon(success: result.failedRows == 0),
              const SizedBox(height: 24),
              Text(
                result.failedRows == 0 ? '导入成功！' : '导入完成',
                textAlign: TextAlign.center,
                style: Theme.of(context).textTheme.headlineSmall,
              ),
              const SizedBox(height: 32),
              _ResultSummary(result: result),
              if (result.errors != null && result.errors!.isNotEmpty) ...[
                const SizedBox(height: 24),
                _ErrorList(errors: result.errors!),
              ],
              const Spacer(flex: 2),
              FilledButton(
                onPressed: () {
                  // Navigate back or reset
                  Navigator.of(context).pop();
                },
                child: const Text('完成'),
              ),
            ],
          ),
        );
      },
    );
  }
}

class _ResultIcon extends StatelessWidget {
  final bool success;

  const _ResultIcon({required this.success});

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Container(
      width: 100,
      height: 100,
      decoration: BoxDecoration(
        color: success
            ? colorScheme.primaryContainer
            : colorScheme.errorContainer,
        shape: BoxShape.circle,
      ),
      child: Icon(
        success ? Icons.check : Icons.warning,
        size: 60,
        color: success
            ? colorScheme.onPrimaryContainer
            : colorScheme.onErrorContainer,
      ),
    );
  }
}

class _ResultSummary extends StatelessWidget {
  final models.ImportResult result;

  const _ResultSummary({required this.result});

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          children: [
            _SummaryRow(
              label: '总计',
              value: '${result.totalRows} 条',
              color: colorScheme.onSurface,
            ),
            const SizedBox(height: 12),
            _SummaryRow(
              label: '成功',
              value: '${result.importedRows} 条',
              color: colorScheme.primary,
            ),
            const SizedBox(height: 12),
            _SummaryRow(
              label: '跳过',
              value: '${result.skippedRows} 条',
              color: colorScheme.onSurfaceVariant,
            ),
            if (result.failedRows > 0) ...[
              const SizedBox(height: 12),
              _SummaryRow(
                label: '失败',
                value: '${result.failedRows} 条',
                color: colorScheme.error,
              ),
            ],
          ],
        ),
      ),
    );
  }
}

class _SummaryRow extends StatelessWidget {
  final String label;
  final String value;
  final Color color;

  const _SummaryRow({
    required this.label,
    required this.value,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          label,
          style: Theme.of(context).textTheme.bodyMedium,
        ),
        Text(
          value,
          style: Theme.of(context).textTheme.titleMedium?.copyWith(
            color: color,
            fontWeight: FontWeight.bold,
          ),
        ),
      ],
    );
  }
}

class _ErrorList extends StatelessWidget {
  final List<models.ImportError> errors;

  const _ErrorList({required this.errors});

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Card(
      color: colorScheme.errorContainer,
      child: ExpansionTile(
        title: Text(
          '${errors.length} 条错误',
          style: TextStyle(
            color: colorScheme.onErrorContainer,
          ),
        ),
        leading: Icon(
          Icons.error_outline,
          color: colorScheme.onErrorContainer,
        ),
        children: [
          ListView.separated(
            shrinkWrap: true,
            padding: const EdgeInsets.symmetric(horizontal: 16),
            itemCount: errors.length,
            separatorBuilder: (context, index) => const Divider(),
            itemBuilder: (context, index) {
              final error = errors[index];
              return ListTile(
                dense: true,
                contentPadding: EdgeInsets.zero,
                title: Text(
                  '第 ${error.lineNumber} 行',
                  style: TextStyle(
                    color: colorScheme.onErrorContainer,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                subtitle: Text(
                  error.error,
                  style: TextStyle(
                    color: colorScheme.onErrorContainer.withOpacity(0.8),
                  ),
                ),
              );
            },
          ),
        ],
      ),
    );
  }
}
