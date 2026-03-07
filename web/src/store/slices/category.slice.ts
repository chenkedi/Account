import { StateCreator } from 'zustand';
import type { Category } from '../../data/models/category.model';
import { categoryApi } from '../../data/api';

export interface CategoryState {
  categories: Category[];
  isLoading: boolean;
  error: string | null;
}

export interface CategoryActions {
  fetchCategories: () => Promise<void>;
  fetchCategoriesByType: (type: 'income' | 'expense') => Promise<Category[]>;
  clearError: () => void;
}

export interface CategorySlice {
  categoryState: CategoryState;
  categoryActions: CategoryActions;
}

const initialCategoryState: CategoryState = {
  categories: [],
  isLoading: false,
  error: null,
};

export const createCategorySlice: StateCreator<CategorySlice> = (set) => ({
  categoryState: initialCategoryState,

  categoryActions: {
    fetchCategories: async () => {
      set((state) => ({
        categoryState: { ...state.categoryState, isLoading: true, error: null },
      }));

      try {
        const categories = await categoryApi.getAll();
        set((state) => ({
          categoryState: { ...state.categoryState, categories, isLoading: false },
        }));
      } catch (error) {
        const errorMessage = (error as Error).message;
        set((state) => ({
          categoryState: { ...state.categoryState, isLoading: false, error: errorMessage },
        }));
        throw error;
      }
    },

    fetchCategoriesByType: async (type) => {
      set((state) => ({
        categoryState: { ...state.categoryState, isLoading: true, error: null },
      }));

      try {
        const categories = await categoryApi.getByType(type);
        set((state) => ({
          categoryState: { ...state.categoryState, categories, isLoading: false },
        }));
        return categories;
      } catch (error) {
        const errorMessage = (error as Error).message;
        set((state) => ({
          categoryState: { ...state.categoryState, isLoading: false, error: errorMessage },
        }));
        throw error;
      }
    },

    clearError: () => {
      set((state) => ({
        categoryState: { ...state.categoryState, error: null },
      }));
    },
  },
});
