import Dexie from 'dexie';
import type { Account, Category, Transaction } from '../models';

export class AppDatabase extends Dexie {
  accounts!: Dexie.Table<Account, string>;
  categories!: Dexie.Table<Category, string>;
  transactions!: Dexie.Table<Transaction, string>;
  syncState!: Dexie.Table<{ deviceId: string; lastSyncAt: string }, string>;

  constructor() {
    super('AccountDB');

    this.version(1).stores({
      accounts: '&id, userId, lastModifiedAt, isDeleted',
      categories: '&id, userId, type, lastModifiedAt, isDeleted',
      transactions: '&id, userId, accountId, categoryId, transactionDate, lastModifiedAt, isDeleted',
      syncState: '&deviceId',
    });
  }
}
