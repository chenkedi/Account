import 'package:drift/drift.dart' hide isNotNull;
import 'package:flutter_test/flutter_test.dart';
import 'package:account/data/datasources/local/database/app_database.dart';
import '../../helpers/test_helpers.dart';

void main() {
  group('TransactionDao', () {
    late AppDatabase db;

    setUp(() {
      db = createTestDatabase();
    });

    // No need for tearDown - createTestDatabase registers it automatically

    test('insertTransaction and getTransactionById work correctly', () async {
      const id = 'test-transaction-1';
      const userId = 'user-1';
      const accountId = 'account-1';
      final now = DateTime.now().toUtc();

      final companion = TransactionsCompanion(
        id: const Value(id),
        userId: const Value(userId),
        accountId: const Value(accountId),
        type: const Value('expense'),
        amount: const Value(100.0),
        transactionDate: Value(now),
        createdAt: Value(now),
        updatedAt: Value(now),
        lastModifiedAt: Value(now),
      );

      await db.transactionDao.insertTransaction(companion);
      final transaction = await db.transactionDao.getTransactionById(id);

      expect(transaction, isNotNull);
      expect(transaction?.id, id);
      expect(transaction?.userId, userId);
      expect(transaction?.accountId, accountId);
      expect(transaction?.amount, 100.0);
    });

    test('getAllTransactions returns non-deleted transactions ordered by date', () async {
      final now = DateTime.now().toUtc();

      await db.transactionDao.insertTransaction(TransactionsCompanion(
        id: const Value('t1'),
        userId: const Value('user1'),
        accountId: const Value('acc1'),
        type: const Value('expense'),
        amount: const Value(100),
        transactionDate: Value(now),
        createdAt: Value(now),
        updatedAt: Value(now),
        lastModifiedAt: Value(now),
      ));

      await db.transactionDao.insertTransaction(TransactionsCompanion(
        id: const Value('t2'),
        userId: const Value('user1'),
        accountId: const Value('acc1'),
        type: const Value('income'),
        amount: const Value(200),
        transactionDate: Value(now.add(const Duration(days: 1))),
        createdAt: Value(now),
        updatedAt: Value(now),
        lastModifiedAt: Value(now),
      ));

      await db.transactionDao.insertTransaction(TransactionsCompanion(
        id: const Value('t3'),
        userId: const Value('user1'),
        accountId: const Value('acc1'),
        type: const Value('expense'),
        amount: const Value(300),
        transactionDate: Value(now.subtract(const Duration(days: 1))),
        createdAt: Value(now),
        updatedAt: Value(now),
        lastModifiedAt: Value(now),
        isDeleted: const Value(true),
      ));

      final transactions = await db.transactionDao.getAllTransactions();

      expect(transactions.length, 2);
      expect(transactions[0].id, 't2'); // newest first
      expect(transactions[1].id, 't1');
    });

    test('watchAllTransactions emits updates', () async {
      final now = DateTime.now().toUtc();

      expectLater(
        db.transactionDao.watchAllTransactions(),
        emits(isEmpty),
      );

      await db.transactionDao.insertTransaction(TransactionsCompanion(
        id: const Value('t1'),
        userId: const Value('user1'),
        accountId: const Value('acc1'),
        type: const Value('expense'),
        amount: const Value(100),
        transactionDate: Value(now),
        createdAt: Value(now),
        updatedAt: Value(now),
        lastModifiedAt: Value(now),
      ));
    });

    test('updateTransaction updates existing transaction', () async {
      final now = DateTime.now().toUtc();
      const id = 'update-test';

      await db.transactionDao.insertTransaction(TransactionsCompanion(
        id: const Value(id),
        userId: const Value('user1'),
        accountId: const Value('acc1'),
        type: const Value('expense'),
        amount: const Value(100),
        transactionDate: Value(now),
        createdAt: Value(now),
        updatedAt: Value(now),
        lastModifiedAt: Value(now),
      ));

      await db.transactionDao.updateTransaction(const TransactionsCompanion(
        id: Value(id),
        amount: Value(200),
        note: Value('Updated note'),
      ));

      final updated = await db.transactionDao.getTransactionById(id);
      expect(updated!.amount, 200);
      expect(updated.note, 'Updated note');
    });

    test('deleteTransaction soft deletes transaction', () async {
      final now = DateTime.now().toUtc();
      const id = 'delete-test';

      await db.transactionDao.insertTransaction(TransactionsCompanion(
        id: const Value(id),
        userId: const Value('user1'),
        accountId: const Value('acc1'),
        type: const Value('expense'),
        amount: const Value(100),
        transactionDate: Value(now),
        createdAt: Value(now),
        updatedAt: Value(now),
        lastModifiedAt: Value(now),
      ));

      await db.transactionDao.deleteTransaction(id);
      final transaction = await db.transactionDao.getTransactionById(id);
      final all = await db.transactionDao.getAllTransactions();

      expect(transaction?.isDeleted, true);
      expect(all, isEmpty);
    });

    test('getModifiedSince returns transactions modified after date', () async {
      final base = DateTime(2024, 1, 1).toUtc();

      await db.transactionDao.insertTransaction(TransactionsCompanion(
        id: const Value('old'),
        userId: const Value('user1'),
        accountId: const Value('acc1'),
        type: const Value('expense'),
        amount: const Value(100),
        transactionDate: Value(base),
        createdAt: Value(base),
        updatedAt: Value(base),
        lastModifiedAt: Value(base),
      ));

      await db.transactionDao.insertTransaction(TransactionsCompanion(
        id: const Value('new'),
        userId: const Value('user1'),
        accountId: const Value('acc1'),
        type: const Value('expense'),
        amount: const Value(200),
        transactionDate: Value(base),
        createdAt: Value(base),
        updatedAt: Value(base),
        lastModifiedAt: Value(base.add(const Duration(days: 10))),
      ));

      final modified = await db.transactionDao.getModifiedSince(
        base.add(const Duration(days: 5)),
      );

      expect(modified.length, 1);
      expect(modified.first.id, 'new');
    });

    test('getTotalIncomeForPeriod sums income transactions', () async {
      final start = DateTime(2024, 1, 1).toUtc();
      final end = DateTime(2024, 1, 31).toUtc();

      await db.transactionDao.insertTransaction(TransactionsCompanion(
        id: const Value('income1'),
        userId: const Value('user1'),
        accountId: const Value('acc1'),
        type: const Value('income'),
        amount: const Value(1000),
        transactionDate: Value(start.add(const Duration(days: 5))),
        createdAt: Value(start),
        updatedAt: Value(start),
        lastModifiedAt: Value(start),
      ));

      await db.transactionDao.insertTransaction(TransactionsCompanion(
        id: const Value('income2'),
        userId: const Value('user1'),
        accountId: const Value('acc1'),
        type: const Value('income'),
        amount: const Value(500),
        transactionDate: Value(start.add(const Duration(days: 10))),
        createdAt: Value(start),
        updatedAt: Value(start),
        lastModifiedAt: Value(start),
      ));

      await db.transactionDao.insertTransaction(TransactionsCompanion(
        id: const Value('expense1'),
        userId: const Value('user1'),
        accountId: const Value('acc1'),
        type: const Value('expense'),
        amount: const Value(300),
        transactionDate: Value(start.add(const Duration(days: 15))),
        createdAt: Value(start),
        updatedAt: Value(start),
        lastModifiedAt: Value(start),
      ));

      final total = await db.transactionDao.getTotalIncomeForPeriod(start, end);
      expect(total, 1500);
    });

    test('getTotalExpenseForPeriod sums expense transactions', () async {
      final start = DateTime(2024, 1, 1).toUtc();
      final end = DateTime(2024, 1, 31).toUtc();

      await db.transactionDao.insertTransaction(TransactionsCompanion(
        id: const Value('expense1'),
        userId: const Value('user1'),
        accountId: const Value('acc1'),
        type: const Value('expense'),
        amount: const Value(100),
        transactionDate: Value(start.add(const Duration(days: 5))),
        createdAt: Value(start),
        updatedAt: Value(start),
        lastModifiedAt: Value(start),
      ));

      await db.transactionDao.insertTransaction(TransactionsCompanion(
        id: const Value('expense2'),
        userId: const Value('user1'),
        accountId: const Value('acc1'),
        type: const Value('expense'),
        amount: const Value(200),
        transactionDate: Value(start.add(const Duration(days: 10))),
        createdAt: Value(start),
        updatedAt: Value(start),
        lastModifiedAt: Value(start),
      ));

      final total = await db.transactionDao.getTotalExpenseForPeriod(start, end);
      expect(total, 300);
    });
  });
}
