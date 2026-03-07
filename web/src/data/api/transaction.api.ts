import { ApiClient } from './client';
import type { Transaction, TransactionCreateInput, TransactionUpdateInput } from '../models/transaction.model';
import { API_CONSTANTS } from '../../core/constants/api.constants';

// Backend detailed stats response (snake_case)
interface BackendDetailedStatsResponse {
  summary: {
    income_total: number;
    expense_total: number;
    net_total: number;
  };
  by_category: Array<{
    category_id: string;
    category_name: string;
    category_type: string;
    total_amount: number;
    transaction_count: number;
    percentage: number;
  }>;
  monthly_trend: Array<{
    year: number;
    month: number;
    income_total: number;
    expense_total: number;
    net_total: number;
  }>;
}

// Frontend stats response (camelCase)
export interface StatsResponse {
  totalIncome: number;
  totalExpense: number;
  netIncome: number;
  byCategory: Array<{ categoryId: string; categoryName: string; amount: number; type: string }>;
  monthlyTrend: Array<{ month: string; income: number; expense: number }>;
}

// Transform snake_case to camelCase
function transformStatsResponse(data: BackendDetailedStatsResponse): StatsResponse {
  return {
    totalIncome: data.summary.income_total,
    totalExpense: data.summary.expense_total,
    netIncome: data.summary.net_total,
    byCategory: data.by_category?.map((item) => ({
      categoryId: item.category_id,
      categoryName: item.category_name,
      amount: item.total_amount,
      type: item.category_type,
    })) || [],
    monthlyTrend: data.monthly_trend?.map((item) => ({
      month: `${item.year}-${String(item.month).padStart(2, '0')}`,
      income: item.income_total,
      expense: item.expense_total,
    })) || [],
  };
}

export class TransactionApi {
  constructor(private apiClient: ApiClient) {}

  async getAll(params?: { limit?: number; offset?: number }): Promise<Transaction[]> {
    return this.apiClient.get<Transaction[]>(API_CONSTANTS.transactions, { params });
  }

  async getByDateRange(start: Date, end: Date): Promise<Transaction[]> {
    return this.apiClient.get<Transaction[]>(API_CONSTANTS.transactionsRange, {
      params: {
        start_date: start.toISOString(),
        end_date: end.toISOString(),
      },
    });
  }

  async getStats(start: Date, end: Date): Promise<StatsResponse> {
    const backendData = await this.apiClient.get<BackendDetailedStatsResponse>(API_CONSTANTS.transactionsStatsDetailed, {
      params: {
        start_date: start.toISOString(),
        end_date: end.toISOString(),
      },
    });
    return transformStatsResponse(backendData);
  }

  async getById(id: string): Promise<Transaction> {
    return this.apiClient.get<Transaction>(`${API_CONSTANTS.transactions}/${id}`);
  }

  async create(data: TransactionCreateInput): Promise<Transaction> {
    return this.apiClient.post<Transaction>(API_CONSTANTS.transactions, data);
  }

  async update(id: string, data: TransactionUpdateInput): Promise<Transaction> {
    return this.apiClient.put<Transaction>(`${API_CONSTANTS.transactions}/${id}`, data);
  }

  async delete(id: string): Promise<void> {
    return this.apiClient.delete(`${API_CONSTANTS.transactions}/${id}`);
  }
}
