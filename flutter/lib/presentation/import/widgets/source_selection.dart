import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../data/models/import.dart';
import '../bloc/import_bloc.dart';

class SourceSelection extends StatelessWidget {
  const SourceSelection({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<ImportBloc, ImportState>(
      builder: (context, state) {
        if (state.sourceStatus == ImportStepStatus.loading) {
          return const Center(child: CircularProgressIndicator());
        }

        if (state.sourceStatus == ImportStepStatus.error) {
          return Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text(state.errorMessage ?? '加载失败'),
                const SizedBox(height: 16),
                ElevatedButton(
                  onPressed: () {
                    context.read<ImportBloc>().add(ImportLoadSources());
                  },
                  child: const Text('重试'),
                ),
              ],
            ),
          );
        }

        return Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                '选择导入来源',
                style: Theme.of(context).textTheme.headlineSmall,
              ),
              const SizedBox(height: 8),
              Text(
                '请选择要导入的账单类型',
                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                  color: Theme.of(context).colorScheme.onSurfaceVariant,
                ),
              ),
              const SizedBox(height: 24),
              Expanded(
                child: ListView.separated(
                  itemCount: state.availableSources?.length ?? 0,
                  separatorBuilder: (context, index) => const SizedBox(height: 12),
                  itemBuilder: (context, index) {
                    final source = state.availableSources![index];
                    return _SourceCard(
                      source: source,
                      onTap: () {
                        final importSource = _parseSource(source.id);
                        context.read<ImportBloc>().add(
                          ImportSourceSelected(importSource),
                        );
                      },
                    );
                  },
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}

ImportSource _parseSource(String id) {
  switch (id) {
    case 'alipay':
      return ImportSource.alipay;
    case 'wechat':
      return ImportSource.wechat;
    case 'bank':
      return ImportSource.bank;
    case 'generic':
    default:
      return ImportSource.generic;
  }
}

class _SourceCard extends StatelessWidget {
  final ImportSourceInfo source;
  final VoidCallback onTap;

  const _SourceCard({required this.source, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return Card(
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Row(
            children: [
              _buildIcon(context),
              const SizedBox(width: 16),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      source.name,
                      style: Theme.of(context).textTheme.titleMedium,
                    ),
                    const SizedBox(height: 4),
                    Text(
                      source.description,
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        color: Theme.of(context).colorScheme.onSurfaceVariant,
                      ),
                    ),
                  ],
                ),
              ),
              Icon(
                Icons.arrow_forward_ios,
                size: 16,
                color: Theme.of(context).colorScheme.onSurfaceVariant,
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildIcon(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    switch (source.id) {
      case 'alipay':
        return Container(
          width: 48,
          height: 48,
          decoration: BoxDecoration(
            color: const Color(0xFF1677FF),
            borderRadius: BorderRadius.circular(12),
          ),
          child: const Icon(Icons.payment, color: Colors.white),
        );
      case 'wechat':
        return Container(
          width: 48,
          height: 48,
          decoration: BoxDecoration(
            color: const Color(0xFF07C160),
            borderRadius: BorderRadius.circular(12),
          ),
          child: const Icon(Icons.chat_bubble, color: Colors.white),
        );
      case 'bank':
        return Container(
          width: 48,
          height: 48,
          decoration: BoxDecoration(
            color: colorScheme.primary,
            borderRadius: BorderRadius.circular(12),
          ),
          child: const Icon(Icons.account_balance, color: Colors.white),
        );
      case 'generic':
      default:
        return Container(
          width: 48,
          height: 48,
          decoration: BoxDecoration(
            color: colorScheme.secondary,
            borderRadius: BorderRadius.circular(12),
          ),
          child: const Icon(Icons.file_upload, color: Colors.white),
        );
    }
  }
}
