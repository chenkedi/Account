import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../data/models/import.dart' as models;
import '../../../data/models/transaction.dart';
import '../../../core/utils/date_utils.dart' as date_utils;
import '../../../core/utils/amount_utils.dart' as amount_utils;
import '../bloc/import_bloc.dart';

class ImportPreviewWidget extends StatelessWidget {
  const ImportPreviewWidget({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<ImportBloc, ImportState>(
      builder: (context, state) {
        if (state.preview == null) {
          return const Center(child: Text('没有数据'));
        }

        final preview = state.preview!;

        return Column(
          children: [
            _SummaryHeader(preview: preview),
            const Divider(height: 1),
            Expanded(
              child: _TransactionList(
                transactions: preview.transactions,
                categories: preview.categories,
                accountSuggestions: preview.accountSuggestions,
              ),
            ),
            const Divider(height: 1),
            _BottomActions(preview: preview),
          ],
        );
      },
    );
  }
}

class _SummaryHeader extends StatelessWidget {
  final models.ImportPreview preview;

  const _SummaryHeader({required this.preview});

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Container(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '预览导入',
            style: Theme.of(context).textTheme.headlineSmall,
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              _SummaryItem(
                label: '总计',
                value: '${preview.totalRows} 条',
                color: colorScheme.onSurface,
              ),
              const SizedBox(width: 16),
              _SummaryItem(
                label: '可导入',
                value: '${preview.validRows} 条',
                color: colorScheme.primary,
              ),
              const SizedBox(width: 16),
              _SummaryItem(
                label: '重复',
                value: '${preview.duplicateRows} 条',
                color: colorScheme.error,
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _SummaryItem extends StatelessWidget {
  final String label;
  final String value;
  final Color color;

  const _SummaryItem({
    required this.label,
    required this.value,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: Theme.of(context).textTheme.bodySmall?.copyWith(
            color: Theme.of(context).colorScheme.onSurfaceVariant,
          ),
        ),
        const SizedBox(height: 4),
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

class _TransactionList extends StatelessWidget {
  final List<models.ParsedTransaction> transactions;
  final List<dynamic> categories;
  final Map<String, List<dynamic>>? accountSuggestions;

  const _TransactionList({
    required this.transactions,
    required this.categories,
    this.accountSuggestions,
  });

  @override
  Widget build(BuildContext context) {
    return ListView.separated(
      padding: const EdgeInsets.symmetric(vertical: 8),
      itemCount: transactions.length,
      separatorBuilder: (context, index) => const Divider(height: 1),
      itemBuilder: (context, index) {
        return _TransactionTile(
          index: index,
          transaction: transactions[index],
          categories: categories,
        );
      },
    );
  }
}

class _TransactionTile extends StatelessWidget {
  final int index;
  final models.ParsedTransaction transaction;
  final List<dynamic> categories;

  const _TransactionTile({
    required this.index,
    required this.transaction,
    required this.categories,
  });

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Opacity(
      opacity: transaction.canBeImported ? 1.0 : 0.5,
      child: ExpansionTile(
        leading: _buildLeadingIcon(context),
        title: Text(
          transaction.counterparty ?? transaction.note ?? '未知交易',
          maxLines: 1,
          overflow: TextOverflow.ellipsis,
        ),
        subtitle: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              date_utils.DateUtils.formatDateTime(transaction.transactionDate),
              style: Theme.of(context).textTheme.bodySmall,
            ),
            if (transaction.isDuplicate)
              Text(
                '可能重复',
                style: Theme.of(context).textTheme.bodySmall?.copyWith(
                  color: colorScheme.error,
                ),
              ),
          ],
        ),
        trailing: Text(
          amount_utils.AmountUtils.formatWithSign(
            transaction.amount,
            transaction.type,
          ),
          style: Theme.of(context).textTheme.titleMedium?.copyWith(
            color: transaction.type == TransactionType.income
                ? colorScheme.primary
                : colorScheme.error,
            fontWeight: FontWeight.bold,
          ),
        ),
        children: [
          _TransactionDetails(
            index: index,
            transaction: transaction,
            categories: categories,
          ),
        ],
      ),
    );
  }

  Widget _buildLeadingIcon(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    if (!transaction.canBeImported) {
      return Icon(Icons.block, color: colorScheme.error);
    }

    if (transaction.isDuplicate) {
      return Icon(Icons.warning, color: colorScheme.error);
    }

    return Icon(
      transaction.type == TransactionType.income
          ? Icons.add_circle
          : Icons.remove_circle,
      color: transaction.type == TransactionType.income
          ? colorScheme.primary
          : colorScheme.error,
    );
  }
}

class _TransactionDetails extends StatelessWidget {
  final int index;
  final models.ParsedTransaction transaction;
  final List<dynamic> categories;

