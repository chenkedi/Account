import { StateCreator } from 'zustand';
import type { ImportSource, ImportPreview, ImportResult } from '../../data/models/import.model';
import { importApi } from '../../data/api';

export interface ImportState {
  importSource: ImportSource | null;
  preview: ImportPreview | null;
  result: ImportResult | null;
  isLoading: boolean;
  error: string | null;
}

export interface ImportActions {
  selectSource: (source: ImportSource) => void;
  uploadFile: (source: ImportSource, file: File) => Promise<ImportPreview>;
  confirmImport: (
    transactions: Array<{
      rawIndex: number;
      accountId: string | null;
      categoryId: string | null;
      note: string | null;
    }>
  ) => Promise<ImportResult>;
  resetImport: () => void;
  clearError: () => void;
}

export interface ImportSlice {
  importState: ImportState;
  importActions: ImportActions;
}

const initialImportState: ImportState = {
  importSource: null,
  preview: null,
  result: null,
  isLoading: false,
  error: null,
};

export const createImportSlice: StateCreator<ImportSlice> = (set, get) => ({
  importState: initialImportState,

  importActions: {
    selectSource: (source) => {
      set((state) => ({
        importState: { ...state.importState, importSource: source },
      }));
    },

    uploadFile: async (source, file) => {
      set((state) => ({
        importState: { ...state.importState, isLoading: true, error: null },
      }));

      try {
        const preview = await importApi.uploadAndParse(source, file);
        set((state) => ({
          importState: { ...state.importState, preview, isLoading: false },
        }));
        return preview;
      } catch (error) {
        const errorMessage = (error as Error).message;
        set((state) => ({
          importState: { ...state.importState, isLoading: false, error: errorMessage },
        }));
        throw error;
      }
    },

    confirmImport: async (transactions) => {
      set((state) => ({
        importState: { ...state.importState, isLoading: true, error: null },
      }));

      try {
        const preview = get().importState.preview;
        if (!preview) {
          throw new Error('No import preview available');
        }

        // Convert to snake_case for API
        const apiTransactions = transactions.map((t) => ({
          raw_index: t.rawIndex,
          selected_account_id: t.accountId,
          selected_category_id: t.categoryId,
          note: t.note,
        }));

        const result = await importApi.executeImport(preview.jobId, apiTransactions);
        set((state) => ({
          importState: { ...state.importState, result, isLoading: false },
        }));
        return result;
      } catch (error) {
        const errorMessage = (error as Error).message;
        set((state) => ({
          importState: { ...state.importState, isLoading: false, error: errorMessage },
        }));
        throw error;
      }
    },

    resetImport: () => {
      set({
        importState: initialImportState,
      });
    },

    clearError: () => {
      set((state) => ({
        importState: { ...state.importState, error: null },
      }));
    },
  },
});
