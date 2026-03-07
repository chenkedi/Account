import { StateCreator } from 'zustand';
import type { Account, AccountCreateInput, AccountUpdateInput } from '../../data/models/account.model';
import { accountApi } from '../../data/api';

export interface AccountState {
  accounts: Account[];
  isLoading: boolean;
  error: string | null;
}

export interface AccountActions {
  fetchAccounts: () => Promise<void>;
  createAccount: (data: AccountCreateInput) => Promise<Account>;
  updateAccount: (id: string, data: AccountUpdateInput) => Promise<Account>;
  deleteAccount: (id: string) => Promise<void>;
  clearError: () => void;
}

export interface AccountSlice {
  accountState: AccountState;
  accountActions: AccountActions;
}

const initialAccountState: AccountState = {
  accounts: [],
  isLoading: false,
  error: null,
};

export const createAccountSlice: StateCreator<AccountSlice> = (set) => ({
  accountState: initialAccountState,

  accountActions: {
    fetchAccounts: async () => {
      set((state) => ({
        accountState: { ...state.accountState, isLoading: true, error: null },
      }));

      try {
        const accounts = await accountApi.getAll();
        set((state) => ({
          accountState: { ...state.accountState, accounts, isLoading: false },
        }));
      } catch (error) {
        const errorMessage = (error as Error).message;
        set((state) => ({
          accountState: { ...state.accountState, isLoading: false, error: errorMessage },
        }));
        throw error;
      }
    },

    createAccount: async (data) => {
      set((state) => ({
        accountState: { ...state.accountState, isLoading: true, error: null },
      }));

      try {
        const account = await accountApi.create(data);
        set((state) => ({
          accountState: {
            ...state.accountState,
            accounts: [...state.accountState.accounts, account],
            isLoading: false,
          },
        }));
        return account;
      } catch (error) {
        const errorMessage = (error as Error).message;
        set((state) => ({
          accountState: { ...state.accountState, isLoading: false, error: errorMessage },
        }));
        throw error;
      }
    },

    updateAccount: async (id, data) => {
      set((state) => ({
        accountState: { ...state.accountState, isLoading: true, error: null },
      }));

      try {
        const updated = await accountApi.update(id, data);
        set((state) => ({
          accountState: {
            ...state.accountState,
            accounts: state.accountState.accounts.map((a) => (a.id === id ? updated : a)),
            isLoading: false,
          },
        }));
        return updated;
      } catch (error) {
        const errorMessage = (error as Error).message;
        set((state) => ({
          accountState: { ...state.accountState, isLoading: false, error: errorMessage },
        }));
        throw error;
      }
    },

    deleteAccount: async (id) => {
      set((state) => ({
        accountState: { ...state.accountState, isLoading: true, error: null },
      }));

      try {
        await accountApi.delete(id);
        set((state) => ({
          accountState: {
            ...state.accountState,
            accounts: state.accountState.accounts.filter((a) => a.id !== id),
            isLoading: false,
          },
        }));
      } catch (error) {
        const errorMessage = (error as Error).message;
        set((state) => ({
          accountState: { ...state.accountState, isLoading: false, error: errorMessage },
        }));
        throw error;
      }
    },

    clearError: () => {
      set((state) => ({
        accountState: { ...state.accountState, error: null },
      }));
    },
  },
});
