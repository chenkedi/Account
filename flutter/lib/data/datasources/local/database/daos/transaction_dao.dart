import 'package:drift/drift.dart';

import '../app_database.dart';
import '../tables/transaction_table.dart';

part 'transaction_dao.g.dart';

@DriftAccessor(tables: [Transactions])
class TransactionDao extends DatabaseAccessor<AppDatabase> with _$TransactionDaoMixin {
  TransactionDao(AppDatabase db) : super(db);

  Future<List<Transaction>> getAllTransactions({int limit = 100, int offset = 0}) {
    return (select(transactions)
          ..where((t) => t.isDeleted.equals(false))
          ..orderBy([(t) => OrderingTerm(expression: t.transactionDate, mode: OrderingMode.desc)])
          ..limit(limit, offset: offset))
        .get();
  }

  Stream<List<Transaction>> watchAllTransactions({int limit = 100}) {
    return (select(transactions)
          ..where((t) => t.isDeleted.equals(false))
          ..orderBy([(t) => OrderingTerm(expression: t.transactionDate, mode: OrderingMode.desc)])
          ..limit(limit))
        .watch();
  }

  Future<List<Transaction>> getTransactionsByDateRange(DateTime start, DateTime end) {
    return (select(transactions)
          ..where((t) => t.isDeleted.equals(false) & t.transactionDate.isBetweenValues(start, end))
          ..orderBy([(t) => OrderingTerm(expression: t.transactionDate, mode: OrderingMode.desc)]))
        .get();
  }

  Future<List<Transaction>> getTransactionsByAccountId(String accountId, {int limit = 50}) {
    return (select(transactions)
          ..where((t) => t.isDeleted.equals(false) & t.accountId.equals(accountId))
          ..orderBy([(t) => OrderingTerm(expression: t.transactionDate, mode: OrderingMode.desc)])
          ..limit(limit))
        .get();
  }

  Future<Transaction?> getTransactionById(String id) {
    return (select(transactions)..where((t) => t.id.equals(id))).getSingleOrNull();
  }

  Future<void> insertTransaction(TransactionsCompanion transaction) {
    return into(transactions).insert(transaction, mode: InsertMode.insertOrReplace);
  }

  Future<void> insertTransactions(List<TransactionsCompanion> transactionList) {
    return batch((batch) {
      batch.insertAll(transactions, transactionList, mode: InsertMode.insertOrReplace);
    });
  }

  Future<void> updateTransaction(TransactionsCompanion transaction) {
    return (update(transactions)..where((t) => t.id.equals(transaction.id.value))).write(transaction);
  }

  Future<void> deleteTransaction(String id) {
    return (update(transactions)..where((t) => t.id.equals(id)))
        .write(TransactionsCompanion(isDeleted: const Value(true)));
  }

  Future<List<Transaction>> getModifiedSince(DateTime since) {
    return (select(transactions)
          ..where((t) => t.lastModifiedAt.isBiggerThanValue(since)))
        .get();
  }

  Future<double> getTotalIncomeForPeriod(DateTime start, DateTime end) async {
    final result = await (select(transactions)
          ..where((t) =>
              t.isDeleted.equals(false) &
              t.type.equals('income') &
              t.transactionDate.isBetweenValues(start, end)))
        .get();
    double total = 0;
    for (final t in result) {
      total += t.amount;
    }
    return total;
  }

  Future<double> getTotalExpenseForPeriod(DateTime start, DateTime end) async {
    final result = await (select(transactions)
          ..where((t) =>
              t.isDeleted.equals(false) &
              t.type.equals('expense') &
              t.transactionDate.isBetweenValues(start, end)))
        .get();
    double total = 0;
    for (final t in result) {
      total += t.amount;
    }
    return total;
  }
}
