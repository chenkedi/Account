import '../datasources/local/database/app_database.dart';
import '../datasources/local/database/daos/account_dao.dart';
import '../../core/network/api_client.dart';

class AccountRepository {
  final AccountDao _accountDao;
  final ApiClient _apiClient;

  AccountRepository(this._accountDao, this._apiClient);

  Stream<List<Account>> watchAllAccounts() => _accountDao.watchAllAccounts();

  Future<List<Account>> getAllAccounts() => _accountDao.getAllAccounts();

  Future<Account?> getAccountById(String id) => _accountDao.getAccountById(id);

  Future<void> addAccount(AccountsCompanion account) => _accountDao.insertAccount(account);

  Future<void> updateAccount(AccountsCompanion account) => _accountDao.updateAccount(account);

  Future<void> deleteAccount(String id) => _accountDao.deleteAccount(id);
}
