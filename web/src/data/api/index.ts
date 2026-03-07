import { apiClient } from './client';
import { AuthApi } from './auth.api';
import { AccountApi } from './account.api';
import { TransactionApi } from './transaction.api';
import { CategoryApi } from './category.api';
import { SyncApi } from './sync.api';
import { ImportApi } from './import.api';

// Initialize all API instances with the shared apiClient
export const authApi = new AuthApi(apiClient);
export const accountApi = new AccountApi(apiClient);
export const transactionApi = new TransactionApi(apiClient);
export const categoryApi = new CategoryApi(apiClient);
export const syncApi = new SyncApi(apiClient);
export const importApi = new ImportApi(apiClient);

// Export everything
export * from './client';
export * from './auth.api';
export * from './account.api';
export * from './transaction.api';
export * from './category.api';
export * from './sync.api';
export * from './import.api';
