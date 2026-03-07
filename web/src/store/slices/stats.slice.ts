import { StateCreator } from 'zustand';
import { transactionApi } from '../../data/api';

export type TimeRange = 'month' | 'quarter' | 'year';

export interface StatsData {
  totalIncome: number;
  totalExpense: number;
  netIncome: number;
  byCategory: Array<{ categoryId: string; categoryName: string; amount: number; type: string }>;
  monthlyTrend: Array<{ month: string; income: number; expense: number }>;
}

export interface StatsState {
  stats: StatsData | null;
  timeRange: TimeRange;
  isLoading: boolean;
  error: string | null;
}

export interface StatsActions {
  fetchStats: (startDate: Date, endDate: Date) => Promise<void>;
  setTimeRange: (range: TimeRange) => void;
  clearError: () => void;
}

export interface StatsSlice {
  statsState: StatsState;
  statsActions: StatsActions;
}

const initialStatsState: StatsState = {
  stats: null,
  timeRange: 'month',
  isLoading: false,
  error: null,
};

export const createStatsSlice: StateCreator<StatsSlice> = (set) => ({
  statsState: initialStatsState,

  statsActions: {
    fetchStats: async (startDate, endDate) => {
      set((state) => ({
        statsState: { ...state.statsState, isLoading: true, error: null },
      }));

      try {
        const stats = await transactionApi.getStats(startDate, endDate);
        set((state) => ({
          statsState: { ...state.statsState, stats, isLoading: false },
        }));
      } catch (error) {
        const errorMessage = (error as Error).message;
        set((state) => ({
          statsState: { ...state.statsState, isLoading: false, error: errorMessage },
        }));
        throw error;
      }
    },

    setTimeRange: (range) => {
      set((state) => ({
        statsState: { ...state.statsState, timeRange: range },
      }));
    },

    clearError: () => {
      set((state) => ({
        statsState: { ...state.statsState, error: null },
      }));
    },
  },
});
