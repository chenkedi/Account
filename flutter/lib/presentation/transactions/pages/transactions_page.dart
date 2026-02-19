import 'package:flutter/material.dart';

import '../../import/pages/import_page.dart';

class TransactionsPage extends StatelessWidget {
  const TransactionsPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Transactions'),
        actions: [
          PopupMenuButton(
            itemBuilder: (context) => [
              const PopupMenuItem(
                value: 'export',
                child: ListTile(
                  leading: Icon(Icons.file_download),
                  title: Text('导出 CSV'),
                ),
              ),
            ],
            onSelected: (value) {
              if (value == 'export') {
                _exportTransactions(context);
              }
            },
          ),
        ],
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.swap_horiz,
              size: 64,
              color: Theme.of(context).colorScheme.primary,
            ),
            const SizedBox(height: 16),
            Text(
              'Transactions',
              style: Theme.of(context).textTheme.headlineMedium,
            ),
            const SizedBox(height: 8),
            Text(
              'Your transaction history',
              style: Theme.of(context).textTheme.bodyLarge,
            ),
          ],
        ),
      ),
      floatingActionButton: Column(
        mainAxisAlignment: MainAxisAlignment.end,
        children: [
          FloatingActionButton(
            heroTag: 'import_btn',
            onPressed: () {
              Navigator.of(context).push(
                MaterialPageRoute(builder: (_) => const ImportPage()),
              );
            },
            child: const Icon(Icons.file_upload),
          ),
          const SizedBox(height: 16),
          FloatingActionButton(
            heroTag: 'add_transaction_btn',
            onPressed: () {
              // TODO: Navigate to add transaction page
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('添加交易功能开发中')),
              );
            },
            child: const Icon(Icons.add),
          ),
        ],
      ),
    );
  }

  void _exportTransactions(BuildContext context) {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('导出功能开发中')),
    );
  }
}
