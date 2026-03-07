import { StateCreator } from 'zustand';
import type { Transaction, TransactionCreateInput, TransactionUpdateInput } from '../../data/models/transaction.model';
import { transactionApi } from '../../data/api';

export interface TransactionState {
  transactions: Transaction[];
  isLoading: boolean;
  error: string | null;
}

export interface TransactionActions {
  fetchTransactions: (params?: { limit?: number; offset?: number }) => Promise<void>;
  createTransaction: (data: TransactionCreateInput) => Promise<Transaction>;
  updateTransaction: (id: string, data: TransactionUpdateInput) => Promise<Transaction>;
  deleteTransaction: (id: string) => Promise<void>;
  clearError: () => void;
}

export interface TransactionSlice {
  transactionState: TransactionState;
  transactionActions: TransactionActions;
}

const initialTransactionState: TransactionState = {
  transactions: [],
  isLoading: false,
  error: null,
};

export const createTransactionSlice: StateCreator<TransactionSlice> = (set) => ({
  transactionState: initialTransactionState,

  transactionActions: {
    fetchTransactions: async (params) => {
      set((state) => ({
        transactionState: { ...state.transactionState, isLoading: true, error: null },
      }));

      try {
        const transactions = await transactionApi.getAll(params);
        set((state) => ({
          transactionState: { ...state.transactionState, transactions, isLoading: false },
        }));
      } catch (error) {
        const errorMessage = (error as Error).message;
        set((state) => ({
          transactionState: { ...state.transactionState, isLoading: false, error: errorMessage },
        }));
        throw error;
      }
    },

    createTransaction: async (data) => {
      set((state) => ({
        transactionState: { ...state.transactionState, isLoading: true, error: null },
      }));

      try {
        const transaction = await transactionApi.create(data);
        set((state) => ({
          transactionState: {
            ...state.transactionState,
            transactions: [transaction, ...state.transactionState.transactions],
            isLoading: false,
          },
        }));
        return transaction;
      } catch (error) {
        const errorMessage = (error as Error).message;
        set((state) => ({
          transactionState: { ...state.transactionState, isLoading: false, error: errorMessage },
        }));
        throw error;
      }
    },

    updateTransaction: async (id, data) => {
      set((state) => ({
        transactionState: { ...state.transactionState, isLoading: true, error: null },
      }));

      try {
        const updated = await transactionApi.update(id, data);
        set((state) => ({
          transactionState: {
            ...state.transactionState,
            transactions: state.transactionState.transactions.map((t) => (t.id === id ? updated : t)),
            isLoading: false,
          },
        }));
        return updated;
      } catch (error) {
        const errorMessage = (error as Error).message;
        set((state) => ({
          transactionState: { ...state.transactionState, isLoading: false, error: errorMessage },
        }));
        throw error;
      }
    },

    deleteTransaction: async (id) => {
      set((state) => ({
        transactionState: { ...state.transactionState, isLoading: true, error: null },
      }));

      try {
        await transactionApi.delete(id);
        set((state) => ({
          transactionState: {
            ...state.transactionState,
            transactions: state.transactionState.transactions.filter((t) => t.id !== id),
            isLoading: false,
          },
        }));
      } catch (error) {
        const errorMessage = (error as Error).message;
        set((state) => ({
          transactionState: { ...state.transactionState, isLoading: false, error: errorMessage },
        }));
        throw error;
      }
    },

    clearError: () => {
      set((state) => ({
        transactionState: { ...state.transactionState, error: null },
      }));
    },
  },
});
