import axios, { AxiosInstance, AxiosError, InternalAxiosRequestConfig } from 'axios';
import { API_CONSTANTS } from '../../core/constants/api.constants';
import { authEvents, AUTH_EVENTS } from '../../core/events/auth.events';

interface ApiClientConfig {
  baseURL?: string;
  timeout?: number;
}

export class ApiClient {
  private client: AxiosInstance;

  constructor(config: ApiClientConfig = {}) {
    this.client = axios.create({
      baseURL: config.baseURL || API_CONSTANTS.baseUrl,
      timeout: config.timeout || 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.setupInterceptors();
  }

  setAuthToken(token: string | null): void {
    if (token) {
      localStorage.setItem('auth_token', token);
    } else {
      localStorage.removeItem('auth_token');
    }
  }

  private setupInterceptors(): void {
    // Request interceptor
    this.client.interceptors.request.use(
      (config: InternalAxiosRequestConfig) => {
        // Always get the latest token from localStorage
        const token = localStorage.getItem('auth_token');
        console.log('[ApiClient] Request to:', config.url, 'Token present:', !!token);

        // Check if we have an authenticated route but no token
        const isAuthRoute = config.url?.includes('/auth/');
        if (!token && !isAuthRoute) {
          console.log('[ApiClient] No token found for authenticated route, notifying session expired');
          // Emit event immediately for faster response, but still let the request go through
          // to get a 401 from server as a fallback
          authEvents.emit(AUTH_EVENTS.SESSION_EXPIRED);
        }

        if (token && config.headers) {
          config.headers.Authorization = `Bearer ${token}`;
          console.log('[ApiClient] Added Authorization header');
        }
        return config;
      },
      (error: AxiosError) => Promise.reject(error)
    );

    // Response interceptor
    this.client.interceptors.response.use(
      (response) => response.data,
      (error: AxiosError) => {
        console.log('[ApiClient] Response error:', error.response?.status, error.config?.url);
        if (error.response?.status === 401) {
          // Handle unauthorized - clear token and notify app
          console.log('[ApiClient] 401 received, clearing token and notifying session expired');
          localStorage.removeItem('auth_token');
          authEvents.emit(AUTH_EVENTS.SESSION_EXPIRED);
        }
        return Promise.reject(error);
      }
    );
  }

  get instance(): AxiosInstance {
    return this.client;
  }

  async get<T>(url: string, config?: Record<string, unknown>): Promise<T> {
    return this.client.get(url, config);
  }

  async post<T>(url: string, data?: unknown, config?: Record<string, unknown>): Promise<T> {
    return this.client.post(url, data, config);
  }

  async put<T>(url: string, data?: unknown, config?: Record<string, unknown>): Promise<T> {
    return this.client.put(url, data, config);
  }

  async delete<T>(url: string, config?: Record<string, unknown>): Promise<T> {
    return this.client.delete(url, config);
  }

  async upload<T>(url: string, file: File, data?: Record<string, unknown>): Promise<T> {
    const formData = new FormData();
    formData.append('file', file);
    if (data) {
      Object.entries(data).forEach(([key, value]) => {
        formData.append(key, String(value));
      });
    }
    return this.client.post(url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  }
}

export const apiClient = new ApiClient();
