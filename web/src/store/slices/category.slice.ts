import { StateCreator } from 'zustand';
import type { Category, CategoryCreateInput, CategoryUpdateInput } from '../../data/models/category.model';
import { categoryApi } from '../../data/api';

export interface CategoryState {
  categories: Category[];
  isLoading: boolean;
  error: string | null;
}

export interface CategoryActions {
  fetchCategories: () => Promise<void>;
  fetchCategoriesByType: (type: 'income' | 'expense') => Promise<Category[]>;
  createCategory: (data: CategoryCreateInput) => Promise<Category>;
  updateCategory: (id: string, data: CategoryUpdateInput) => Promise<Category>;
  deleteCategory: (id: string) => Promise<void>;
  getCategoryById: (id: string) => Category | undefined;
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

export const createCategorySlice: StateCreator<CategorySlice> = (set, get) => ({
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

    createCategory: async (data) => {
      set((state) => ({
        categoryState: { ...state.categoryState, isLoading: true, error: null },
      }));

      try {
        const category = await categoryApi.create(data);
        set((state) => ({
          categoryState: {
            ...state.categoryState,
            categories: [...state.categoryState.categories, category],
            isLoading: false,
          },
        }));
        return category;
      } catch (error) {
        const errorMessage = (error as Error).message;
        set((state) => ({
          categoryState: { ...state.categoryState, isLoading: false, error: errorMessage },
        }));
        throw error;
      }
    },

    updateCategory: async (id, data) => {
      set((state) => ({
        categoryState: { ...state.categoryState, isLoading: true, error: null },
      }));

      try {
        const updated = await categoryApi.update(id, data);
        set((state) => ({
          categoryState: {
            ...state.categoryState,
            categories: state.categoryState.categories.map((c) => (c.id === id ? updated : c)),
            isLoading: false,
          },
        }));
        return updated;
      } catch (error) {
        const errorMessage = (error as Error).message;
        set((state) => ({
          categoryState: { ...state.categoryState, isLoading: false, error: errorMessage },
        }));
        throw error;
      }
    },

    deleteCategory: async (id) => {
      set((state) => ({
        categoryState: { ...state.categoryState, isLoading: true, error: null },
      }));

      try {
        await categoryApi.delete(id);
        set((state) => ({
          categoryState: {
            ...state.categoryState,
            categories: state.categoryState.categories.filter((c) => c.id !== id),
            isLoading: false,
          },
        }));
      } catch (error) {
        const errorMessage = (error as Error).message;
        set((state) => ({
          categoryState: { ...state.categoryState, isLoading: false, error: errorMessage },
        }));
        throw error;
      }
    },

    getCategoryById: (id) => {
      return get().categoryState.categories.find((c) => c.id === id);
    },

    clearError: () => {
      set((state) => ({
        categoryState: { ...state.categoryState, error: null },
      }));
    },
  },
});
