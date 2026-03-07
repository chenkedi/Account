import { create } from 'zustand';
import { devtools, persist, subscribeWithSelector } from 'zustand/middleware';
import { immer } from 'zustand/middleware/immer';
import { AuthSlice, createAuthSlice } from './slices/auth.slice';
import { TransactionSlice, createTransactionSlice } from './slices/transaction.slice';
import { AccountSlice, createAccountSlice } from './slices/account.slice';
import { CategorySlice, createCategorySlice } from './slices/category.slice';
import { StatsSlice, createStatsSlice } from './slices/stats.slice';
import { ImportSlice, createImportSlice } from './slices/import.slice';

type AppStore = AuthSlice &
  TransactionSlice &
  AccountSlice &
  CategorySlice &
  StatsSlice &
  ImportSlice;

export const useAppStore = create<AppStore>()(
  subscribeWithSelector(
    persist(
      immer(
        devtools((...a) => ({
          ...createAuthSlice(...a),
          ...createTransactionSlice(...a),
          ...createAccountSlice(...a),
          ...createCategorySlice(...a),
          ...createStatsSlice(...a),
          ...createImportSlice(...a),
        }))
      ),
      {
        name: 'account-app-storage',
        partialize: (state) => ({
          auth: {
            user: state.auth.user,
            isAuthenticated: state.auth.isAuthenticated,
          },
        }),
      }
    )
  )
);

// 导出便捷 hooks
export const useAuth = () => useAppStore((state) => state.auth);
export const useAuthActions = () => useAppStore((state) => state.actions);

export const useTransactions = () => useAppStore((state) => state.transactionState.transactions);
export const useTransactionsState = () => useAppStore((state) => state.transactionState);
export const useTransactionActions = () => useAppStore((state) => state.transactionActions);

export const useAccounts = () => useAppStore((state) => state.accountState.accounts);
export const useAccountsState = () => useAppStore((state) => state.accountState);
export const useAccountActions = () => useAppStore((state) => state.accountActions);

export const useCategories = () => useAppStore((state) => state.categoryState.categories);
export const useCategoriesState = () => useAppStore((state) => state.categoryState);
export const useCategoryActions = () => useAppStore((state) => state.categoryActions);

export const useStats = () => useAppStore((state) => state.statsState.stats);
export const useStatsState = () => useAppStore((state) => state.statsState);
export const useStatsActions = () => useAppStore((state) => state.statsActions);

export const useImportState = () => useAppStore((state) => state.importState);
export const useImportActions = () => useAppStore((state) => state.importActions);
