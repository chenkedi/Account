import 'package:drift/drift.dart';

import '../app_database.dart';
import '../tables/account_table.dart';

part 'account_dao.g.dart';

@DriftAccessor(tables: [Accounts])
class AccountDao extends DatabaseAccessor<AppDatabase> with _$AccountDaoMixin {
  AccountDao(AppDatabase db) : super(db);

  Future<List<Account>> getAllAccounts() {
    return (select(accounts)
          ..where((a) => a.isDeleted.equals(false))
          ..orderBy([(a) => OrderingTerm(expression: a.name)]))
        .get();
  }

  Stream<List<Account>> watchAllAccounts() {
    return (select(accounts)
          ..where((a) => a.isDeleted.equals(false))
          ..orderBy([(a) => OrderingTerm(expression: a.name)]))
        .watch();
  }

  Future<Account?> getAccountById(String id) {
    return (select(accounts)..where((a) => a.id.equals(id))).getSingleOrNull();
  }

  Future<void> insertAccount(AccountsCompanion account) {
    return into(accounts).insert(account, mode: InsertMode.insertOrReplace);
  }

  Future<void> insertAccounts(List<AccountsCompanion> accountList) {
    return batch((batch) {
      batch.insertAll(accounts, accountList, mode: InsertMode.insertOrReplace);
    });
  }

  Future<void> updateAccount(AccountsCompanion account) {
    return (update(accounts)..where((a) => a.id.equals(account.id.value))).write(account);
  }

  Future<void> deleteAccount(String id) {
    return (update(accounts)..where((a) => a.id.equals(id)))
        .write(AccountsCompanion(isDeleted: const Value(true)));
  }

  Future<List<Account>> getModifiedSince(DateTime since) {
    return (select(accounts)
          ..where((a) => a.lastModifiedAt.isBiggerThanValue(since)))
        .get();
  }
}
