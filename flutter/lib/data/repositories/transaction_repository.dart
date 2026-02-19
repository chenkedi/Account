import '../datasources/local/database/app_database.dart';
import '../datasources/local/database/daos/transaction_dao.dart';
import '../../core/network/api_client.dart';

class TransactionRepository {
  final TransactionDao _transactionDao;
  final ApiClient _apiClient;

  TransactionRepository(this._transactionDao, this._apiClient);

  Stream<List<Transaction>> watchAllTransactions({int limit = 100}) =>
      _transactionDao.watchAllTransactions(limit: limit);

  Future<List<Transaction>> getAllTransactions({int limit = 100, int offset = 0}) =>
      _transactionDao.getAllTransactions(limit: limit, offset: offset);

  Future<List<Transaction>> getTransactionsByDateRange(DateTime start, DateTime end) =>
      _transactionDao.getTransactionsByDateRange(start, end);

  Future<Transaction?> getTransactionById(String id) =>
      _transactionDao.getTransactionById(id);

  Future<void> addTransaction(TransactionsCompanion transaction) =>
      _transactionDao.insertTransaction(transaction);

  Future<void> updateTransaction(TransactionsCompanion transaction) =>
      _transactionDao.updateTransaction(transaction);

  Future<void> deleteTransaction(String id) =>
      _transactionDao.deleteTransaction(id);

  Future<double> getTotalIncomeForPeriod(DateTime start, DateTime end) =>
      _transactionDao.getTotalIncomeForPeriod(start, end);

  Future<double> getTotalExpenseForPeriod(DateTime start, DateTime end) =>
      _transactionDao.getTotalExpenseForPeriod(start, end);
}
