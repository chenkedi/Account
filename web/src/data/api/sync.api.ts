import { ApiClient } from './client';
import type { SyncPullRequest, SyncPullResponse, SyncPushRequest, SyncPushResponse } from '../models/sync.model';
import { API_CONSTANTS } from '../../core/constants/api.constants';

export class SyncApi {
  constructor(private apiClient: ApiClient) {}

  async pull(data: SyncPullRequest): Promise<SyncPullResponse> {
    return this.apiClient.post<SyncPullResponse>(API_CONSTANTS.syncPull, {
      deviceId: data.deviceId,
      lastSyncAt: data.lastSyncAt.toISOString(),
    });
  }

  async push(data: SyncPushRequest): Promise<SyncPushResponse> {
    return this.apiClient.post<SyncPushResponse>(API_CONSTANTS.syncPush, {
      deviceId: data.deviceId,
      accounts: data.accounts,
      categories: data.categories,
      transactions: data.transactions,
      lastSyncAt: data.lastSyncAt.toISOString(),
    });
  }
}
