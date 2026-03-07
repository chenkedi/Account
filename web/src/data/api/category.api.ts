import { ApiClient } from './client';
import type { Category, CategoryCreateInput, CategoryUpdateInput, CategoryType } from '../models/category.model';
import { API_CONSTANTS } from '../../core/constants/api.constants';

export class CategoryApi {
  constructor(private apiClient: ApiClient) {}

  async getAll(): Promise<Category[]> {
    return this.apiClient.get<Category[]>(API_CONSTANTS.categories);
  }

  async getByType(type: CategoryType): Promise<Category[]> {
    return this.apiClient.get<Category[]>(`${API_CONSTANTS.categories}/type/${type}`);
  }

  async getById(id: string): Promise<Category> {
    return this.apiClient.get<Category>(`${API_CONSTANTS.categories}/${id}`);
  }

  async create(data: CategoryCreateInput): Promise<Category> {
    return this.apiClient.post<Category>(API_CONSTANTS.categories, data);
  }

  async update(id: string, data: CategoryUpdateInput): Promise<Category> {
    return this.apiClient.put<Category>(`${API_CONSTANTS.categories}/${id}`, data);
  }

  async delete(id: string): Promise<void> {
    return this.apiClient.delete(`${API_CONSTANTS.categories}/${id}`);
  }
}
