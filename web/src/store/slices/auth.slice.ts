import { StateCreator } from 'zustand';
import type { User, LoginInput, RegisterInput, AuthResponse } from '../../data/models/user.model';
import { authApi, apiClient } from '../../data/api';

export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

export interface AuthActions {
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  checkAuth: () => Promise<boolean>;
  clearError: () => void;
}

export interface AuthSlice extends AuthState {
  auth: AuthState;
  actions: AuthActions;
}

const initialAuthState: AuthState = {
  user: null,
  isAuthenticated: false,
  isLoading: false,
  error: null,
};

export const createAuthSlice: StateCreator<AuthSlice> = (set) => ({
  auth: initialAuthState,
  user: null,
  isAuthenticated: false,
  isLoading: false,
  error: null,

  actions: {
    login: async (email: string, password: string) => {
      console.log('[AuthSlice] Login called with:', email);
      set((state) => ({
        auth: { ...state.auth, isLoading: true, error: null },
        isLoading: true,
        error: null,
      }));

      try {
        const input: LoginInput = { email, password };
        console.log('[AuthSlice] Calling authApi.login...');
        const response: AuthResponse = await authApi.login(input);
        console.log('[AuthSlice] Login response received, access_token:', response.access_token ? 'yes' : 'no');
        console.log('[AuthSlice] user:', response.user ? 'yes' : 'no');

        // Save token to localStorage and set in API client
        apiClient.setAuthToken(response.access_token);
        console.log('[AuthSlice] Token saved to localStorage');

        // Verify token was saved
        const verifyToken = localStorage.getItem('auth_token');
        console.log('[AuthSlice] Verified token in localStorage:', verifyToken ? 'yes' : 'no');

        // Transform user date - map snake_case to camelCase for frontend
        const user = {
          id: response.user.id,
          email: response.user.email,
          createdAt: response.user.created_at,
        };

        set((state) => ({
          auth: {
            ...state.auth,
            user,
            isAuthenticated: true,
            isLoading: false,
          },
          user,
          isAuthenticated: true,
          isLoading: false,
        }));
      } catch (error) {
        const errorMessage = (error as Error).message;
        set((state) => ({
          auth: {
            ...state.auth,
            isLoading: false,
            error: errorMessage,
          },
          isLoading: false,
          error: errorMessage,
        }));
        throw error;
      }
    },

    register: async (email: string, password: string) => {
      set((state) => ({
        auth: { ...state.auth, isLoading: true, error: null },
        isLoading: true,
        error: null,
      }));

      try {
        const input: RegisterInput = { email, password };
        const response: AuthResponse = await authApi.register(input);

        // Save token to localStorage and set in API client
        apiClient.setAuthToken(response.access_token);

        // Transform user date - map snake_case to camelCase for frontend
        const user = {
          id: response.user.id,
          email: response.user.email,
          createdAt: response.user.created_at,
        };

        set((state) => ({
          auth: {
            ...state.auth,
            user,
            isAuthenticated: true,
            isLoading: false,
          },
          user,
          isAuthenticated: true,
          isLoading: false,
        }));
      } catch (error) {
        const errorMessage = (error as Error).message;
        set((state) => ({
          auth: {
            ...state.auth,
            isLoading: false,
            error: errorMessage,
          },
          isLoading: false,
          error: errorMessage,
        }));
        throw error;
      }
    },

    logout: async () => {
      // Clear token and user data
      apiClient.setAuthToken(null);
      localStorage.removeItem('user');
      set({
        auth: initialAuthState,
        ...initialAuthState,
      });
    },

    checkAuth: async () => {
      const token = localStorage.getItem('auth_token');
      if (token) {
        apiClient.setAuthToken(token);
        try {
          // Try to get current user to verify token is valid
          const userData = await authApi.getMe();
          const user: User = {
            id: userData.user_id,
            email: userData.email,
            createdAt: new Date(),
          };
          set((state) => ({
            auth: {
              ...state.auth,
              user,
              isAuthenticated: true,
              isLoading: false,
            },
            user,
            isAuthenticated: true,
            isLoading: false,
          }));
          return true;
        } catch (error) {
          // Token invalid, clear it
          localStorage.removeItem('auth_token');
          apiClient.setAuthToken(null);
        }
      }
      return false;
    },

    clearError: () => {
      set((state) => ({
        auth: { ...state.auth, error: null },
        error: null,
      }));
    },
  },
});
