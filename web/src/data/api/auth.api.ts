import { ApiClient } from './client';
import type { LoginInput, RegisterInput, AuthResponse } from '../models/user.model';
import { API_CONSTANTS } from '../../core/constants/api.constants';

export class AuthApi {
  constructor(private apiClient: ApiClient) {}

  async login(data: LoginInput): Promise<AuthResponse> {
    return this.apiClient.post<AuthResponse>(API_CONSTANTS.auth.login, data);
  }

  async register(data: RegisterInput): Promise<AuthResponse> {
    return this.apiClient.post<AuthResponse>(API_CONSTANTS.auth.register, data);
  }

  async getMe(): Promise<{ user_id: string; email: string }> {
    return this.apiClient.get<{ user_id: string; email: string }>('/me');
  }
}
