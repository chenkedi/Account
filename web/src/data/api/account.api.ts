import { ApiClient } from './client';
import type { Account, AccountCreateInput, AccountUpdateInput } from '../models/account.model';
import { API_CONSTANTS } from '../../core/constants/api.constants';

export class AccountApi {
  constructor(private apiClient: ApiClient) {}

  async getAll(): Promise<Account[]> {
    return this.apiClient.get<Account[]>(API_CONSTANTS.accounts);
  }

  async getById(id: string): Promise<Account> {
    return this.apiClient.get<Account>(`${API_CONSTANTS.accounts}/${id}`);
  }

  async create(data: AccountCreateInput): Promise<Account> {
    return this.apiClient.post<Account>(API_CONSTANTS.accounts, data);
  }

  async update(id: string, data: AccountUpdateInput): Promise<Account> {
    return this.apiClient.put<Account>(`${API_CONSTANTS.accounts}/${id}`, data);
  }

  async delete(id: string): Promise<void> {
    return this.apiClient.delete(`${API_CONSTANTS.accounts}/${id}`);
  }
}