  const _TransactionDetails({
    required this.index,
    required this.transaction,
    required this.categories,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _DetailRow(
            label: '账户',
            child: _AccountSelector(
              index: index,
              transaction: transaction,
            ),
          ),
          const SizedBox(height: 12),
          _DetailRow(
            label: '备注',
            value: transaction.note,
          ),
          if (transaction.accountName != null) ...[
            const SizedBox(height: 12),
            _DetailRow(
              label: '原始账户',
              value: transaction.accountName,
            ),
          ],
        ],
      ),
    );
  }
}

class _AccountSelector extends StatelessWidget {
  final int index;
  final models.ParsedTransaction transaction;

  const _AccountSelector({
    required this.index,
    required this.transaction,
  });

  @override
  Widget build(BuildContext context) {
    // TODO: Implement account selector dropdown
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surfaceVariant,
        borderRadius: BorderRadius.circular(8),
      ),
      child: Text(
        transaction.selectedAccountId != null
            ? '已选择账户'
            : '请选择账户',
        style: TextStyle(
          color: transaction.selectedAccountId != null
              ? Theme.of(context).colorScheme.onSurface
              : Theme.of(context).colorScheme.error,
        ),
      ),
    );
  }
}

class _DetailRow extends StatelessWidget {
  final String label;
  final String? value;
  final Widget? child;

  const _DetailRow({
    required this.label,
    this.value,
    this.child,
  }) : assert(value != null || child != null);

  @override
  Widget build(BuildContext context) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SizedBox(
          width: 80,
          child: Text(
            label,
            style: Theme.of(context).textTheme.bodyMedium?.copyWith(
              color: Theme.of(context).colorScheme.onSurfaceVariant,
            ),
          ),
        ),
        Expanded(
          child: child ??
              Text(
                value ?? '-',
                style: Theme.of(context).textTheme.bodyMedium,
              ),
        ),
      ],
    );
  }
}

class _BottomActions extends StatelessWidget {
  final models.ImportPreview preview;

  const _BottomActions({required this.preview});

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<ImportBloc, ImportState>(
      builder: (context, state) {
        final canImport = preview.validRows > 0 &&
            preview.transactions.any((tx) =>
                tx.canBeImported && tx.selectedAccountId != null);

        final isLoading = state.importStatus == ImportStepStatus.loading;

        return Container(
          padding: const EdgeInsets.all(16),
          child: SafeArea(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                // TODO: Add "Set account for all" dropdown
                const SizedBox(height: 12),
                SizedBox(
                  width: double.infinity,
                  child: FilledButton(
                    onPressed: isLoading || !canImport
                        ? null
                        : () {
                            context.read<ImportBloc>().add(ImportConfirmed());
                          },
                    child: isLoading
                        ? const SizedBox(
                            height: 20,
                            width: 20,
                            child: CircularProgressIndicator(strokeWidth: 2),
                          )
                        : Text('导入 ${preview.validRows} 条记录'),
                  ),
                ),
              ],
            ),
          ),
        );
      },
    );
  }
}
